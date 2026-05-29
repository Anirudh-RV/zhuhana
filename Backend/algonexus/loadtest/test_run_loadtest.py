"""Unit tests for the pure parsers in run_loadtest.py (no Docker needed).

Samples are real output captured from actual runs in this project. Run with:
    pytest Backend/algonexus/loadtest/test_run_loadtest.py
"""
import json

import run_loadtest as rl


# Real `redis-cli XINFO GROUPS orderstream:strategy-1` output (clean run, zero lag).
XINFO_CLEAN = """name
group:order:strategy-1
consumers
1
pending
0
last-delivered-id
1780090695343-2
entries-read
150001
lag
0"""

# Hypothetical multi-group case — picks the one whose name matches.
XINFO_TWO_GROUPS = """name
group:order:strategy-1
consumers
1
pending
0
last-delivered-id
1780090695343-2
entries-read
150001
lag
0
name
group:other:strategy-1
consumers
1
pending
42
last-delivered-id
1780090695343-2
entries-read
99
lag
7"""


def test_parse_xinfo_zero_loss():
    x = rl.parse_xinfo_groups(XINFO_CLEAN, group_name="group:order:strategy-1")
    assert x == {"pending": 0, "entries_read": 150001, "lag": 0}


def test_parse_xinfo_no_group_name_falls_back_to_last():
    # When no group_name supplied, the last entry parsed wins (single-group case unaffected).
    x = rl.parse_xinfo_groups(XINFO_CLEAN, group_name=None)
    assert x == {"pending": 0, "entries_read": 150001, "lag": 0}


def test_parse_xinfo_selects_matching_group():
    x = rl.parse_xinfo_groups(XINFO_TWO_GROUPS, group_name="group:order:strategy-1")
    assert x == {"pending": 0, "entries_read": 150001, "lag": 0}
    y = rl.parse_xinfo_groups(XINFO_TWO_GROUPS, group_name="group:other:strategy-1")
    assert y == {"pending": 42, "entries_read": 99, "lag": 7}


def test_parse_xinfo_empty_or_garbage():
    assert rl.parse_xinfo_groups("", "anything") == {"pending": None, "entries_read": None, "lag": None}
    assert rl.parse_xinfo_groups("nonsense\nblob\n", "x") == {
        "pending": None, "entries_read": None, "lag": None,
    }


# Minimal k6 --summary-export shape; values lifted from a real 5k constant run.
K6_SUMMARY_CONSTANT = {
    "metrics": {
        "http_reqs": {"values": {"count": 150001, "rate": 4503.822979}},
        "http_req_failed": {"values": {"rate": 0.0, "passes": 150001, "fails": 0}},
        "http_req_duration": {"values": {
            "avg": 12.1, "min": 0.04185, "med": 0.58127,
            "max": 2110.0, "p(90)": 37.69, "p(95)": 43.4,
        }},
        "vus_max": {"values": {"value": 2000, "max": 2000, "min": 2000}},
    }
}

K6_SUMMARY_SATURATION = {
    "metrics": {
        "http_reqs": {"values": {"count": 696745, "rate": 6901.362531}},
        "http_req_failed": {"values": {"rate": 0.0}},
        "http_req_duration": {"values": {
            "avg": 419.23, "med": 78.11, "p(90)": 1080.0, "p(95)": 1910.0, "max": 6790.0,
        }},
        "dropped_iterations": {"values": {"count": 460570, "rate": 4562.014139}},
    }
}


def test_parse_k6_constant():
    m = rl.parse_k6_summary(K6_SUMMARY_CONSTANT)
    assert m["accepted_requests"] == 150001
    assert abs(m["accepted_rate_per_s"] - 4503.822979) < 0.01
    assert m["http_failed_pct"] == 0.0
    assert m["ingest_med_ms"] == 0.58127
    assert m["ingest_p90_ms"] == 37.69
    assert m["ingest_p95_ms"] == 43.4
    assert m["dropped_iterations"] is None
    assert m["vus_max"] == 2000


def test_parse_k6_saturation_dropped():
    m = rl.parse_k6_summary(K6_SUMMARY_SATURATION)
    assert m["accepted_requests"] == 696745
    assert m["http_failed_pct"] == 0.0
    assert m["ingest_p95_ms"] == 1910.0
    assert m["dropped_iterations"] == 460570


# Real --summary-export from this k6 image — note: stats live directly on the
# metric entry, not under `.values`, and http_req_failed uses `value` not `rate`.
K6_SUMMARY_LEGACY_SHAPE = {
    "metrics": {
        "http_reqs": {"count": 149999, "rate": 4519.06},
        "http_req_failed": {"passes": 0, "fails": 149999, "value": 0.0},
        "http_req_duration": {
            "avg": 8.18, "min": 0.035, "med": 0.18,
            "max": 2111.4, "p(90)": 24.29, "p(95)": 30.48,
        },
        "vus_max": {"value": 2000, "min": 2000, "max": 2000},
    }
}


def test_parse_k6_legacy_shape():
    m = rl.parse_k6_summary(K6_SUMMARY_LEGACY_SHAPE)
    assert m["accepted_requests"] == 149999
    assert m["http_failed_pct"] == 0.0
    assert m["ingest_p95_ms"] == 30.48
    assert m["vus_max"] == 2000


def test_parse_k6_empty():
    # missing/None inputs return all-None, never raise.
    m = rl.parse_k6_summary({})
    assert m["accepted_requests"] is None
    assert m["http_failed_pct"] is None
    assert m["ingest_p95_ms"] is None
    m2 = rl.parse_k6_summary(None)
    assert m2["accepted_requests"] is None


def test_load_env_file(tmp_path):
    p = tmp_path / "x.env"
    p.write_text("# a comment\n\nPORT=8080\nREDIS_PASSWORD=secret\n  KEY = val \n", encoding="utf-8")
    env = rl.load_env_file(p)
    assert env["PORT"] == "8080"
    assert env["REDIS_PASSWORD"] == "secret"
    assert env["KEY"] == "val"
    # missing file -> empty dict, no raise
    assert rl.load_env_file(tmp_path / "nope.env") == {}
