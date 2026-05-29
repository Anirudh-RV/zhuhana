#!/usr/bin/env python3
"""
One-click load test for the AlgoNexus OrderHub order->broker->fill pipeline.

Brings up the pinned (cpus:1.0) Docker stack, resets to a clean state, runs k6 at a
target arrival rate, waits for the pipeline to drain, then collects ground-truth
evidence and saves it as JSON + Markdown.

Ground truth = redis XINFO GROUPS (lag / pending / entries-read), NOT log line
counts: the app's zap logger uses production sampling, so grepping log messages
UNDERCOUNTS. lag==0 & pending==0 & entries-read==accepted & errors==0 ⇒
end-to-end zero loss.

Usage:
  python run_loadtest.py --rate 5000
  python run_loadtest.py --rate 5000 --policy backpressure
  python run_loadtest.py --saturation              # ramp, find the ceiling
  python run_loadtest.py --rate 1000 --no-build    # skip image rebuild

Stdlib only (no pip install). Requires Docker on PATH.

Config sources, in order of precedence:
  1. CLI flags
  2. host env vars (LOADTEST_*)
  3. values read from env/local-env.env (REDIS_PASSWORD, PORT, ...)
  4. defaults in DEFAULTS below
"""
import argparse
import json
import os
import subprocess
import sys
import time
from datetime import datetime
from pathlib import Path

HERE = Path(__file__).resolve().parent          # .../Backend/algonexus/loadtest
ALGO = HERE.parent                              # .../Backend/algonexus
ENVF = ALGO / "env" / "local-env.env"
F1 = ALGO / "docker-compose" / "docker-compose-local.yml"
F2 = ALGO / "docker-compose" / "docker-compose-loadtest.yml"
RESULTS = HERE / "results"

# Defaults — every value here is overridable via LOADTEST_<KEY> env var or CLI flag.
DEFAULTS = {
    "PROJECT":         "algonexus",
    "APP_CONTAINER":   "algonexus-container",
    "REDIS_CONTAINER": "algonexus-redis-container",
    "STREAM":          "orderstream:strategy-1",
    "GROUP":           "group:order:strategy-1",
    "SUBMIT_PATH":     "/v1/ordermanager/submit",
}

ERR_PATTERNS = ["connection pool", "XAdd Failed", "unable to push", "Fail to handle message"]


# ---------------------------------------------------------------- shell helpers

def run(cmd):
    """Run a command, return (rc, combined_output)."""
    p = subprocess.run(cmd, capture_output=True, text=True, encoding="utf-8", errors="replace")
    return p.returncode, (p.stdout or "") + (p.stderr or "")


def die(msg):
    print(f"\n[ERROR] {msg}", file=sys.stderr)
    sys.exit(1)


# ---------------------------------------------------------------- env file load

def load_env_file(path):
    """Parse a KEY=VALUE .env file into a dict (ignoring blanks/comments)."""
    out = {}
    if not path.exists():
        return out
    for raw in path.read_text(encoding="utf-8").splitlines():
        line = raw.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        k, _, v = line.partition("=")
        out[k.strip()] = v.strip()
    return out


# ---------------------------------------------------------------- pure parsers

def parse_xinfo_groups(text, group_name=None):
    """
    Parse `redis-cli XINFO GROUPS <stream>` output into a dict.

    Output is alternating key/value lines (one line each). When multiple groups
    are returned, keys repeat; we keep the entry matching `group_name`, or the
    last one if no name match.
    """
    lines = [l for l in text.splitlines() if l.strip()]
    entries, current = [], {}
    for i in range(0, len(lines) - 1, 2):
        key, val = lines[i].strip(), lines[i + 1].strip()
        if key == "name" and current:
            entries.append(current)
            current = {}
        current[key] = val
    if current:
        entries.append(current)

    if not entries:
        return {"pending": None, "entries_read": None, "lag": None}

    chosen = next((e for e in entries if group_name and e.get("name") == group_name), entries[-1])

    def asint(k):
        try:
            return int(chosen.get(k))
        except (TypeError, ValueError):
            return None

    return {
        "pending":      asint("pending"),
        "entries_read": asint("entries-read"),
        "lag":          asint("lag"),
    }


