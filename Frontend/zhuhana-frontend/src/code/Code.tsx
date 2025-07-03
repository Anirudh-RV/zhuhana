import { useEffect, useRef, useState, useCallback, useMemo } from "react";
import CssBaseline from "@mui/material/CssBaseline";
import Box from "@mui/material/Box";
import Divider from "@mui/material/Divider";
import AppTheme from "../shared-ui-theme/AppTheme";
import CodeMirrorEditor from "./components/CodeMirrorEditor";
import CodeSideMenu from "./components/CodeSideMenu";
import LLMPanel from "./components/LLMPanel";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import PlayArrowIcon from "@mui/icons-material/PlayArrow";
import { green } from "@mui/material/colors";

// CodeMirror specific imports for LSP integration
import { EditorView, hoverTooltip } from "@codemirror/view";
import type { Extension } from "@codemirror/state";
import { EditorState } from "@codemirror/state";
import { autocompletion } from "@codemirror/autocomplete";
import { linter, Diagnostic as CodeMirrorDiagnostic, lintGutter } from "@codemirror/lint";
import { textDocument } from "codemirror-languageservice";



// Correct imports from codemirror-languageservice based on its API
import {
  createCompletionSource,
  createHoverTooltipSource,
} from "codemirror-languageservice";

import {
  InitializeParams,
  InitializeResult,
  Message,
  RequestMessage,
  PublishDiagnosticsParams,
  CompletionParams,
  HoverParams,
  CompletionList,
  Hover,
  MarkupContent,
  MarkupKind,
  MarkedString,
  TextDocumentIdentifier,
  Position,
  TextDocumentItem,
  TextDocumentContentChangeEvent,
  CompletionContext as LSPCompletionContext,
  TextDocument,
  Diagnostic as LSPDiagnostic,
} from "vscode-languageserver-protocol";

import MarkdownIt from 'markdown-it';
import DOMPurify from 'dompurify';
import ReconnectingWebSocket from 'reconnecting-websocket';

import { TextDecoder } from 'text-encoding'; // Ensure this is imported if needed globally

const md = MarkdownIt();

// Utility to convert LSP MarkupContent/MarkedString to HTML
function lspMarkupContentToHtml(content: string | MarkupContent | MarkedString | MarkedString[] | undefined): string {
  if (!content) return '';

  if (typeof content === 'string') {
    return DOMPurify.sanitize(md.render(content));
  } else if (Array.isArray(content)) {
    return DOMPurify.sanitize(content.map(item => {
      if (typeof item === 'string') {
        return md.render(item);
      } else { // MarkedString { language: string, value: string }
        return item.language ? `<pre><code class="language-${item.language}">${item.value}</code></pre>` : md.render(item.value);
      }
    }).join('\n'));
  } else if ('kind' in content) {
    if (content.kind === MarkupKind.Markdown) {
      return DOMPurify.sanitize(md.render(content.value));
    } else { // MarkupKind.PlainText
      return DOMPurify.sanitize(`<pre>${content.value}</pre>`);
    }
  } else if ('value' in content) {
    if (content.language) {
      return DOMPurify.sanitize(`<pre><code class="language-${content.language}">${content.value}</code></pre>`);
    } else {
      return DOMPurify.sanitize(md.render(content.value));
    }
  }
  return '';
}

// Helper to convert LSP Diagnostic to CodeMirror Diagnostic
function lspDiagnosticToCmDiagnostic(lspDiag: LSPDiagnostic, view: EditorView): CodeMirrorDiagnostic {
  const lineStart = view.state.doc.line(lspDiag.range.start.line + 1);
  const from = lineStart.from + lspDiag.range.start.character;

  const lineEnd = view.state.doc.line(lspDiag.range.end.line + 1);
  const to = lineEnd.from + lspDiag.range.end.character;

  return {
    from: from,
    to: to,
    severity: lspDiag.severity === 1 ? "error" : lspDiag.severity === 2 ? "warning" : "info",
    message: lspDiag.message,
    source: lspDiag.source,
  };
}

// --- Pyodide setup ---
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

type LLMMessage = {
  role: "user" | "assistant" | "system";
  content: string;
};

const defaultPythonCode = `import zhuhana
from zhuhana.types import OrderInstruction, OrderSide, OrderType, OrderMode, OrderTIF, OrderDomain, OHLCData



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
            quantity=100
        )

    def condition_for_buy(self, current_data: OHLCData) -> OrderInstruction:
        return OrderInstruction(
            side=OrderSide.BUY,
            type=OrderType.MARKET,
            mode=OrderMode.INTRADAY,
            tif=OrderTIF.DAY,
            domain=OrderDomain.BACKTEST,
            quantity=100
        )
`;

