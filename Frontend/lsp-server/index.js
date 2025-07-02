const { WebSocketServer } = require("ws");
const { createWebSocketStream } = require("ws");
const {
  createMessageConnection,
  StreamMessageReader,
  StreamMessageWriter,
} = require("vscode-jsonrpc");

const { spawn } = require("child_process");

const wss = new WebSocketServer({ port: 3001 });
console.log("LSP WebSocket server listening on ws://localhost:3001");

wss.on("connection", (socket) => {
  const pyright = spawn("npx", ["pyright-langserver", "--stdio"]);

  const socketStream = createWebSocketStream(socket, {
    encoding: "utf8",
    decodeStrings: false,
  });
  const reader = new StreamMessageReader(socketStream);
  const writer = new StreamMessageWriter(socketStream);

  const lspReader = new StreamMessageReader(pyright.stdout);
  const lspWriter = new StreamMessageWriter(pyright.stdin);

  // Pipe data between client and Pyright
  const clientToServer = reader.listen((msg) => lspWriter.write(msg));
  const serverToClient = lspReader.listen((msg) => writer.write(msg));

  socket.on("close", () => pyright.kill());
});