def parse_k6_summary(obj):
    """
    Pull the fields we care about out of k6's --summary-export JSON.

    Handles both schemas k6 emits:
      - legacy --summary-export:  metrics.<name>.<stat>
      - --new-machine-readable-summary:  metrics.<name>.values.<stat>
    Times are in milliseconds (k6 native unit). Failure ratio is 0.0..1.0;
    legacy stores it under "value", new schema under "rate".
    """
    metrics = (obj or {}).get("metrics", {}) or {}

    def m(name, *stats, default=None):
        """Read first matching stat from either schema."""
        entry = metrics.get(name) or {}
        layers = [entry.get("values") or {}, entry]
        for layer in layers:
            for s in stats:
                if s in layer and layer[s] is not None:
                    return layer[s]
        return default

    accepted = m("http_reqs", "count")
    rate_s   = m("http_reqs", "rate")
    failed   = m("http_req_failed", "rate", "value")  # 0.0..1.0
    p50      = m("http_req_duration", "med")
    p90      = m("http_req_duration", "p(90)")
    p95      = m("http_req_duration", "p(95)")
    avg      = m("http_req_duration", "avg")
    maxv     = m("http_req_duration", "max")
    dropped  = m("dropped_iterations", "count")
    vus_max  = m("vus_max", "max", "value")

    return {
        "accepted_requests":   int(accepted) if accepted is not None else None,
        "accepted_rate_per_s": float(rate_s) if rate_s is not None else None,
        "http_failed_pct":     round(float(failed) * 100, 4) if failed is not None else None,
        "ingest_med_ms":       p50,
        "ingest_p90_ms":       p90,
        "ingest_p95_ms":       p95,
        "ingest_avg_ms":       avg,
        "ingest_max_ms":       maxv,
        "dropped_iterations":  int(dropped) if dropped is not None else None,
        "vus_max":             int(vus_max) if vus_max is not None else None,
    }


# ---------------------------------------------------------------- context obj

class Ctx:
    """Bundles the resolved config so we don't smuggle globals into every helper."""
    def __init__(self, args, env_file):
        # Layered resolve: CLI > host env (LOADTEST_KEY) > env_file > defaults.
        def cfg(key, fallback=None):
            return (
                os.environ.get(f"LOADTEST_{key}")
                or env_file.get(f"LOADTEST_{key}")
                or fallback
            )

        self.project        = cfg("PROJECT",         DEFAULTS["PROJECT"])
        self.app_container  = cfg("APP_CONTAINER",   DEFAULTS["APP_CONTAINER"])
        self.redis_cont     = cfg("REDIS_CONTAINER", DEFAULTS["REDIS_CONTAINER"])
        self.stream         = cfg("STREAM",          DEFAULTS["STREAM"])
        self.group          = cfg("GROUP",           DEFAULTS["GROUP"])
        self.submit_path    = cfg("SUBMIT_PATH",     DEFAULTS["SUBMIT_PATH"])

        # Pulled straight from the env_file so we share one source of truth.
        self.redis_password = env_file.get("REDIS_PASSWORD", "password")
        self.app_port       = env_file.get("PORT", "8080")

        self.network        = f"{self.project}_default"
        self.target_url     = f"http://algonexus:{self.app_port}"
        self.args           = args

    def compose(self, *a):
        return run(["docker", "compose", "-p", self.project, "--env-file", str(ENVF),
                    "-f", str(F1), "-f", str(F2), *a])

    def docker(self, *a):
        return run(["docker", *a])

    def redis_cli(self, *a):
        rc, out = self.docker("exec", self.redis_cont, "redis-cli",
                              "-a", self.redis_password, *a)
        # strip redis-cli's "Warning: Using a password..." stderr line
        return "\n".join(l for l in out.splitlines() if "Using a password" not in l).strip()


