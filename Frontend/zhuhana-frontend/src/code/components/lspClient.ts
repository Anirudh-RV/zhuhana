import { MonacoLanguageClient } from "monaco-languageclient";

import { MessageTransports } from "vscode-languageclient";
import type { ErrorHandlerResult } from "vscode-languageclient";
import type { CloseHandlerResult } from "vscode-languageclient";

import {
  toSocket,
  WebSocketMessageReader,
  WebSocketMessageWriter,
} from "vscode-ws-jsonrpc";

import normalizeUrl from "normalize-url";
import ReconnectingWebSocket from "reconnecting-websocket";

export function createPythonLanguageClient(): MonacoLanguageClient {
  const url = normalizeUrl("ws://localhost:3001");
  const webSocket = new ReconnectingWebSocket(url, [], {
    maxRetries: 10,
    connectionTimeout: 10000,
  });

  const socket = toSocket(webSocket as WebSocket);
  const reader = new WebSocketMessageReader(socket);
  const writer = new WebSocketMessageWriter(socket);

  const messageTransports: MessageTransports = { reader, writer };

  const languageClient = new MonacoLanguageClient({
    id: "python-client",
    name: "Python LSP",
    messageTransports,
    clientOptions: {
      documentSelector: ["python"],
      errorHandler: {
        error: (): ErrorHandlerResult => ({
          action: 1,
          message: "Continuing after error",
          handled: true,
        }),
        closed: (): CloseHandlerResult => ({
          action: 1,
        }),
      },
    },
  });

  return languageClient;
}
