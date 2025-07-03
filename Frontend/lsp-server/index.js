const { WebSocketServer } = require("ws");
const { createWebSocketStream } = require("ws");
const { StreamMessageReader, StreamMessageWriter } = require("vscode-jsonrpc");
const { TextDecoder } = require("util");
const { spawn } = require("child_process");

const wss = new WebSocketServer({ port: 3001 });
console.log("[Server] LSP WebSocket server listening on ws://localhost:3001");

let connectionCounter = 0;

wss.on("connection", (socket) => {
  connectionCounter++;
  const connId = connectionCounter;

  console.log(`[Server ${connId}] Client connected.`);

  // Start Pyright LSP server
  const pyright = spawn("npx", ["pyright-langserver", "--stdio"]);
  console.log(`[Server ${connId}] Spawned Pyright (PID: ${pyright.pid})`);

  // Pyright diagnostic output
  pyright.stderr.on("data", (data) => {
    const msg = new TextDecoder().decode(data);
    console.error(`[Server ${connId}] [Pyright STDERR]: ${msg}`);
  });

  pyright.stdout.on("data", (data) => {
    // Optional debug:
    // console.log(`[Server ${connId}] [Pyright STDOUT]: ${data.toString()}`);
  });

  pyright.on("exit", (code, signal) => {
    console.log(
      `[Server ${connId}] Pyright exited. Code=${code}, Signal=${signal}`
    );
    socket.close();
  });

  pyright.on("error", (err) => {
    console.error(`[Server ${connId}] Pyright error: ${err.message}`);
    socket.close();
  });

  // WebSocket ↔ Pyright stream
  const socketStream = createWebSocketStream(socket);
  const clientReader = new StreamMessageReader(socketStream);
  const clientWriter = new StreamMessageWriter(socketStream);

  const pyrightReader = new StreamMessageReader(pyright.stdout);
  const pyrightWriter = new StreamMessageWriter(pyright.stdin);

  // Forward client → pyright
  clientReader.listen((message) => {
    logMessage("Client → Pyright", message, connId);
    pyrightWriter.write(message);
  });

  // Forward pyright → client
  pyrightReader.listen((message) => {
    clientWriter.write(message);
    logMessage("Pyright → Client", message, connId);
  });

  // Handle WebSocket close
  socket.on("close", () => {
    console.log(`[Server ${connId}] WebSocket closed. Cleaning up...`);
    clientReader.dispose();
    clientWriter.dispose();
    pyrightReader.dispose();
    pyrightWriter.dispose();
    if (!pyright.killed) {
      pyright.kill();
    }
  });

  // Handle WebSocket error
  socket.on("error", (err) => {
    console.error(`[Server ${connId}] WebSocket error: ${err.message}`);
    socket.destroy();
  });
});

function logMessage(direction, message, connId) {
  const type = message.method
    ? `Method: ${message.method}`
    : message.id
    ? `Response to ID: ${message.id}`
    : "Other";
  console.log(
    `[Server ${connId}] ${direction} (${type}): ${JSON.stringify(message)}`
  );
}
