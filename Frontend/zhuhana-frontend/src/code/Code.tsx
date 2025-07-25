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
import LLMPanel, { LLMPanelHandle } from "./components/LLMPanel";
import Toolbar from "@mui/material/Toolbar";
import Tooltip from "@mui/material/Tooltip";
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
import { CompletionContext, CompletionResult } from "@codemirror/autocomplete";
import { completionKeymap, acceptCompletion } from "@codemirror/autocomplete";
import { indentMore } from "@codemirror/commands";
import EditableFileName, {
  type EditableFileNameHandle,
} from "./components/EditableFileName";
import { keymap } from "@codemirror/view";

import MarkdownIt from "markdown-it";
import DOMPurify from "dompurify";

import { Decoration, ViewPlugin, ViewUpdate } from "@codemirror/view";
import { RangeSetBuilder } from "@codemirror/state";

import { useNavigate } from "react-router-dom";
import ColorModeIconDropdown from "../shared-ui-theme/ColorModeIconDropdown";
import SaveIcon from "@mui/icons-material/Save";

import { initializeLspClient } from "./components/lspClient";
import { useSearchParams } from "react-router-dom";
import {
  USER_PYTHON_ALGORITHM_UPLOAD_V1_ENDPOINT,
  USER_PYTHON_ALGORITHM_EDIT_V1_ENDPOINT,
  USER_PYTHON_ALGORITHMS_INFORMATION_V1_ENDPOINT,
  CREATE_CHAT_SESSION_V1_ENDPOINT,
  USER_PYTHON_ALGORITHM_INFORMATION_V1_ENDPOINT,
  ASK_LLM_V1_ENDPOINT,
  ADD_MESSAGES_V1_ENDPOINT,
} from "../constants";
import { useColorScheme } from "@mui/material/styles";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import { Snackbar, Alert } from "@mui/material";

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

  const { user, accessToken } = useAuth();
  const navigate = useNavigate();
  const [menuAnchorEl, setMenuAnchorEl] = React.useState<null | HTMLElement>(
    null
  );

  const { mode: themeMode, systemMode: themeSystemMode } = useColorScheme();
  const resolvedThemeMode =
    themeMode === "system" ? themeSystemMode : themeMode;

  const [saveStatus, setSaveStatus] = useState<
    "idle" | "saving" | "success" | "error" | "codeRun"
  >("idle");

  const [isNewSession, setIsNewSession] = useState(false);

  const [searchParams, setSearchParams] = useSearchParams();
  const initialAlgorithmId = searchParams.get("algorithm_id");
  const initialSessionId = searchParams.get("session_id");

  const [algorithmId, setAlgorithmId] = useState<string | null>(
    initialAlgorithmId
  );

  const [isSuccessDialogOpen, setSuccessDialogOpen] = useState(false);

  const [sessionId, setSessionId] = useState<string | null>(initialSessionId);

  const fileNameRef = useRef<EditableFileNameHandle>(null);

  useEffect(() => {
    if (!initialAlgorithmId) {
      fileNameRef.current?.focusEditMode();
    }
  }, [initialAlgorithmId]);

  const saveTimeout = useRef<NodeJS.Timeout | null>(null);

  const handleRename = async (newName: string) => {
    setFilename(newName);
    await handleSaveAlgorithm(newName); // Pass the new name to save
  };

  useEffect(() => {
    const fetchAlgorithmDetails = async () => {
      if (!initialAlgorithmId || !accessToken) return;

      try {
        const response = await fetch(
          `${USER_PYTHON_ALGORITHM_INFORMATION_V1_ENDPOINT}?algorithm_id=${initialAlgorithmId}`,
          {
            headers: {
              ...(accessToken ? { USER_TOKEN: accessToken } : {}),
            },
          }
        );

        if (!response.ok) throw new Error("Failed to fetch algorithm details");

        const result = await response.json();
        const scriptName = result?.user_algorithm?.scriptName;

        if (scriptName) {
          setFilename(scriptName);
        }
      } catch (err) {
        console.error("Error fetching algorithm:", err);
      }
    };

    fetchAlgorithmDetails();
  }, [initialAlgorithmId, accessToken]);

  const [showSuccess, setShowSuccess] = useState(false);

  const handleSaveAlgorithm = async (nameOverride?: string) => {
    if (!user || !accessToken) {
      console.error("User not authenticated");
      return;
    }
    setSaveStatus("saving");

    const nameToUse = nameOverride || filename;

    const formData = new FormData();
    formData.append("algorithmName", nameToUse);
    formData.append(
      "algorithm",
      new Blob([code], { type: "text/plain" }),
      `${filename}.py`
    );

    if (algorithmId) {
      formData.append("algorithmID", algorithmId); // Only required for PUT
    }

    // Determine method and endpoint
    const isEdit = !!algorithmId;
    const method = isEdit ? "PUT" : "POST";
    const url = isEdit
      ? USER_PYTHON_ALGORITHM_EDIT_V1_ENDPOINT
      : USER_PYTHON_ALGORITHM_UPLOAD_V1_ENDPOINT;

    try {
      const response = await fetch(url, {
        method,
        headers: {
          ...(accessToken ? { USER_TOKEN: accessToken } : {}),
        },
        body: formData,
      });

      if (!response.ok) throw new Error("Failed to save algorithm");

      const result = await response.json();

      // If it was a new upload, capture the new ID
      if (!algorithmId && result.user_algorithm?.ID) {
        const newId = result.user_algorithm.ID;
        setAlgorithmId(newId);
        searchParams.set("algorithm_id", newId);
        setSearchParams(searchParams);
      }
      setSaveStatus("success");
      setTimeout(() => {
        setSaveStatus("idle");
      }, 5000);
    } catch (err) {
      setSaveStatus("error");
      console.error("❌ Failed to save:", err);
    }
  };

  const [llmPanelWidth, setLlmPanelWidth] = useState(25); // in percentage
  const dragLlmInfo = useRef<{ startX: number; startWidth: number } | null>(
    null
  );

  const [aiPrompt, setAiPrompt] = useState<string>("");

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
  const containerRef = useRef<HTMLDivElement>(null);
  const lspClientRef = useRef<any>(null);
  const [filename, setFilename] = useState("New Algorithm");
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
    });
    lspClientRef.current = client;
  }, []);

  useEffect(() => {
    return () => {
      if (lspClientRef.current?.dispose) {
        lspClientRef.current.dispose();
      }
    };
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

  const completionSource = async (
    context: CompletionContext
  ): Promise<CompletionResult | null> => {
    const offset = context.pos;
    const line = context.state.doc.lineAt(offset);

    const position = {
      line: line.number - 1,
      character: offset - line.from,
    };

    const result = await lspClientRef.current?.completion(position);
    if (!result) return null;

    const items = Array.isArray(result.items) ? result.items : result;

    const before = line.text.slice(0, offset - line.from);
    const match = /[a-zA-Z0-9_]+$/.exec(before);
    const from = match ? offset - match[0].length : offset;

    return {
      from,
      to: offset,
      options: items.map((item: any) => ({
        label: item.label,
        type: "variable", // fallback type; can improve later
        info: item.documentation?.value || item.detail,
        apply: item.insertText || item.label,
      })),
    };
  };

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
      .then(() => {})
      .catch((err) => {
        console.error("❌ Failed to copy:", err);
      });
  };

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

      const latestUserMessage = [...messages]
        .reverse()
        .find((msg) => msg.role === "user");

      if (!latestUserMessage) {
        onChunk("[Error: No user message found]");
        return;
      }

      // Local variable to ensure correct sessionId usage
      let currentSessionId = sessionId;

      if (!currentSessionId) {
        const createSessionResponse = await fetch(
          CREATE_CHAT_SESSION_V1_ENDPOINT,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              ...(accessToken ? { USER_TOKEN: accessToken } : {}),
            },
            body: JSON.stringify({
              algorithm_id: algorithmId,
              title: latestUserMessage.content.slice(0, 60),
            }),
          }
        );

        if (!createSessionResponse.ok) {
          const errorText = await createSessionResponse.text();
          throw new Error(`Failed to create session: ${errorText}`);
        }

        const sessionData = await createSessionResponse.json();
        currentSessionId = sessionData.Result.id;

        setSessionId(currentSessionId);
        setIsNewSession(true);

        searchParams.set("session_id", currentSessionId!);
        setSearchParams(searchParams);
      }

      const queryParams = new URLSearchParams({
        q: prompt,
        current_user_q: latestUserMessage.content,
        ...(currentSessionId ? { session_id: currentSessionId } : {}),
      });

      const response = await fetch(`${ASK_LLM_V1_ENDPOINT}?${queryParams}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          ...(accessToken ? { USER_TOKEN: accessToken } : {}),
        },
        signal,
      });

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

        if (isNewSession) setIsNewSession(false);
      }
    } catch (err: any) {
      if (err.name !== "AbortError") {
        onChunk(`[LLM Error]: ${err.message}`);
      }
    }
  };

  const handleRunCode = async () => {
    if (!pyodideInstanceRef.current) return;

    setRuntimeDiagnostics([]);
    setTerminalOutput([
      { text: ">> Terminal ready...", type: "info" },
      { text: ">> Executing Python code...", type: "info" },
    ]);

    try {
      pyodideInstanceRef.current.FS.writeFile("/main_editor_code.py", code);
      await pyodideInstanceRef.current.runPythonAsync(code);

      setRuntimeDiagnostics([]);
      setTerminalOutput((prev) => [
        ...prev,
        { text: ">> Code execution finished.", type: "success" },
      ]);
      setSaveStatus("codeRun");
      setTimeout(() => {
        setSaveStatus("idle");
      }, 5000);
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

      // ❌ Send error message to LLMPanel
      const formattedPrompt = `While running your code, it encountered the following error:\n\n\`\`\`python\n${message}\n\`\`\`\nDo you want me to help you fix it?\nJust response with a 'Yes' and I'll the fix it for you!`;
      const defaultTitle = "Zhuhana AI code fix";

      try {
        let currentSessionId = sessionId;

        if (!currentSessionId) {
          // Create session
          const createRes = await fetch(CREATE_CHAT_SESSION_V1_ENDPOINT, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              ...(accessToken ? { USER_TOKEN: accessToken } : {}),
            },
            body: JSON.stringify({
              algorithm_id: algorithmId,
              title: defaultTitle,
              system_q: formattedPrompt,
            }),
          });

          const sessionJson = await createRes.json();
          if (sessionJson.Status !== 1)
            throw new Error("Failed to create session");

          currentSessionId = sessionJson.Result.id;
          setSessionId(currentSessionId);
          window.history.replaceState(
            {},
            "",
            `/code/?algorithm_id=${algorithmId}&session_id=${currentSessionId}`
          );
        } else {
          // Send message to session
          await fetch(ADD_MESSAGES_V1_ENDPOINT, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              ...(accessToken ? { USER_TOKEN: accessToken } : {}),
            },
            body: JSON.stringify({
              session_id: currentSessionId,
              system_q: formattedPrompt,
            }),
          });

          setLlmMessages((prev) => [
            ...prev,
            { role: "system", content: formattedPrompt },
          ]);
        }
      } catch (e) {
        console.error("Failed to handle LLM error message:", e);
      }
    }
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
            <IconButton
              size="small"
              onClick={() => navigate("/dashboard")}
              sx={{
                "&:hover": {
                  backgroundColor: "action.hover",
                },
              }}
            >
              <ArrowBackIcon fontSize="small" />
            </IconButton>
          </Box>

          {/* Center: Title */}
          <Box sx={{ flexGrow: 1, display: "flex", justifyContent: "center" }}>
            <Stack
              direction="row"
              alignItems="center"
              spacing={1}
              sx={{
                minWidth: 300, // or whatever minimum width keeps layout stable
              }}
            >
              <EditableFileName
                ref={fileNameRef}
                name={filename}
                onRename={handleRename}
              />
              <Tooltip title="Save File" arrow>
                <IconButton
                  size="small"
                  onClick={() => handleSaveAlgorithm()}
                  sx={{ ml: 1 }}
                >
                  <SaveIcon fontSize="small" />
                </IconButton>
              </Tooltip>
              <Tooltip title="Run the code locally">
                <span>
                  <IconButton
                    onClick={handleRunCode}
                    disabled={isLoadingPyodide}
                    sx={{
                      color:
                        resolvedThemeMode === "dark"
                          ? "#00e676"
                          : "rgba(0, 100, 0, 0.85)",
                      backgroundColor:
                        resolvedThemeMode === "dark"
                          ? "rgba(0, 230, 118, 0.1)"
                          : "rgba(76, 175, 80, 0.15)",
                      border: `1px solid ${
                        resolvedThemeMode === "dark"
                          ? "rgba(0, 230, 118, 0.4)"
                          : "rgba(76, 175, 80, 0.4)"
                      }`,
                      backdropFilter: "blur(4px)",
                      "&:hover": {
                        backgroundColor:
                          resolvedThemeMode === "dark"
                            ? "rgba(0, 230, 118, 0.2)"
                            : "rgba(76, 175, 80, 0.25)",
                      },
                    }}
                    size="small"
                  >
                    <PlayArrowIcon fontSize="small" />
                  </IconButton>
                </span>
              </Tooltip>
              {saveStatus !== "idle" && (
                <Typography
                  variant="body2"
                  noWrap
                  sx={{
                    ml: 1,
                    fontWeight: 500,
                    overflowX: "visible",
                  }}
                >
                  {saveStatus === "saving"
                    ? "Saving..."
                    : saveStatus === "success"
                    ? "Saved!"
                    : saveStatus === "codeRun"
                    ? "✅ Code executed successfully"
                    : "Failed to Save"}
                </Typography>
              )}
            </Stack>
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
                height: "100%",
              }}
            >
              <CodeSideMenu
                onClose={() => setIsSidebarOpen(false)}
                terminalOutput={terminalOutput}
                isLoadingPyodide={isLoadingPyodide}
                onRunCode={handleRunCode}
              />
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
              display: "flex",
              flexDirection: "column",
              minWidth: 0,
            }}
          >
            <Box
              sx={{
                flexGrow: 1,
                border: "1px solid #ccc",
                minHeight: 0,
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
                messages={llmMessages}
                setMessages={setLlmMessages}
                algorithmId={algorithmId}
                sessionId={sessionId}
                setSessionId={setSessionId}
                isNewSession={isNewSession}
                setIsNewSession={setIsNewSession}
              />
            </Box>
          </>
        </Box>
      </Box>
    </AppTheme>
  );
}
