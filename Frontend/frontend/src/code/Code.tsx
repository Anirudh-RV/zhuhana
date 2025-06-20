import { useEffect, useRef, useState } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
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

  const containerRef = useRef<HTMLDivElement>(null);
  const dragInfo = useRef<{ startY: number; startHeight: number } | null>(null);

  // Default editorHeight = 80% of 100vh - navbar & padding estimate
  const [editorHeight, setEditorHeight] = useState(
    () => window.innerHeight * 0.8
  );

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (!dragInfo.current) return;
      const delta = e.clientY - dragInfo.current.startY;
      const newHeight = dragInfo.current.startHeight + delta;
      setEditorHeight(
        Math.max(100, Math.min(newHeight, window.innerHeight - 100))
      ); // clamp values
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
    setLlmOutput(""); // reset output

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
            onChunk(token);
            setTerminalOutput((prev) => prev + token);
          }
        }
      }
    });
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
            height: "100vh", // ensure it's relative to full height
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

          {/* Terminal Output */}
          <Box
            sx={{
              height: `calc(100% - ${editorHeight + 6 + 16}px)`, // 6px divider + 16px (0.5rem x 2) padding offset
              border: "1px solid #333",
              borderRadius: 1,
              backgroundColor: "#111",
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
          }}
        >
          <LLMPanel output={llmOutput} onSend={handleSendToLLM} />
        </Box>
      </Box>
    </AppTheme>
  );
}