# ---------------------------------------------------------------- live probes

def wait_ready(ctx, timeout=180):
    """The app creates the consumer group on startup, so this readiness check
    works even with LOG_LEVEL=error."""
    print("  waiting for app/orderhub ready", end="", flush=True)
    t0 = time.time()
    while time.time() - t0 < timeout:
        if ctx.group in ctx.redis_cli("XINFO", "GROUPS", ctx.stream):
            print(" - ready")
            return
        print(".", end="", flush=True)
        time.sleep(2)
    die(f"app did not become ready (consumer group {ctx.group} not on {ctx.stream})")


def stats(ctx):
    """Return (mem_used_str, cpu_float). Drops `docker stats` MemUsage's
    `/ <limit>` tail — the limit is the host RAM when no cgroup mem cap is set,
    which is misleading on the report."""
    _, out = ctx.docker("stats", "--no-stream",
                        "--format", "{{.MemUsage}};{{.CPUPerc}}",
                        ctx.app_container)
    line = out.strip().splitlines()[-1] if out.strip() else "?;0%"
    mem, _, cpu = line.partition(";")
    mem_used = mem.partition("/")[0].strip() or mem.strip()
    try:
        cpu_f = float(cpu.replace("%", "").strip())
    except ValueError:
        cpu_f = 0.0
    return mem_used, cpu_f


def wait_drain(ctx, timeout=150):
    print("  waiting for pipeline drain", end="", flush=True)
    t0 = time.time()
    while time.time() - t0 < timeout:
        if stats(ctx)[1] < 5:
            print(" - idle")
            return
        print(".", end="", flush=True)
        time.sleep(3)
    print(" - (still busy)")


def count_errors(ctx):
    _, logs = ctx.docker("logs", ctx.app_container)
    n = 0
    for line in logs.splitlines():
        if any(p in line for p in ERR_PATTERNS):
            n += 1
    return n


# ---------------------------------------------------------------- k6

def run_k6(ctx, script, env_overrides):
    """
    Run k6 in a sidecar container; write its machine-readable summary to a host
    file via -v mount, then load+parse it.
    """
    summary_name = f"k6-summary-{datetime.now().strftime('%Y%m%d-%H%M%S')}.json"
    summary_host = RESULTS / summary_name

    cmd = ["docker", "run", "--rm", "--network", ctx.network,
           "-v", f"{HERE.as_posix()}:/scripts",
           "-v", f"{RESULTS.as_posix()}:/out"]
    for k, v in env_overrides.items():
        cmd += ["-e", f"{k}={v}"]
    cmd += ["grafana/k6", "run",
            "--summary-export", f"/out/{summary_name}",
            f"/scripts/{script}"]

    rc, out = run(cmd)
    # echo k6's headline lines for the operator
    for line in out.splitlines():
        if any(tok in line for tok in ("http_req", "iterations", "dropped", "vus_max", "p(9")):
            print("  " + line.strip())

    summary = {}
    if summary_host.exists():
        try:
            summary = json.loads(summary_host.read_text(encoding="utf-8"))
        except json.JSONDecodeError as e:
            print(f"  WARN: could not parse k6 summary JSON: {e}")
    else:
        print("  WARN: k6 did not produce a summary file (rc=%d)" % rc)
    return summary, summary_host


# ---------------------------------------------------------------- main

