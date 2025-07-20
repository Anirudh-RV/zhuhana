// Updated CodeEditorDashboard.tsx with LSP logic extracted to lspClient.ts
import { useEffect, useRef, useState, useCallback } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import * as React from "react";
import TerminalPanel, { TerminalLine } from "./components/TerminalPanel";
import Avatar from "@mui/material/Avatar";
import MenuIcon from "@mui/icons-material/Menu";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import AppTheme from "../shared-ui-theme/AppTheme";
import CodeMirrorEditor from "./components/CodeMirrorEditor";
import CodeSideMenu from "./components/CodeSideMenu";
import LLMPanel from "./components/LLMPanel";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import { useAuth } from "../AuthContext";
import OptionsMenu from "../dashboard/components/OptionsMenu";
import Stack from "@mui/material/Stack";
import { EditorView, hoverTooltip } from "@codemirror/view";
import { autocompletion } from "@codemirror/autocomplete";
import {
  linter,
  lintGutter,
  Diagnostic as CodeMirrorDiagnostic,
} from "@codemirror/lint";
import { textDocument } from "codemirror-languageservice";
import {
  createCompletionSource,
  createHoverTooltipSource,
} from "codemirror-languageservice";

import EditableFileName from "./components/EditableFileName";

import MarkdownIt from "markdown-it";
import DOMPurify from "dompurify";

import { Decoration, ViewPlugin, ViewUpdate } from "@codemirror/view";
import { RangeSetBuilder } from "@codemirror/state";

import ContentCopyIcon from "@mui/icons-material/ContentCopy";
import ColorModeIconDropdown from "../shared-ui-theme/ColorModeIconDropdown";

import { initializeLspClient } from "./components/lspClient";

const md = new MarkdownIt();
const FILE_URI =
  "file:///Users/anirudhrv/Desktop/zhuana-trading/Frontend/lsp-server/main_editor_code.py";
const LANGUAGE_ID = "python";
type Message = {
  role: "user" | "assistant" | "system";
  content: string;
};

const markdownToDom = (markdown: string): DocumentFragment => {
  const html = DOMPurify.sanitize(md.render(markdown));
  const fragment = document.createDocumentFragment();
  const wrapper = document.createElement("div");
  wrapper.innerHTML = html;
  fragment.appendChild(wrapper);
  return fragment;
};

function highlightErrorLines(diagnostics: CodeMirrorDiagnostic[]) {
  const plugin = ViewPlugin.fromClass(
    class {
      decorations;

      constructor(view: EditorView) {
        this.decorations = this.buildDecorations(view);
      }

      update(update: ViewUpdate) {
        if (update.docChanged || update.viewportChanged) {
          this.decorations = this.buildDecorations(update.view);
        }
      }

      buildDecorations(view: EditorView) {
        const builder = new RangeSetBuilder<Decoration>();

        // Sort diagnostics by position before adding
        for (const d of diagnostics.sort((a, b) => a.from - b.from)) {
          const line = view.state.doc.lineAt(d.from);
          builder.add(
            line.from,
            line.from,
            Decoration.line({
              attributes: { class: "cm-error-line" },
            })
          );
        }

        return builder.finish();
      }

      destroy() {}
    },
    {
      decorations: (v) => v.decorations,
    }
  );

  return [plugin]; // Extension[]
}

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

type LLMMessage = { role: "user" | "assistant" | "system"; content: string };
const defaultPythonCode = `import zhuhana
from zhuhana.types import (
    OHLCData,
    OrderDomain,
    OrderInstruction,
    OrderMode,
    OrderSide,
    OrderTIF,
    OrderType,
)


class ZhuhanaStrategy:
    def __init__(self, zhuhana_sdk: zhuhana.ZhuhanaClass):
      self.zhuhana_sdk: zhuhana.ZhuhanaClass = zhuhana_sdk

    def on_data(self, current_data: OHLCData):
      pass

    def condition_for_sell(self, current_data: OHLCData) -> OrderInstruction:
      return OrderInstruction(
            side=OrderSide.SELL,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=100,
        )

    def condition_for_buy(self, current_data: OHLCData) -> OrderInstruction:
      return OrderInstruction(
            side=OrderSide.BUY,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=100,
        )`;

