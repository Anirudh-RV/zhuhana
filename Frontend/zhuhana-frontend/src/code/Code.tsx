import { useEffect, useRef, useState } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import AppTheme from "../shared-ui-theme/AppTheme";
import MonacoEditor from "./components/MonacoEditor";
import CodeSideMenu from "./components/CodeSideMenu";
import LLMPanel from "./components/LLMPanel";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import { green } from "@mui/material/colors";
import type { Monaco } from "@monaco-editor/react";

declare global {
  interface Window {
    loadPyodide: (config: {
      indexURL: string;
      stdout?: (msg: string) => void;
      stderr?: (msg: string) => void;
    }) => Promise<any>;
    pyodide?: any;
  }
}

type TerminalLine = { text: string; type: "info" | "success" | "error" };

type Message = {
  role: "user" | "assistant" | "system";
  content: string;
};

const defaultPythonCode = `from algorithm.models import OrderInstruction

class ZhuhanaStrategy:
    def __init__(self, zhuhana_sdk):
        self.zhuhana_sdk = zhuhana_sdk

    def on_data(self, current_data):
       pass

    def condition_for_sell(self, current_data) -> OrderInstruction:
        pass

    def condition_for_buy(self, current_data) -> OrderInstruction:
        pass
`;

export default function CodeEditorDashboard(props: {
  disableCustomTheme?: boolean;
}) {
  useEffect(() => {
    document.title = "Zhuhana - Algorithm IDE";
  }, []);

  const [code, setCode] = useState(defaultPythonCode);
  const [terminalOutput, setTerminalOutput] = useState<TerminalLine[]>([
    { text: ">> Terminal ready...", type: "info" },
  ]);
  const pyodideInstanceRef = useRef<any>(null);
  const [isLoadingPyodide, setIsLoadingPyodide] = useState(true);

  const containerRef = useRef<HTMLDivElement>(null);
  const dragInfo = useRef<{ startY: number; startHeight: number } | null>(null);
  const [errorLines, setErrorLines] = useState<number[]>([]);
  const [editorHeight, setEditorHeight] = useState(
    () => window.innerHeight * 0.6
  );

  const appendStdout = useRef((msg: string) => {
    setTerminalOutput((prev) => [...prev, { text: msg, type: "success" }]);
  });

  const appendStderr = useRef((msg: string) => {
    setTerminalOutput((prev) => [
      ...prev,
      { text: `[Python Error]: ${msg}`, type: "error" },
    ]);
  });

  const editorRef = useRef<any>(null);
  const monacoRef = useRef<Monaco | null>(null);
  const [decorations, setDecorations] = useState<string[]>([]);

  useEffect(() => {
    const PYODIDE_BASE_URL = "https://cdn.jsdelivr.net/pyodide/v0.26.1/full/";
    const script = document.createElement("script");
    script.src = PYODIDE_BASE_URL + "pyodide.js";
    script.async = true;
    script.setAttribute("data-pyodide-base-url", PYODIDE_BASE_URL);
    document.head.appendChild(script);

    script.onload = async () => {
      try {
        if (typeof window.loadPyodide === "function") {
          const pyodide = await window.loadPyodide({
            indexURL: PYODIDE_BASE_URL,
            stdout: appendStdout.current,
            stderr: appendStderr.current,
          });
          pyodideInstanceRef.current = pyodide;
          setTerminalOutput((prev) => [
            ...prev,
            {
              text: ">> Python runtime loaded. Run the script to test...",
              type: "info",
            },
          ]);
        } else {
          throw new Error("window.loadPyodide is not a function.");
        }
      } catch (err: any) {
        setTerminalOutput((prev) => [
          ...prev,
          {
            text: `>> Failed to initialize Pyodide: ${err.message || err}`,
            type: "error",
          },
        ]);
      } finally {
        setIsLoadingPyodide(false);
      }
    };

    script.onerror = (err) => {
      setTerminalOutput((prev) => [
        ...prev,
        { text: `>> Failed to load Pyodide: ${err}`, type: "error" },
      ]);
      setIsLoadingPyodide(false);
    };

    const handleMouseMove = (e: MouseEvent) => {
      if (!dragInfo.current) return;
      const delta = e.clientY - dragInfo.current.startY;
      const newHeight = Math.max(100, dragInfo.current.startHeight + delta);
      setEditorHeight(newHeight);
    };

    const handleMouseUp = () => {
      dragInfo.current = null;
      document.body.style.cursor = "default";
      document.body.style.userSelect = "auto";
    };

    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);

    return () => {
      document.head.removeChild(script);
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  const handleRunCode = async () => {
    if (!pyodideInstanceRef.current) {
      setTerminalOutput((prev) => [
        ...prev,
        { text: ">> Python runtime not ready yet.", type: "info" },
      ]);
      return;
    }

    setTerminalOutput([
      { text: ">> Terminal ready...", type: "info" },
      { text: ">> Executing Python code...", type: "info" },
    ]);

    // Clear previous decorations
    if (editorRef.current) {
      setDecorations(editorRef.current.deltaDecorations(decorations, []));
    }

    console.log("[Current Code]:\n" + code);
    console.log("[Code Lines]:", code.split("\n").length);

    try {
      await pyodideInstanceRef.current.runPythonAsync(code);
      setTerminalOutput((prev) => [
        ...prev,
        { text: ">> Code execution finished.", type: "success" },
      ]);
    } catch (error: any) {
      const message = error.message || String(error);
      console.error("[Pyodide Error]", message);

      setTerminalOutput((prev) => [
        ...prev,
        { text: `>> Error: ${message}`, type: "error" },
      ]);

      const lineMatch = message.match(/File "<exec>", line (\d+)/);
      console.log("[Parsed lineMatch]:", lineMatch);

      if (lineMatch && editorRef.current && monacoRef.current) {
        const lineNumber = parseInt(lineMatch[1], 10);
        setErrorLines([lineNumber]);
        const model = editorRef.current.getModel();
        const lineText = model?.getLineContent(lineNumber);
        console.log(`[Content at line ${lineNumber}]:`, lineText);

        const range = new monacoRef.current.Range(lineNumber, 1, lineNumber, 1);
        console.log("[Monaco Range]:", range);

        const newDecorations = editorRef.current.deltaDecorations(decorations, [
          {
            range,
            options: {
              isWholeLine: true,
              className: "errorLineHighlight",
              glyphMarginClassName: "errorGlyphMargin",
              hoverMessage: { value: `**SyntaxError**: ${message}` },
            },
          },
        ]);
        setDecorations(newDecorations);
      }
    }
  };

  const handleSendToLLM = async (
    messages: Message[],
    onChunk: (token: string) => void,
    signal: AbortSignal
  ) => {
    try {
      // 👇 Format messages into a full prompt string
      const prompt = messages
        .map(
          (msg) =>
            `${msg.role === "user" ? "User" : "Assistant"}: ${msg.content}`
        )
        .join("\n");

      const response = await fetch(
        `http://localhost:3000/v1/ask?q=${encodeURIComponent(prompt)}`,
        { signal }
      );

      const reader = response.body?.getReader();
      const decoder = new TextDecoder();

      if (!reader) {
        onChunk("[Error: No stream]");
        return;
      }

      while (true) {
        if (signal.aborted) {
          try {
            await reader.cancel();
          } catch {}
          throw new DOMException("Aborted", "AbortError");
        }

        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        const lines = chunk.split("\n").filter((line) => line.trim() !== "");

        for (const line of lines) {
          try {
            const data = JSON.parse(line);
            if (data.done) return;
            onChunk(data.response || "");
          } catch {
            onChunk(chunk);
          }
        }
      }
    } catch (err: any) {
      if (err.name === "AbortError") {
        console.log("Request aborted");
      } else {
        onChunk(`[LLM Error]: ${err.message}`);
      }
    }
  };

  const handleCodeChange = (newCode: string | undefined) => {
    setCode(newCode ?? "");

    // Clear decorations and error lines on code change
    if (editorRef.current) {
      setDecorations(editorRef.current.deltaDecorations(decorations, []));
    }
    setErrorLines([]); // 👈 This resets the ❗ on next render
  };

  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: "flex", height: "100vh" }}>
        <CodeSideMenu />

        <Box
          ref={containerRef}
          sx={{
            flex: 2,
            p: 1,
            display: "flex",
            flexDirection: "column",
            minWidth: 0,
          }}
        >
          <Box
            sx={{
              height: `${editorHeight}px`,
              border: "1px solid #ccc",
              borderRadius: 1,
            }}
          >
            <MonacoEditor
              code={code}
              onChange={handleCodeChange}
              onMount={(editor, monaco) => {
                editorRef.current = editor;
                monacoRef.current = monaco;
              }}
              errorLines={errorLines}
            />
          </Box>

          <Divider
            sx={{
              height: "6px",
              backgroundColor: "divider",
              my: 0.5,
              cursor: "row-resize",
            }}
            onMouseDown={(e) => {
              dragInfo.current = {
                startY: e.clientY,
                startHeight: editorHeight,
              };
              document.body.style.cursor = "row-resize";
              document.body.style.userSelect = "none";
            }}
          />

          <Box
            sx={{
              flexGrow: 1,
              display: "flex",
              flexDirection: "column",
              border: "1px solid #333",
              borderRadius: 1,
              backgroundColor: "#111",
              minHeight: "100px",
            }}
          >
            <Toolbar
              variant="dense"
              sx={{
                backgroundColor: "#222",
                borderBottom: "1px solid #444",
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                px: 1,
              }}
            >
              <Typography variant="subtitle2" sx={{ color: "#bbb" }}>
                Terminal
              </Typography>
              <IconButton
                onClick={handleRunCode}
                size="small"
                disabled={isLoadingPyodide}
                sx={{
                  color: isLoadingPyodide ? "grey" : green[500],
                  "&:hover": {
                    backgroundColor: isLoadingPyodide
                      ? "transparent"
                      : green[900],
                  },
                }}
              >
                <PlayArrowIcon />
              </IconButton>
            </Toolbar>

            <Box sx={{ flexGrow: 1, p: 1, overflowY: "auto" }}>
              {terminalOutput.map((line, index) => (
                <Typography
                  key={index}
                  sx={{
                    color:
                      line.type === "success"
                        ? "#0f0"
                        : line.type === "error"
                        ? "#f55"
                        : "#aaa",
                    fontFamily: "monospace",
                    fontSize: "0.875rem",
                  }}
                >
                  {line.text}
                </Typography>
              ))}
            </Box>
          </Box>
        </Box>

        <Box
          sx={{
            width: "25%",
            display: "flex",
            flexDirection: "column",
            borderLeft: "1px solid",
            borderColor: "divider",
            backgroundColor: "background.paper",
            p: 2,
          }}
        >
          <LLMPanel onSend={handleSendToLLM} />
        </Box>
      </Box>
    </AppTheme>
  );
}