def main():
    ap = argparse.ArgumentParser(description="One-click AlgoNexus OrderHub load test.")
    ap.add_argument("--rate", type=int, default=5000, help="constant arrival rate (orders/s)")
    ap.add_argument("--duration", default="30s", help="constant-test hold duration")
    ap.add_argument("--policy", choices=["failfast", "backpressure"], default="failfast")
    ap.add_argument("--logs", choices=["off", "on"], default="off",
                    help="app logging: off (default) is the #1 throughput lever; on = gin debug + INFO logs")
    ap.add_argument("--pre-vus", type=int, default=None,
                    help="k6 preAllocatedVUs (default: max(200, rate*0.4))")
    ap.add_argument("--max-vus", type=int, default=None,
                    help="k6 maxVUs (default: max(pre, rate*2))")
    ap.add_argument("--saturation", action="store_true",
                    help="run ramp instead of constant probe")
    ap.add_argument("--no-build", action="store_true", help="do not rebuild the app image")
    ap.add_argument("--no-up", action="store_true", help="assume the stack is already up")
    args = ap.parse_args()

    if not ENVF.exists() or not F1.exists() or not F2.exists():
        die(f"compose/env files not found under {ALGO} (run from the repo).")
    if run(["docker", "version"])[0] != 0:
        die("Docker not available on PATH. Start Docker Desktop / Engine.")

    env_file = load_env_file(ENVF)
    ctx = Ctx(args, env_file)

    script = "submit_saturation.js" if args.saturation else "submit_constant.js"

    # App-side env (compose interpolates these into the algonexus service)
    os.environ["BROKER_OVERFLOW_POLICY"] = args.policy
    if args.logs == "on":
        os.environ.update(GIN_MODE="debug", LOG_LEVEL="info", HTTP_LOG="on")
    else:
        os.environ.update(GIN_MODE="release", LOG_LEVEL="error", HTTP_LOG="off")
    RESULTS.mkdir(exist_ok=True)
    ts = datetime.now().strftime("%Y%m%d-%H%M%S")

    print(f"=== AlgoNexus load test | script={script} rate={args.rate} policy={args.policy} ===")

    if not args.no_up:
        print("=== 1) bring up stack (cpus:1.0) ===")
        rc, out = ctx.compose(*(["up", "-d"] + ([] if args.no_build else ["--build"])))
        if rc != 0:
            print(out); die("compose up failed")
        wait_ready(ctx)

    print(f"=== 2) reset clean (flush redis + recreate app, policy={args.policy}) ===")
    ctx.redis_cli("FLUSHALL")
    rc, out = ctx.compose("up", "-d", "--force-recreate", "--no-deps", "algonexus")
    if rc != 0:
        print(out); die("compose recreate failed")
    wait_ready(ctx)
    mem_before, _ = stats(ctx)
    print(f"  idle baseline: mem={mem_before}")

    print(f"=== 3) k6 {script} (rate={args.rate} dur={args.duration}) ===")
    k6_env = {
        "TARGET": ctx.target_url,
        "PATH_":  ctx.submit_path,
        "RATE":   str(args.rate),
        "DUR":    args.duration,
    }
    if args.pre_vus is not None:
        k6_env["PRE_VUS"] = str(args.pre_vus)
    if args.max_vus is not None:
        k6_env["MAX_VUS"] = str(args.max_vus)
    summary, summary_path = run_k6(ctx, script, k6_env)
    k6m = parse_k6_summary(summary)

    wait_drain(ctx)

    print("=== 4) evidence (ground truth) ===")
    xinfo_raw = ctx.redis_cli("XINFO", "GROUPS", ctx.stream)
    xinfo = parse_xinfo_groups(xinfo_raw, ctx.group)
    xlen_raw = ctx.redis_cli("XLEN", ctx.stream)
    mem_after, cpu_after = stats(ctx)
    errs = count_errors(ctx)

    accepted = k6m.get("accepted_requests")
    lag      = xinfo.get("lag")
    pending  = xinfo.get("pending")
    entries  = xinfo.get("entries_read")
    xlen     = int(xlen_raw) if xlen_raw.isdigit() else xlen_raw

    # Inline echo so the operator sees the raw evidence in the console (not just
    # in the saved report). XLEN = how many entries are on the stream; entries-read
    # = how many the consumer group has read+ACKed; pending = read-but-not-ACKed;
    # lag = stream backlog the group still owes (XLEN - last-delivered position).
    print(f"  redis stream  : {ctx.stream}")
    print(f"  consumer group: {ctx.group}")
    print("  XINFO GROUPS:")
    for line in xinfo_raw.splitlines():
        print(f"    {line}")
    print(f"  XLEN {ctx.stream} = {xlen}")
    print(f"  -> lag={lag}  pending={pending}  entries-read={entries}  xlen={xlen}  "
          f"accepted={accepted}  errors={errs}")
    print(f"  -> memory: {mem_before} -> {mem_after}  cpu_after={cpu_after}%")

    zero_loss = (
        lag == 0 and pending == 0 and errs == 0
        and entries is not None and accepted is not None
        and entries >= accepted
    )

    result = {
        "timestamp": ts, "script": script, "policy": args.policy, "logs": args.logs,
        "target_rate": args.rate, "duration": args.duration,
        "target_url": ctx.target_url + ctx.submit_path,
        "k6": k6m,
        "redis": {**xinfo, "xlen": xlen,
                  "stream": ctx.stream, "group": ctx.group},
        "memory": {"idle_before": mem_before, "after": mem_after, "cpu_after_pct": cpu_after},
        "pipeline_errors": errs,
        "verdict_zero_loss": bool(zero_loss),
        "k6_summary_file": summary_path.name,
    }

    (RESULTS / f"result_{ts}.json").write_text(json.dumps(result, indent=2), encoding="utf-8")

    def fmt_ms(x):
        return f"{x:.2f}ms" if isinstance(x, (int, float)) else str(x)
    md = f"""# Load-test result {ts}

| metric | value |
|---|---|
| target | `{ctx.target_url}{ctx.submit_path}` |
| script / policy / logs | `{script}` / `{args.policy}` / `{args.logs}` |
| target rate | {args.rate}/s |
| **accepted (HTTP)** | {accepted} reqs @ {k6m.get('accepted_rate_per_s')}/s |
| HTTP failed | {k6m.get('http_failed_pct')}% |
| ingest latency med / p90 / p95 | {fmt_ms(k6m.get('ingest_med_ms'))} / {fmt_ms(k6m.get('ingest_p90_ms'))} / {fmt_ms(k6m.get('ingest_p95_ms'))} |
| dropped iterations (k6 VU cap) | {k6m.get('dropped_iterations')} |
| **redis stream / group** | `{ctx.stream}` / `{ctx.group}` |
| **XINFO GROUPS - lag** | {lag} |
| XINFO GROUPS - pending | {pending} |
| **XINFO GROUPS - entries-read** | {entries} |
| XLEN (stream length) | {xlen} |
| memory idle -> after | {mem_before} -> {mem_after} (cpu {cpu_after}%) |
| pipeline errors | {errs} |
| **VERDICT: end-to-end zero loss** | {'YES' if zero_loss else 'NO (see lag/errors)'} |

```
$ redis-cli XINFO GROUPS {ctx.stream}
{xinfo_raw}

$ redis-cli XLEN {ctx.stream}
{xlen}
```

Ground truth = redis XINFO (lag/pending/entries-read) + XLEN, not sampled log counts.
Zero-loss criterion: lag==0 AND pending==0 AND entries-read>=accepted AND errors==0.
"""
    (RESULTS / f"result_{ts}.md").write_text(md, encoding="utf-8")

    print("\n" + md)
    print(f"saved: {RESULTS / f'result_{ts}.json'}")
    print(f"saved: {RESULTS / f'result_{ts}.md'}")
    print(f"\nVERDICT: {'end-to-end ZERO LOSS' if zero_loss else 'NOT clean - inspect lag/pending/errors'}")


if __name__ == "__main__":
    main()