// LSP Server URL
const LSP_SERVER_URL = "ws://localhost:3001";
const FILE_URI = "file:///main_editor_code.py"; // A fixed URI for your single editor file
const LANGUAGE_ID = "python"; // LSP language ID for Python

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
  const editorViewRef = useRef<EditorView | null>(null);
  const [isLoadingPyodide, setIsLoadingPyodide] = useState(true);
  const [diagnostics, setDiagnostics] = useState<CodeMirrorDiagnostic[]>([]);

  // LSP WebSocket client
  const lspSocketRef = useRef<ReconnectingWebSocket | null>(null);
  const lspRequestIdCounter = useRef(0);
  const lspPendingRequests = useRef(new Map<number, { resolve: Function, reject: Function }>());
  const isLspInitialized = useRef(false);
  const documentVersion = useRef(1);

  const containerRef = useRef<HTMLDivElement>(null);
  const dragInfo = useRef<{ startY: number; startHeight: number } | null>(null);
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



  // --- LSP Communication Helpers ---
  // IMPORTANT: Ensure this useCallback has an empty dependency array for stability
  const sendLSPMessage = useCallback((message: Message) => {
    if (lspSocketRef.current && lspSocketRef.current.readyState === WebSocket.OPEN) {
      try {
        console.log("🔌 Sending LSP request:", message);
        const json = JSON.stringify(message);
        const encoder = new TextEncoder();
        const contentBytes = encoder.encode(json);
        const header = `Content-Length: ${contentBytes.length}\r\n\r\n`;
        const fullBytes = encoder.encode(header + json);
        lspSocketRef.current.send(fullBytes);

      } catch (e) {
        console.error("Failed to send LSP message:", e);
      }
    } else {
      console.warn("LSP WebSocket not open, cannot send message:", message);
    }
  }, []); // <--- Empty dependency array

  const sendLSPRequest = useCallback((method: string, params: any): Promise<any> => {
    return new Promise((resolve, reject) => {
      const id = lspRequestIdCounter.current++;
      lspPendingRequests.current.set(id, { resolve, reject });

      const request: RequestMessage = {
        jsonrpc: "2.0",
        id: id,
        method: method,
        params: params,
      };
      console.log("🔌 Sending LSP request:", request);
      sendLSPMessage(request);
    });
  }, [sendLSPMessage]);

  const markdownToDom = (markdown: string): DocumentFragment => {
  const html = DOMPurify.sanitize(md.render(markdown));
  const fragment = document.createDocumentFragment();
  const wrapper = document.createElement("div");
  wrapper.innerHTML = html;
  fragment.appendChild(wrapper);
  return fragment;
};


  const completionSource = createCompletionSource({
    markdownToDom,
    doComplete: async (textDocument, position, context) => {
      if (!isLspInitialized.current) return null;
      return await sendLSPRequest("textDocument/completion", {
        textDocument: { uri: textDocument.uri },
        position,
        context,
      });
    },
});

