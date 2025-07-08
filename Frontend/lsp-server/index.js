const { WebSocketServer } = require("ws");
const { createWebSocketStream } = require("ws");
const { StreamMessageReader, StreamMessageWriter } = require("vscode-jsonrpc");
const { TextDecoder } = require("util");
const { spawn } = require("child_process");
const fs = require("fs");
const { URL } = require("url");

const wss = new WebSocketServer({ port: 3001 });
console.log("[Server] LSP WebSocket server listening on ws://localhost:3001");

let connectionCounter = 0;

wss.on("connection", (socket) => {
  connectionCounter++;
  const connId = connectionCounter;

  console.log(`[Server ${connId}] Client connected.`);

  const pyright = spawn("npx", ["pyright-langserver", "--stdio"], {
    cwd: __dirname,
    env: {
      ...process.env,
      PYTHONPATH:
        "/Library/Frameworks/Python.framework/Versions/3.10/lib/python3.10/site-packages",
    },
  });

  console.log(`[Server ${connId}] Spawned Pyright (PID: ${pyright.pid})`);

  pyright.stderr.on("data", (data) => {
    const msg = new TextDecoder().decode(data);
    console.error(`[Server ${connId}] [Pyright STDERR]: ${msg}`);
  });

  pyright.stdout.on("data", (data) => {
    // Optional debug
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

  const socketStream = createWebSocketStream(socket);
  const clientReader = new StreamMessageReader(socketStream);
  const clientWriter = new StreamMessageWriter(socketStream);
  const pyrightReader = new StreamMessageReader(pyright.stdout);
  const pyrightWriter = new StreamMessageWriter(pyright.stdin);

  clientReader.listen(async (message) => {
    logMessage("Client → Pyright", message, connId);

    if (
      message.method === "textDocument/didOpen" ||
      message.method === "textDocument/didChange"
    ) {
      try {
        const uri = message.params.textDocument.uri;
        const text =
          message.method === "textDocument/didOpen"
            ? message.params.textDocument.text
            : message.params.contentChanges[0]?.text;

        const filePath = new URL(uri).pathname;

        await fs.promises.writeFile(filePath, text, "utf8");
        console.log(`[Server ${connId}] Wrote file to disk: ${filePath}`);
      } catch (err) {
        console.error(`[Server ${connId}] Failed to write file:`, err);
      }
    }

    pyrightWriter.write(message);
  });

  pyrightReader.listen((message) => {
    console.log("🟢 Response from Pyright:", message);
    clientWriter.write(message);
    logMessage("Pyright → Client", message, connId);
  });

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