export default function CodeEditorDashboard(props: {
  disableCustomTheme?: boolean;
}) {
  useEffect(() => {
    document.title = "Zhuhana - Algorithm IDE";
  }, []);

  const { user } = useAuth();
  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(
    null
  );

  const [llmPanelWidth, setLlmPanelWidth] = useState(25); // in percentage
  const dragLlmInfo = useRef<{ startX: number; startWidth: number } | null>(
    null
  );

  const handleLlmMouseMove = (e: MouseEvent) => {
    if (!dragLlmInfo.current) return;
    const delta = dragLlmInfo.current.startX - e.clientX;
    const newWidth = Math.min(
      60,
      Math.max(
        10,
        dragLlmInfo.current.startWidth + (delta / window.innerWidth) * 100
      )
    );
    setLlmPanelWidth(newWidth);
  };

  const handleLlmMouseUp = () => {
    dragLlmInfo.current = null;
    document.body.style.cursor = "default";
    document.body.style.userSelect = "auto";
    window.removeEventListener("mousemove", handleLlmMouseMove);
    window.removeEventListener("mouseup", handleLlmMouseUp);
  };

  const handleAvatarClick = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };

  const [code, setCode] = useState(defaultPythonCode);
  const [terminalOutput, setTerminalOutput] = useState<TerminalLine[]>([
    { text: ">> Terminal ready...", type: "info" },
  ]);

  const pyodideInstanceRef = useRef<any>(null);
  const editorViewRef = useRef<EditorView | null>(null);
  const [isLoadingPyodide, setIsLoadingPyodide] = useState(true);
  const [diagnostics, setDiagnostics] = useState<CodeMirrorDiagnostic[]>([]);
  const containerRef = useRef<HTMLDivElement>(null);
  const dragInfo = useRef<{ startY: number; startHeight: number } | null>(null);
  const [editorHeight, setEditorHeight] = useState(
    () => window.innerHeight * 0.82
  );
  const lspClientRef = useRef<any>(null);
  const [filename, setFilename] = useState("NewAlgorithm");
  const [runtimeDiagnostics, setRuntimeDiagnostics] = useState<
    CodeMirrorDiagnostic[]
  >([]);
  const [lspDiagnostics, setLspDiagnostics] = useState<CodeMirrorDiagnostic[]>(
    []
  );
  const [llmMessages, setLlmMessages] = useState<Message[]>([]);

  const appendStdout = useRef((msg: string) => {
    setTerminalOutput((prev) => [...prev, { text: msg, type: "success" }]);
  });

  const appendStderr = useRef((msg: string) => {
    console.error("[Pyodide stderr]", msg); // Debug log
    setTerminalOutput((prev) => [
      ...prev,
      { text: `[Python Error]: ${msg}`, type: "error" },
    ]);
  });

  useEffect(() => {
    const PYODIDE_BASE_URL = "https://cdn.jsdelivr.net/pyodide/v0.26.1/full/";
    const script = document.createElement("script");
    script.src = PYODIDE_BASE_URL + "pyodide.js";
    script.async = true;
    script.setAttribute("data-pyodide-base-url", PYODIDE_BASE_URL);
    document.head.appendChild(script);

    script.onload = async () => {
      try {
        const pyodide = await window.loadPyodide({
          indexURL: PYODIDE_BASE_URL,
          stdout: appendStdout.current,
          stderr: appendStderr.current,
        });
        pyodideInstanceRef.current = pyodide;

        await pyodide.loadPackage("micropip");
        await pyodide.runPythonAsync(`
          import micropip
          await micropip.install(["https://test-files.pythonhosted.org/packages/0d/a6/9442d90a4f723583e7d625107129af7c61151f2c4c2749f04791edfeed96/zhuhana-0.1.0-py3-none-any.whl"])
        `);
      } catch (err) {
        appendStderr.current("Failed to initialize Pyodide");
      } finally {
        setIsLoadingPyodide(false);
      }
    };
  }, []);

  useEffect(() => {
    const client = initializeLspClient({
      uri: FILE_URI,
      languageId: LANGUAGE_ID,
      code,
      onDiagnostics: setLspDiagnostics,
      getEditorView: () => editorViewRef.current,
      onInitialized: () => console.log("✅ LSP ready"),
    });
    lspClientRef.current = client;
  }, []);

  const handleEditorCreation = useCallback((view: EditorView) => {
    editorViewRef.current = view;
  }, []);

  const handleCodeChange = useCallback((newCode: string | undefined) => {
    const currentCode = newCode ?? "";
    setCode(currentCode);
    setRuntimeDiagnostics([]); // ✅ clear runtime errors on edit
    lspClientRef.current?.sendDidChange(currentCode);
  }, []);

  const completionSource = createCompletionSource({
    markdownToDom,
    doComplete: async (_doc, position, context) => {
      return await lspClientRef.current?.completion(position, context);
    },
  });

  const hoverSource = createHoverTooltipSource({
    markdownToDom,
    doHover: async (_doc, position) => {
      return await lspClientRef.current?.hover(position);
    },
  });

  const handleCopyTerminal = () => {
    const filtered = terminalOutput
      .filter((line) => line.type === "success" || line.type === "error")
      .map((line) => line.text)
      .join("\n");

    navigator.clipboard
      .writeText(filtered)
      .then(() => {
        console.log("✅ Filtered terminal output copied to clipboard");
      })
      .catch((err) => {
        console.error("❌ Failed to copy:", err);
      });
  };

  const [isLLMOpen, setIsLLMOpen] = useState(true);

  const handleSendToLLM = async (
    messages: LLMMessage[],
    onChunk: (token: string) => void,
    signal: AbortSignal
  ) => {
    try {
      const prompt = messages
        .map(
          (msg) =>
            `${msg.role === "user" ? "User" : "Assistant"}: ${msg.content}`
        )
        .join("\n");

      const sessionId = "e3268d6d-c776-42d8-a04b-c45dc87dc86b";
      const userToken =
        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTQ5ODUxNzAsImlhdCI6MTc1MjM5MzE3MCwidXNlcl9pZCI6IjdlNWEwOTQ0LTM3ODAtNGI3Ni04NWE1LWQ3ZmY1YWUyYTJhOSIsInVzZXJfdHlwZSI6InVzZXIifQ.nkq9c0rGJv2H92MyJc93vo05-XPLTYabVEXTRrTeJb0";

      const response = await fetch(
        `http://localhost:3000/v1/ask/?q=${encodeURIComponent(
          prompt
        )}&session_id=${sessionId}`,
        {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
            USER_TOKEN: userToken,
          },
          signal, // if you have an AbortSignal for cancellation
        }
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
        const lines = chunk
          .split("\n")
          .filter((line: string) => line.trim() !== "");

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

  const handleRunCode = async () => {
    if (!pyodideInstanceRef.current) return;

    // Clear previous runtime errors
    setRuntimeDiagnostics([]);

    setTerminalOutput([
      { text: ">> Terminal ready...", type: "info" },
      { text: ">> Executing Python code...", type: "info" },
    ]);

    try {
      pyodideInstanceRef.current.FS.writeFile("/main_editor_code.py", code);
      await pyodideInstanceRef.current.runPythonAsync(code);

      setRuntimeDiagnostics([]); // ✅ no error
      setTerminalOutput((prev) => [
        ...prev,
        { text: ">> Code execution finished.", type: "success" },
      ]);
    } catch (error: any) {
      const message = error.message || String(error);
      const execLineMatch = message.match(/File "<exec>", line (\d+)/);
      const line = execLineMatch
        ? parseInt(execLineMatch[1], 10) - 1
        : undefined;

      if (line !== undefined && editorViewRef.current) {
        const from = editorViewRef.current.state.doc.line(line + 1).from;
        const to = editorViewRef.current.state.doc.line(line + 1).to;

        setRuntimeDiagnostics([
          {
            from,
            to,
            severity: "error",
            message,
            source: "Pyodide Runtime",
          },
        ]);
      }

      setTerminalOutput((prev) => [
        ...prev,
        { text: `>> Error: ${message}`, type: "error" },
      ]);
    }
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
    window.removeEventListener("mousemove", handleMouseMove);
    window.removeEventListener("mouseup", handleMouseUp);
  };

  const [isSidebarOpen, setIsSidebarOpen] = useState(true);
  const collapseBreakpoint = 900; // px
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth < collapseBreakpoint) {
        setIsSidebarOpen(false);
      } else {
        setIsSidebarOpen(true);
      }
    };

    // Run on mount
    handleResize();

    // Add listener
    window.addEventListener("resize", handleResize);

    return () => window.removeEventListener("resize", handleResize);
  }, []);

  return (
    <AppTheme {...props}>
      <CssBaseline enableColorScheme />
      <Box sx={{ display: "flex", flexDirection: "column", height: "100vh" }}>
        {/* Top Toolbar */}
        <Toolbar
          variant="dense"
          sx={{
            backgroundColor: "background.paper",
            borderBottom: "1px solid",
            borderColor: "divider",
            px: 2,
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between",
          }}
        >
          {/* Left: Back button */}
          <Box sx={{ display: "flex", alignItems: "center", minWidth: "60px" }}>
            <IconButton size="small" onClick={() => console.log("Go Back")}>
              <ArrowBackIcon fontSize="small" />
            </IconButton>
          </Box>

          {/* Center: Title */}
          <Box sx={{ flexGrow: 1, display: "flex", justifyContent: "center" }}>
            <Typography variant="subtitle1" fontWeight="bold" component="div">
              <EditableFileName name={filename} onRename={setFilename} />
            </Typography>
          </Box>

          {/* Right: Empty space to balance layout */}
          <Stack direction="row" sx={{ gap: 1 }}>
            <ColorModeIconDropdown />
            <Avatar
              alt={user?.FirstName}
              src="/static/images/avatar/7.jpg"
              sx={{ width: 36, height: 36, cursor: "pointer" }}
              onClick={handleAvatarClick}
            />

            <OptionsMenu anchorEl={menuAnchorEl} onClose={handleMenuClose} />
          </Stack>
        </Toolbar>

        {/* Three-Panel Layout */}
        <Box sx={{ display: "flex", flexGrow: 1, minHeight: 0 }}>
          {/* Left Sidebar */}
          {isSidebarOpen ? (
            <Box
              sx={{
                width: "20%",
                transition: "width 0.3s ease",
                overflow: "hidden",
                backgroundColor: "background.paper",
                borderRight: "1px solid",
                borderColor: "divider",
                display: "flex",
                flexDirection: "column",
              }}
            >
              <CodeSideMenu onClose={() => setIsSidebarOpen(false)} />
            </Box>
          ) : (
            <Box
              sx={{
                width: "60px",
                display: "flex",
                flexDirection: "column",
                alignItems: "center",
                pt: 1,
                gap: 1,
                borderRight: "1px solid",
                borderColor: "divider",
                backgroundColor: "background.paper",
              }}
            >
              <IconButton size="small" onClick={() => setIsSidebarOpen(true)}>
                <MenuIcon fontSize="small" />
              </IconButton>
            </Box>
          )}

          {/* Center Code + Terminal */}
          <Box
            ref={containerRef}
            sx={{
              flexGrow: 1,
              px: 0.5,
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
              <CodeMirrorEditor
                code={code}
                onChange={handleCodeChange}
                onCreateEditor={handleEditorCreation}
                extraExtensions={[
                  textDocument(FILE_URI),
                  autocompletion({ override: [completionSource] }),
                  hoverTooltip(hoverSource),
                  linter(() => [...lspDiagnostics, ...runtimeDiagnostics]),
                  lintGutter(),
                  ...(runtimeDiagnostics.length > 0
                    ? highlightErrorLines(runtimeDiagnostics)
                    : []),
                ]}
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
                window.addEventListener("mousemove", handleMouseMove);
                window.addEventListener("mouseup", handleMouseUp);
              }}
            />

            <TerminalPanel
              terminalOutput={terminalOutput}
              isLoadingPyodide={isLoadingPyodide}
              onRunCode={handleRunCode}
              onCopyTerminal={handleCopyTerminal}
            />
          </Box>

          {/* Right LLM Panel */}
          <>
            <Divider
              sx={{
                width: "6px",
                cursor: "col-resize",
                backgroundColor: "divider",
              }}
              onMouseDown={(e) => {
                dragLlmInfo.current = {
                  startX: e.clientX,
                  startWidth: llmPanelWidth,
                };
                document.body.style.cursor = "col-resize";
                document.body.style.userSelect = "none";
                window.addEventListener("mousemove", handleLlmMouseMove);
                window.addEventListener("mouseup", handleLlmMouseUp);
              }}
            />
            <Box
              sx={{
                width: `${llmPanelWidth}%`,
                minWidth: "280px",
                display: "flex",
                flexDirection: "column",
                borderLeft: "1px solid",
                borderColor: "divider",
                backgroundColor: "background.paper",
              }}
            >
              <LLMPanel
                onSend={handleSendToLLM}
                onClose={() => setIsLLMOpen(false)}
                messages={llmMessages}
                setMessages={setLlmMessages}
              />
            </Box>
          </>
        </Box>
      </Box>
    </AppTheme>
  );
}
