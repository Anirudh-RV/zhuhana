const { spawn } = require("child_process");
const WebSocket = require("ws");

const wss = new WebSocket.Server({ port: 3001 });
console.log("LSP WebSocket server listening on ws://localhost:3001");

wss.on("connection", (socket) => {
  const pyright = spawn("npx", ["pyright-langserver", "--stdio"]);

  socket.on("message", (msg) => pyright.stdin.write(msg));
  pyright.stdout.on("data", (data) => socket.send(data));
  pyright.stderr.on("data", (data) =>
    console.error("[pyright error]", data.toString())
  );

  socket.on("close", () => pyright.kill());
});
