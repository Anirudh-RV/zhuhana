import { useEffect, useRef, useState } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import AppBar from "@mui/material/AppBar"; // New import
import Toolbar from "@mui/material/Toolbar"; // New import
import Typography from "@mui/material/Typography"; // New import
import IconButton from "@mui/material/IconButton"; // New import
import PlayArrowIcon from "@mui/icons-material/PlayArrow"; // New import
import AppTheme from "../shared-ui-theme/AppTheme";
import MonacoEditor from "./components/MonacoEditor";
import CodeSideMenu from "./components/CodeSideMenu";
import LLMPanel from "./components/LLMPanel";

const defaultPythonCode = `def greet(name):\n    return f"Hello, {name}"\n\nprint(greet("World"))`;

export default function CodeEditorDashboard(props: {
  disableCustomTheme?: boolean;
}) {
  const [code, setCode] = useState(defaultPythonCode);
  const [llmOutput, setLlmOutput] = useState("");
  const [terminalOutput, setTerminalOutput] = useState(
    ">> Terminal ready...\n"
  );
  const [terminalTitle] = useState("Terminal"); // New state for terminal title

  const containerRef = useRef<HTMLHtmlElement>(null);
  const dragInfo = useRef<{ startY: number; startHeight: number } | null>(null);

  const [editorHeight, setEditorHeight] = useState(0);

  useEffect(() => {
    if (containerRef.current) {
      const parentHeight = containerRef.current.offsetHeight;
      const nonEditorHeight = 6 + 16 + 48; // Divider (6px) + terminal box p:1 (16px) + AppBar (48px)
      setEditorHeight((parentHeight - nonEditorHeight) * 0.8);
    }

    const handleMouseMove = (e: MouseEvent) => {
      if (!dragInfo.current) return;
      const delta = e.clientY - dragInfo.current.startY;
      const newHeight = dragInfo.current.startHeight + delta;

      if (containerRef.current) {
        const parentHeight = containerRef.current.offsetHeight;
        const nonEditorHeight = 6 + 16 + 48; // Divider (6px) + terminal box p:1 (16px) + AppBar (48px)
        setEditorHeight(
          Math.max(
            100,
            Math.min(newHeight, parentHeight - nonEditorHeight - 100)
          )
        );
      }
    };

    const handleMouseUp = () => {
      dragInfo.current = null;
      document.body.style.cursor = "default";
      document.body.style.userSelect = "auto";
    };

    window.addEventListener("mousemove", handleMouseMove);
    window.addEventListener("mouseup", handleMouseUp);
    return () => {
      window.removeEventListener("mousemove", handleMouseMove);
      window.removeEventListener("mouseup", handleMouseUp);
    };
  }, []);

  const handleSendToLLM = (input: string, onChunk: (token: string) => void) => {
    setLlmOutput("");

    fetch("http://localhost:8002/stream-llm", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ input }),
    }).then(async (res) => {
      const reader = res.body?.getReader();
      const decoder = new TextDecoder();

      while (reader) {
        const { done, value } = await reader.read();
        if (done) break;

        const chunk = decoder.decode(value);
        for (const line of chunk.split("\n")) {
          if (line.startsWith("data: ")) {
            const token = line.replace("data: ", "");
            setLlmOutput((prev) => prev + token);
            setTerminalOutput((prev) => prev + token);
            onChunk(token);
          }
        }
      }
    });
  };

  // New function to handle running the script
  const handleRunScript = () => {
    setTerminalOutput(
      (prev) => prev + `\n>> Running script...\n${code}\n>> Script finished.\n`
    );
    // In a real application, you would execute the 'code' here.
    // For example, sending it to a backend endpoint that runs Python.
  };

  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: "flex", height: "100vh" }}>
        {/* Left Panel - Sidebar */}
        <CodeSideMenu />

        {/* Middle Panel - Editor + Terminal */}
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
          {/* Code Editor */}
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

          {/* Draggable Divider */}
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

          {/* Terminal Output with Top Bar */}
          <Box
            sx={{
              flexGrow: 1,
              border: "1px solid #333",
              borderRadius: 1,
              backgroundColor: "#111",
              overflow: "hidden", // Important for containing AppBar and content
              display: "flex",
              flexDirection: "column",
            }}
          >
            {/* Top Bar for Terminal */}
            <AppBar
              position="static"
              sx={{
                bgcolor: "#252526", // VS Code-like dark background
                borderBottom: "1px solid #333",
                boxShadow: "none", // Remove shadow
                minHeight: "48px", // Standard App Bar height
              }}
            >
              <Toolbar variant="dense" sx={{ minHeight: "48px" }}>
                <Typography
                  variant="subtitle1"
                  component="div"
                  sx={{ flexGrow: 1, color: "#cccccc" }} // VS Code-like text color
                >
                  {terminalTitle}
                </Typography>
                <IconButton
                  aria-label="run script"
                  onClick={handleRunScript}
                  sx={{ color: "#00ff00" }} // Green arrow color
                >
                  <PlayArrowIcon />
                </IconButton>
              </Toolbar>
            </AppBar>
            {/* Terminal Content */}
            <Box
              sx={{
                flexGrow: 1, // Let terminal content take remaining space
                color: "#0f0",
                fontFamily: "monospace",
                fontSize: "0.875rem",
                p: 1,
                overflowY: "auto",
                whiteSpace: "pre-wrap",
              }}
            >
              {terminalOutput}
            </Box>
          </Box>
        </Box>

        {/* Right Panel - LLM */}
        <Box
          sx={{
            width: "25%",
            display: "flex",
            flexDirection: "column",
            borderLeft: "1px solid",
            borderColor: "divider",
            backgroundColor: "background.paper",
            p: 2,
            flexShrink: 0,
          }}
        >
          <LLMPanel output={llmOutput} onSend={handleSendToLLM} />
        </Box>
      </Box>
    </AppTheme>
  );
}
