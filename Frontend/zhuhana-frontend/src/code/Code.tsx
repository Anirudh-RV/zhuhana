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

  // Default editorHeight = 80% of the available height for the middle panel
  // We'll calculate this dynamically now.
  const [editorHeight, setEditorHeight] = useState(0); // Initialize to 0

  useEffect(() => {
    // Calculate initial editor height based on container's available space
    if (containerRef.current) {
      const parentHeight = containerRef.current.offsetHeight;
      // Subtract the fixed height of the divider and padding/margin estimates
      // Adjust these values based on your actual Divider height and Box padding
      const nonEditorHeight = 6 + 16; // Divider height (6px) + terminal box p:1 (8px top + 8px bottom = 16px)
      setEditorHeight((parentHeight - nonEditorHeight) * 0.8);
    }

    const handleMouseMove = (e: MouseEvent) => {
      if (!dragInfo.current) return;
      const delta = e.clientY - dragInfo.current.startY;
      const newHeight = dragInfo.current.startHeight + delta;

      // Ensure newHeight is within reasonable bounds
      if (containerRef.current) {
        const parentHeight = containerRef.current.offsetHeight;
        const nonEditorHeight = 6 + 16; // Divider height + terminal padding
        setEditorHeight(
          Math.max(
            100,
            Math.min(newHeight, parentHeight - nonEditorHeight - 100)
          )
        ); // clamp values, leaving space for terminal
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
            // height property removed - let flexbox handle it
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
              flexGrow: 1, // Let terminal take remaining space
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
            display: "flex", // Keep display flex if LLMPanel itself has flex children
            flexDirection: "column", // Keep column if LLMPanel children are stacked
            borderLeft: "1px solid",
            borderColor: "divider",
            backgroundColor: "background.paper",
            p: 2,
            flexShrink: 0, // Prevent shrinking
            // height property removed - let flexbox handle it
          }}
        >
          <LLMPanel output={llmOutput} onSend={handleSendToLLM} />
        </Box>
      </Box>
    </AppTheme>
  );
}
