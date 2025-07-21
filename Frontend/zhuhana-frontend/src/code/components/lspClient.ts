import {
  Message,
  RequestMessage,
  InitializeParams,
  InitializeResult,
  PublishDiagnosticsParams,
  MarkupKind,
  Diagnostic as LSPDiagnostic,
  TextDocumentIdentifier,
  Position,
  Hover,
} from "vscode-languageserver-protocol";

import {
  linter,
  Diagnostic as CodeMirrorDiagnostic,
  lintGutter,
} from "@codemirror/lint";
import ReconnectingWebSocket from "reconnecting-websocket";
import { TextDecoder } from "text-encoding";
import { EditorView } from "@codemirror/view";

type LspClientOptions = {
  uri: string;
  languageId: string;
  code: string;
  onDiagnostics: (diags: CodeMirrorDiagnostic[]) => void;
  getEditorView: () => EditorView | null;
  onInitialized?: () => void;
};

export function initializeLspClient(options: LspClientOptions) {
  const { uri, languageId, code, onDiagnostics, getEditorView, onInitialized } =
    options;

  const socket = new ReconnectingWebSocket("ws://localhost:3001");
  const pendingRequests = new Map<
    number,
    { resolve: Function; reject: Function }
  >();
  let requestId = 1;
  let version = 1;
  let initialized = false;

  const sendMessage = (msg: any) => {
    const json = JSON.stringify(msg);
    const encoder = new TextEncoder();
    const bytes = encoder.encode(json);
    const header = `Content-Length: ${bytes.length}\r\n\r\n`;
    socket.send(header + json);
  };

  const sendRequest = (method: string, params: any): Promise<any> => {
    const id = requestId++;
    const request: RequestMessage = { jsonrpc: "2.0", id, method, params };
    sendMessage(request);
    return new Promise((resolve, reject) => {
      pendingRequests.set(id, { resolve, reject });
    });
  };

  const sendNotification = (method: string, params: any) => {
    sendMessage({ jsonrpc: "2.0", method, params });
  };

  const sendDidChange = (newCode: string) => {
    version += 1;
    sendNotification("textDocument/didChange", {
      textDocument: { uri, version },
      contentChanges: [{ text: newCode }],
    });
  };

  const hover = async (position: Position): Promise<Hover | null> => {
    if (!initialized) return null;
    return await sendRequest("textDocument/hover", {
      textDocument: { uri },
      position,
    });
  };

  const completion = async (position: Position, context?: any) => {
    if (!initialized) return null;
    return await sendRequest("textDocument/completion", {
      textDocument: { uri },
      position,
      context,
    });
  };

  const lspDiagnosticToCm = (
    diag: LSPDiagnostic
  ): CodeMirrorDiagnostic | null => {
    const view = getEditorView();
    if (!view) return null;
    const lineStart = view.state.doc.line(diag.range.start.line + 1);
    const lineEnd = view.state.doc.line(diag.range.end.line + 1);
    return {
      from: lineStart.from + diag.range.start.character,
      to: lineEnd.from + diag.range.end.character,
      message: diag.message,
      severity: diag.severity === 1 ? "error" : "warning",
      source: diag.source,
    };
  };

  socket.onopen = async () => {
    const initParams: InitializeParams = {
      processId: null,
      clientInfo: { name: "Zhuhana IDE", version: "1.0" },
      rootUri: "file:///",
      capabilities: {},
    };

    const result: InitializeResult = await sendRequest(
      "initialize",
      initParams
    );
    initialized = true;
    sendNotification("initialized", {});
    sendNotification("textDocument/didOpen", {
      textDocument: { uri, languageId, version, text: code },
    });
    onInitialized?.();
  };

  let messageBuffer = new Uint8Array();

  socket.onmessage = async (event) => {
    const decoder = new TextDecoder("utf-8");

    // Convert Blob to Uint8Array if needed
    const newChunk =
      event.data instanceof Blob
        ? new Uint8Array(await event.data.arrayBuffer())
        : new Uint8Array(event.data);

    // Append to existing buffer
    const combined = new Uint8Array(messageBuffer.length + newChunk.length);
    combined.set(messageBuffer);
    combined.set(newChunk, messageBuffer.length);
    messageBuffer = combined;

    while (true) {
      const text = decoder.decode(messageBuffer);
      const headerEnd = text.indexOf("\r\n\r\n");

      if (headerEnd === -1) {
        break;
      }

      const headerText = text.slice(0, headerEnd);
      const match = headerText.match(/Content-Length: (\d+)/i);

      if (!match) {
        break;
      }

      const contentLength = parseInt(match[1], 10);
      const bodyStart = headerEnd + 4;
      const fullMessageLength = bodyStart + contentLength;

      if (messageBuffer.length < fullMessageLength) {
        break;
      }

      const bodyBytes = messageBuffer.slice(bodyStart, fullMessageLength);
      const bodyText = decoder.decode(bodyBytes);

      try {
        const msg = JSON.parse(bodyText);

        // Handle LSP message
        if (msg.id && pendingRequests.has(msg.id)) {
          const { resolve, reject } = pendingRequests.get(msg.id)!;
          pendingRequests.delete(msg.id);
          msg.error ? reject(msg.error) : resolve(msg.result);
        } else if (msg.method === "textDocument/publishDiagnostics") {
          const params: PublishDiagnosticsParams = msg.params;
          const editor = getEditorView();
          if (editor) {
            const diags = params.diagnostics
              .map((d) => lspDiagnosticToCm(d))
              .filter(Boolean) as CodeMirrorDiagnostic[];
            onDiagnostics(diags);
          }
        }
      } catch (err) {
        console.error("❌ Failed to parse LSP JSON body:", err);
      }

      // Trim processed message from buffer
      messageBuffer = messageBuffer.slice(fullMessageLength);
    }
  };

  return {
    sendRequest,
    sendNotification,
    sendDidChange,
    hover,
    completion,
  };
}
