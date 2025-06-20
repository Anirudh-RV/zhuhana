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

// Declare global window properties for Pyodide
declare global {
  interface Window {
    loadPyodide: (config: {
      indexURL: string;
      stdout?: (msg: string) => void;
      stderr?: (msg: string) => void;
    }) => Promise<any>;
    pyodide?: any; // pyodide instance will be available here after loadPyodide resolves
  }
}

const defaultPythonCode = `def greet(name):\n    return f"Hello, {name}"\n\nprint(greet("World"))`;

export default function CodeEditorDashboard(props: {
  disableCustomTheme?: boolean;
}) {
  useEffect(() => {
    document.title = "Zhuhana - Algorithm IDE";
  }, []);

  const [code, setCode] = useState(defaultPythonCode);
  const [terminalOutput, setTerminalOutput] = useState(
    ">> Terminal ready...\n"
  );
  const pyodideInstanceRef = useRef<any>(null); // Using useRef for the Pyodide instance
  const [isLoadingPyodide, setIsLoadingPyodide] = useState(true);
  const [llmOutput, setLlmOutput] = useState("");

  const containerRef = useRef<HTMLDivElement>(null);
  const dragInfo = useRef<{ startY: number; startHeight: number } | null>(null);
  const [editorHeight, setEditorHeight] = useState(
    () => window.innerHeight * 0.6
  );

  // Use useRef for the output functions to ensure stable references for Pyodide
  const appendStdout = useRef((msg: string) => {
    setTerminalOutput((prev) => prev + msg + "\n");
  });
  const appendStderr = useRef((msg: string) => {
    setTerminalOutput((prev) => prev + `[Python Error]: ${msg}\n`);
  });

  // Effect for Pyodide initialization and script loading
  useEffect(() => {
    // Prevent re-initialization if already loaded (important for React Strict Mode)
    if (pyodideInstanceRef.current) {
      setIsLoadingPyodide(false);
      console.log("Pyodide already initialized in a previous render.");
      return;
    }

    setTerminalOutput((prev) => prev);
    setIsLoadingPyodide(true);

    const PYODIDE_BASE_URL = "https://cdn.jsdelivr.net/pyodide/v0.26.1/full/";
    const PYODIDE_SCRIPT_URL = PYODIDE_BASE_URL + "pyodide.js";

    // Create and append the script element
    const script = document.createElement("script");
    script.src = PYODIDE_SCRIPT_URL;
    script.async = true; // Load asynchronously
    script.type = "text/javascript"; // Explicitly set type

    // IMPORTANT: Set data-pyodide-base-url to help Pyodide's internal loader
    // resolve its JavaScript dependencies (like stackframe.js) correctly from the CDN.
    script.setAttribute("data-pyodide-base-url", PYODIDE_BASE_URL);

    document.head.appendChild(script);

    // Event listener for when the script has loaded
    script.onload = async () => {
      setTerminalOutput((prev) => prev);
      try {
        if (typeof window.loadPyodide === "function") {
          const pyodide = await window.loadPyodide({
            indexURL: PYODIDE_BASE_URL, // This indexURL is for Pyodide's internal packages and WASM
            stdout: appendStdout.current,
            stderr: appendStderr.current,
          });
          pyodideInstanceRef.current = pyodide; // Store instance in ref
          setTerminalOutput(
            (prev) =>
              prev + ">> Python runtime loaded. Run the script to test...\n"
          );
          console.log("Pyodide successfully loaded and initialized.");
        } else {
          throw new Error(
            "window.loadPyodide is not a function after script load."
          );
        }
      } catch (err: any) {
        setTerminalOutput(
          (prev) =>
            prev + `>> Failed to initialize Pyodide: ${err.message || err}\n`
        );
        console.error("Pyodide initialization error:", err);
      } finally {
        setIsLoadingPyodide(false);
      }
    };

    // Event listener for script loading errors
    script.onerror = (err) => {
      setTerminalOutput(
        (prev) =>
          prev +
          `>> Failed to load Pyodide script from CDN: ${
            err || "Unknown error"
          }\n`
      );
      console.error("Pyodide script loading error:", err);
      setIsLoadingPyodide(false);
    };

    // --- Resizable editor/terminal logic (remains unchanged) ---
    const handleMouseMove = (e: MouseEvent) => {
      if (!dragInfo.current) return;
      const delta = e.clientY - dragInfo.current.startY;
      let newHeight = dragInfo.current.startHeight + delta;

      if (containerRef.current) {
        const middleHeight = containerRef.current.clientHeight;
        const toolbarHeight = 48;
        const dividerHeight = 6;
        const padding = 16;
        const minTerminal = 80;
        const maxEditorHeight =
          middleHeight - toolbarHeight - dividerHeight - padding - minTerminal;
        newHeight = Math.max(100, Math.min(newHeight, maxEditorHeight));
      }
      setEditorHeight(newHeight);
    };

    const handleMouseUp = () => {
      dragInfo.current = null;
      document.body.style.cursor = "default";
      document.body.style.userSelect = "auto";
    };

    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);

    // Cleanup function for useEffect
    return () => {
      // Remove the dynamically added script when the component unmounts
      if (document.head.contains(script)) {
        document.head.removeChild(script);
      }
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
    };
  }, []); // Empty dependency array ensures this runs once on mount

  const handleRunCode = async () => {
    if (isLoadingPyodide || !pyodideInstanceRef.current) {
      setTerminalOutput((prev) => prev + ">> Python runtime not ready yet.\n");
      return;
    }

    setTerminalOutput((prev) => prev + "\n>> Executing Python code...\n");
    try {
      // Clear previous outputs before running again, but keep initial message
      setTerminalOutput(">> Terminal ready...\n>> Executing Python code...\n");
      await pyodideInstanceRef.current.runPythonAsync(code);
      setTerminalOutput((prev) => prev + ">> Code execution finished.\n");
    } catch (error: any) {
      setTerminalOutput(
        (prev) => prev + `>> Error: ${error.message || error}\n`
      );
      console.error("Pyodide execution error:", error);
    }
  };

  const handleSendToLLM = (input: string, onChunk: (token: string) => void) => {
    setLlmOutput("");
    // Note: This LLM endpoint URL (`http://localhost:8002/stream-llm`) might need
    // to be adjusted based on your actual backend setup.
    fetch("http://localhost:8002/stream-llm", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ input }),
    })
      .then(async (res) => {
        const reader = res.body?.getReader();
        const decoder = new TextDecoder();

        if (!reader) {
          console.error("No readable stream from LLM response.");
          return;
        }

        while (true) {
          const { done, value } = await reader.read();
          if (done) break;

          const chunk = decoder.decode(value);
          for (const line of chunk.split("\n")) {
            if (line.startsWith("data: ")) {
              const token = line.replace("data: ", "");
              setLlmOutput((prev) => prev + token);
              onChunk(token); // Pass chunk to the LLM panel for real-time display
            }
          }
        }
      })
      .catch((error) => {
        console.error("Error fetching from LLM:", error);
        setLlmOutput(
          (prev) =>
            prev + `\nError communicating with LLM: ${error.message || error}`
        );
      });
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
            height: "100%",
          }}
        >
          <Box
            sx={{
              height: `${editorHeight}px`,
              border: "1px solid #ccc",
              borderRadius: 1,
              overflow: "hidden",
            }}
          >
            <MonacoEditor code={code} onChange={(v) => setCode(v ?? "")} />
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
              overflow: "hidden",
              minHeight: "100px",
            }}
          >
            <Toolbar
              variant="dense"
              sx={{
                minHeight: "48px",
                backgroundColor: "#222",
                borderBottom: "1px solid #444",
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                px: 1,
                borderTopLeftRadius: "inherit",
                borderTopRightRadius: "inherit",
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

            <Box
              sx={{
                flexGrow: 1,
                color: "#0f0",
                fontFamily: "monospace",
                fontSize: "0.875rem",
                overflowY: "auto",
                whiteSpace: "pre-wrap",
                p: 1,
              }}
            >
              {isLoadingPyodide ? (
                <Box sx={{ p: 1, color: "#aaa" }}>
                  Loading Python runtime... please wait.
                  <br />
                  (Check browser network tab for WASM loading errors.)
                </Box>
              ) : (
                terminalOutput
              )}
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
          <LLMPanel output={llmOutput} onSend={handleSendToLLM} />
        </Box>
      </Box>
    </AppTheme>
  );
}