const doHover = async (textDocument: TextDocumentIdentifier, position: Position): Promise<Hover | null> => {
  console.log("➡️ Hover triggered at:", position);
  console.log("WebSocket readyState:", lspSocketRef.current?.readyState);
  console.log("isLspInitialized.current:", isLspInitialized.current);

  if (!isLspInitialized.current || lspSocketRef.current?.readyState !== 1) return null;

  try {
    const result = await sendLSPRequest("textDocument/hover", {
      textDocument,
      position,
    });
    return result as Hover;
  } catch (err) {
    console.error("Hover request failed:", err);
    return null;
  }
};



  const hoverSource = createHoverTooltipSource({
    markdownToDom,
    doHover,
  });




  const sendLSPNotification = useCallback((method: string, params: any) => {
  const notification = {
    jsonrpc: "2.0",
    method,        // ✅ required
    params,        // ✅ required (can be null)
  };
  sendLSPMessage(notification);
}, [sendLSPMessage]);



  // --- Pyodide Initialization (runs once on mount) ---
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
            { text: ">> Loading micropip...", type: "info" },
          ]);
          await pyodide.loadPackage("micropip");
          setTerminalOutput((prev) => [
            ...prev,
            { text: `>> Installing zhuhana==0.1.0...`, type: "info" },
          ]);
          await pyodide.runPythonAsync(`
            import micropip
            await micropip.install(["https://test-files.pythonhosted.org/packages/0d/a6/9442d90a4f723583e7d625107129af7c61151f2c4c2749f04791edfeed96/zhuhana-0.1.0-py3-none-any.whl"])
          `);
          setTerminalOutput((prev) => [
            ...prev,
            { text: `>> Package 'zhuhana' installed successfully.`, type: "info" },
          ]);

          setTerminalOutput((prev) => [
            ...prev,
            { text: ">> Python runtime loaded.", type: "info" },
          ]);
        } else {
          throw new Error("window.loadPyodide is not a function.");
        }
      } catch (err: any) {
        setTerminalOutput((prev) => [
          ...prev,
          { text: `>> Failed to initialize Pyodide: ${err.message || err}`, type: "error" },
        ]);
      } finally {
        setIsLoadingPyodide(false);
      }
    };

    // Cleanup function for Pyodide script
    return () => {
      if (document.head.contains(script)) {
        document.head.removeChild(script);
      }
    };
  }, []); // Empty dependency array for single run on mount


  // --- LSP Client Logic (runs once on mount and handles updates via notifications) ---
  // IMPORTANT: Ensure this useEffect has an empty dependency array for single run on mount
  useEffect(() => {
  console.log("Setting up LSP WebSocket...");
  const ws = new ReconnectingWebSocket(LSP_SERVER_URL);
  lspSocketRef.current = ws;

  let messageBuffer = "";

  ws.onopen = async () => {
    console.log("LSP WebSocket connected.");
    setTerminalOutput((prev) => [...prev, { text: ">> LSP connected.", type: "info" }]);

    const initializeParams: InitializeParams = {
      processId: null,
      clientInfo: { name: "CodeMirror React Client", version: "1.0" },
      rootUri: "file:///",
      capabilities: {
        textDocument: {
          completion: {
            completionItem: {
              documentationFormat: [MarkupKind.Markdown, MarkupKind.PlainText],
              snippetSupport: true,
              resolveSupport: { properties: ["documentation", "detail"] },
            },
            contextSupport: true,
          },
          hover: { contentFormat: [MarkupKind.Markdown, MarkupKind.PlainText] },
          synchronization: {
            didSave: true,
            willSave: true,
            willSaveWaitUntil: true,
            dynamicRegistration: true,
          },
          publishDiagnostics: { relatedInformation: true, tagSupport: { valueSet: [1, 2] } },
        },
        workspace: {
          workspaceFolders: true,
          didChangeWatchedFiles: { dynamicRegistration: true },
        },
      },
      workspaceFolders: [{ uri: "file:///", name: "Workspace" }],
    };

    try {
      const response: InitializeResult = await sendLSPRequest("initialize", initializeParams);
      console.log("LSP initialize response:", response);
      isLspInitialized.current = true;
      sendLSPNotification("initialized", {});
      sendLSPNotification("textDocument/didOpen", {
        textDocument: {
          uri: FILE_URI,
          languageId: LANGUAGE_ID,
          version: documentVersion.current,
          text: code,
        },
      });
    } catch (error: any) {
      console.error("LSP initialization failed:", error);
      setTerminalOutput((prev) => [...prev, { text: `>> LSP initialization failed: ${error.message || error}`, type: "error" }]);
    }
  };

  function handleLspMessage(message: any) {
  if ('id' in message && message.id !== undefined && message.id !== null) {
    const pending = lspPendingRequests.current.get(message.id);
    if (pending) {
      lspPendingRequests.current.delete(message.id);
      if ('error' in message) {
        pending.reject(message.error);
      } else {
        pending.resolve(message.result);
      }
    }
  } else if ('method' in message) {
    if (message.method === "textDocument/publishDiagnostics") {
      const params: PublishDiagnosticsParams = message.params;
      if (editorViewRef.current) {
        const cmDiagnostics: CodeMirrorDiagnostic[] = params.diagnostics.map((diag: LSPDiagnostic) =>
          lspDiagnosticToCmDiagnostic(diag, editorViewRef.current!)
        );
        setDiagnostics(cmDiagnostics);
      }
    }
  }
}


  ws.onmessage = async (event) => {
  let chunk: string;

  if (typeof event.data === "string") {
    chunk = event.data;
  } else if (event.data instanceof Blob) {
    chunk = await event.data.text(); // ✅ decode the Blob properly
  } else {
    const decoder = new TextDecoder("utf-8");
    chunk = decoder.decode(event.data); // for ArrayBuffer (rare in browsers)
  }

  console.log("📥 Decoded chunk from WS:", chunk);

  messageBuffer += chunk;

  while (true) {
    const headerEnd = messageBuffer.indexOf("\r\n\r\n");
    if (headerEnd === -1) break;

    const header = messageBuffer.substring(0, headerEnd);
    const contentLengthMatch = header.match(/Content-Length: (\d+)/i);
    if (!contentLengthMatch) {
      console.error("❌ Invalid LSP header:", header);
      break;
    }

    const contentLength = parseInt(contentLengthMatch[1], 10);
    const fullLength = headerEnd + 4 + contentLength;

    if (messageBuffer.length < fullLength) break;

    const jsonPart = messageBuffer.substring(headerEnd + 4, fullLength);
    messageBuffer = messageBuffer.substring(fullLength);

    try {
      const message = JSON.parse(jsonPart);
      handleLspMessage(message);
    } catch (err) {
      console.error("❌ Failed to parse LSP JSON:", err, jsonPart);
    }
  }
};



  ws.onerror = (error) => {
    console.error("LSP WebSocket error:", error);
    setTerminalOutput((prev) => [...prev, { text: `>> LSP connection error. Make sure your Node.js LSP proxy is running on ${LSP_SERVER_URL}.`, type: "error" }]);
  };

  ws.onclose = (event) => {
    console.log("LSP WebSocket closed:", event.code, event.reason);
    setTerminalOutput((prev) => [...prev, { text: `>> LSP disconnected. Code: ${event.code}, Reason: ${event.reason}`, type: "info" }]);
    isLspInitialized.current = false;
    lspPendingRequests.current.forEach(req => req.reject(new Error("LSP connection closed")));
    lspPendingRequests.current.clear();
    setDiagnostics([]);
  };

  return () => {
    console.log("Cleaning up LSP WebSocket...");
    if (lspSocketRef.current) {
      if (isLspInitialized.current) {
        const id = lspRequestIdCounter.current++;
        const shutdownNotification: RequestMessage = { jsonrpc: "2.0", method: "shutdown", id };
        const exitNotification: RequestMessage = { jsonrpc: "2.0", method: "exit", id };
        try {
          lspSocketRef.current.send(JSON.stringify(shutdownNotification));
          lspSocketRef.current.send(JSON.stringify(exitNotification));
        } catch (e) {
          console.warn("Failed to send LSP shutdown/exit on cleanup:", e);
        }
      }
      lspSocketRef.current.close();
    }
  };
}, []);

  // Callback to get the EditorView instance from CodeMirrorEditor
  const handleEditorCreation = useCallback((view: EditorView) => {
    editorViewRef.current = view;
  }, []);

  // Update the code, and notify the LSP server of changes
  const handleCodeChange = useCallback((newCode: string | undefined) => {
    const currentCode = newCode ?? "";
    setCode(currentCode);

    documentVersion.current += 1;

    if (isLspInitialized.current && lspSocketRef.current?.readyState === WebSocket.OPEN) {
        sendLSPNotification("textDocument/didChange", {
            textDocument: { uri: FILE_URI, version: documentVersion.current },
            contentChanges: [{ text: currentCode }],
        });
    }
  }, [sendLSPNotification]);


  // --- Wrapper functions for createCompletionSource and createHoverTooltipSource ---

  const doCompleteLSP = useCallback(async (
    textDocument: TextDocument,
    position: Position,
    context: LSPCompletionContext
  ): Promise<CompletionList | Iterable<CompletionList> | null | undefined> => {
    if (!isLspInitialized.current || lspSocketRef.current?.readyState !== WebSocket.OPEN) {
        return null;
    }
    const params: CompletionParams = {
        textDocument: { uri: textDocument.uri },
        position: position,
        context: context,
    };
    try {
        const result = await sendLSPRequest('textDocument/completion', params);
        return result as CompletionList;
    } catch (error) {
        console.error("LSP Completion request failed:", error);
        return null;
    }
  }, [sendLSPRequest]);


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

    console.log("[Current Code]:\n" + code);
    console.log("[Code Lines]:", code.split("\n").length);

    try {
      pyodideInstanceRef.current.FS.writeFile("/main_editor_code.py", code);

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
    }
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
        const lines = chunk.split("\n").filter((line: string) => line.trim() !== "");

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
    // These listeners are only added when drag starts, so they should be removed then too.
    window.removeEventListener("mousemove", handleMouseMove);
    window.removeEventListener("mouseup", handleMouseUp);
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
            <CodeMirrorEditor
            code={code}
            onChange={handleCodeChange}
            onCreateEditor={handleEditorCreation}
            extraExtensions={[
            textDocument(FILE_URI),
            autocompletion({ override: [completionSource] }),
            hoverTooltip(hoverSource),
            linter(() => diagnostics),
            lintGutter(),
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
              // Add listeners when drag starts
              window.addEventListener("mousemove", handleMouseMove);
              window.addEventListener("mouseup", handleMouseUp);
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
