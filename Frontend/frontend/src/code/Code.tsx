import { useEffect, useState } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import AppTheme from "../shared-ui-theme/AppTheme";
import MonacoEditor from "./components/MonacoEditor";
import CodeSideMenu from "./components/CodeSideMenu";
import AppNavbar from "../dashboard/components/AppNavbar";
import LLMPanel from "./components/LLMPanel";

const defaultPythonCode = `def greet(name):\n    return f"Hello, {name}"\n\nprint(greet("World"))`;

export default function CodeEditorDashboard(props: {
  disableCustomTheme?: boolean;
}) {
  const [code, setCode] = useState(defaultPythonCode);
  const [llmOutput, setLlmOutput] = useState("");

  useEffect(() => {
    document.title = "Zhuhana - Code Editor";
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

        {/* Middle Panel - Code Editor */}
        <Box sx={{ flex: 2, p: 2, display: "flex", flexDirection: "column" }}>
          <Box
            sx={{
              flexGrow: 1,
              border: "1px solid #ccc",
              borderRadius: 2,
              overflow: "hidden",
            }}
          >
            <MonacoEditor code={code} onChange={(v) => setCode(v ?? "")} />
          </Box>
        </Box>

        {/* Right Panel - LLM Panel */}
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
