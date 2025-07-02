const { WebSocketServer, WebSocket } = require("ws");
const { createWebSocketStream } = require("ws");
const {
  createMessageConnection,
  StreamMessageReader,
  StreamMessageWriter,
  ConnectionError,
  ConnectionErrors,
} = require("vscode-jsonrpc");
const { TextDecoder } = require("util");
const { spawn } = require("child_process");

const wss = new WebSocketServer({ port: 3001 });
console.log("LSP WebSocket server listening on ws://localhost:3001");

let connectionCounter = 0;

wss.on("connection", (socket) => {
  connectionCounter++;
  const currentConnectionId = connectionCounter;
  console.log(`[Server] New Client connected. ID: ${currentConnectionId}`);
  console.log(
    `[Server] Attempting to spawn pyright-langserver process. ID: ${currentConnectionId}`
  );

  const pyright = spawn("npx", ["pyright-langserver", "--stdio"]);
  console.log(
    `[Server] Pyright process spawned (PID: ${pyright.pid}). ID: ${currentConnectionId}`
  );

  pyright.stdout.on("data", (data) => {
    // console.log(`[Pyright RAW STDOUT - ID: ${currentConnectionId}]: ${data.toString()}`);
  });

  pyright.stderr.on("data", (data) => {
    const stderrMessage = new TextDecoder().decode(data);
    console.error(
      `[Pyright STDERR - ID: ${currentConnectionId}]: ${stderrMessage}`
    );
  });

  pyright.on("error", (err) => {
    console.error(
      `[Server] Failed to start/run Pyright process: ${err.message} (ID: ${currentConnectionId})`
    );
    if (socket.readyState === WebSocket.OPEN) {
      socket.send(
        JSON.stringify({
          jsonrpc: "2.0",
          method: "window/showMessage",
          params: {
            type: 1,
            message: `LSP Server Error: Pyright process failed: ${err.message}`,
          },
        })
      );
    }
    console.log(
      `[Server] Closing socket due to Pyright process error. ID: ${currentConnectionId}`
    );
    socket.close();
  });

  pyright.on("exit", (code, signal) => {
    console.log(
      `[Server] Pyright process exited. Code: ${code}, Signal: ${signal} (ID: ${currentConnectionId})`
    );
    if (socket.readyState === WebSocket.OPEN) {
      console.log(
        `[Server] Closing socket due to Pyright process exit. ID: ${currentConnectionId}`
      );
      socket.close();
    }
  });

  // --- WebSocket Stream for Client ---
  const clientSocketStream = createWebSocketStream(socket);

  // --- LSP Message Connection for the Client WebSocket ---
  const clientStreamReader = new StreamMessageReader(clientSocketStream);
  const clientStreamWriter = new StreamMessageWriter(clientSocketStream);

  const clientLspConnection = createMessageConnection(
    clientStreamReader,
    clientStreamWriter
  );

  // --- LSP Message Reader/Writer for Pyright's stdio ---
  const pyrightReader = new StreamMessageReader(pyright.stdout);
  const pyrightWriter = new StreamMessageWriter(pyright.stdin);

  // --- Bridge Logic ---

  // 1. Messages from Pyright (stdio) -> Client (WebSocket)
  pyrightReader.listen((message) => {
    clientLspConnection.send(message);
    if (message && typeof message === "object") {
      const msgType =
        "method" in message
          ? `Method: ${message.method}`
          : "id" in message
          ? `Response to ID: ${message.id}`
          : "Unknown";
      console.log(
        `[Server] Pyright -> Client (LSP Message - ${msgType} - ID: ${currentConnectionId}): ${JSON.stringify(
          message
        )}`
      );
    } else {
      console.log(
        `[Server] Pyright -> Client (Raw Message - ID: ${currentConnectionId}): ${JSON.stringify(
          message
        )}`
      );
    }
  });
  console.log(
    `[Server] pyrightReader.listen() initialized. ID: ${currentConnectionId}`
  );

  // 2. Messages from Client (WebSocket) -> Pyright (stdio)
  clientLspConnection.listen((message) => {
    // This .listen() call is essential for client-to-server messages
    pyrightWriter.write(message);
    if (message && typeof message === "object") {
      const msgType =
        "method" in message
          ? `Method: ${message.method}`
          : "id" in message
          ? `Request ID: ${message.id}`
          : "Unknown";
      console.log(
        `[Server] Client -> Pyright (LSP Message - ${msgType} - ID: ${currentConnectionId}): ${JSON.stringify(
          message
        )}`
      );
    } else {
      console.log(
        `[Server] Client -> Pyright (Raw Message - ID: ${currentConnectionId}): ${JSON.stringify(
          message
        )}`
      );
    }
  });
  console.log(
    `[Server] clientLspConnection.listen() initialized. ID: ${currentConnectionId}`
  );

  // --- REMOVED THE EXPLICIT clientLspConnection.listen() CALL HERE ---
  // The two .listen() calls above (pyrightReader.listen and clientLspConnection.listen)
  // are often sufficient to start the message processing for both directions.
  // The "Connection is already listening" error suggests this explicit call is redundant or problematic.

  socket.on("close", (code, reason) => {
    console.log(
      `[Server] Client socket closed. Code: ${code}, Reason: ${
        reason ? reason.toString() : "N/A"
      } (ID: ${currentConnectionId}).`
    );
    console.log(
      `[Server] Disposing LSP connections and killing Pyright process. ID: ${currentConnectionId}`
    );
    clientLspConnection.dispose();
    clientStreamReader.dispose();
    clientStreamWriter.dispose();
    pyrightReader.dispose();
    pyrightWriter.dispose();
    if (pyright.pid && !pyright.killed) {
      pyright.kill();
      console.log(
        `[Server] Pyright process (PID: ${pyright.pid}) killed. ID: ${currentConnectionId}`
      );
    }
  });

  socket.on("error", (err) => {
    console.error(
      `[Server] WebSocket error on client socket: ${err.message} (ID: ${currentConnectionId})`
    );
    if (socket.readyState === WebSocket.OPEN) {
      socket.close();
    }
  });
});
