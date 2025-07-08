const WebSocket = require("ws");

const socket = new WebSocket("ws://localhost:3001");

function sendLspMessage(socket, json) {
  const body = JSON.stringify(json);
  const contentLength = Buffer.byteLength(body, "utf8");
  const message = `Content-Length: ${contentLength}\r\n\r\n${body}`;
  socket.send(message);
}

let initialized = false;
let opened = false;
let requestedCompletion = false;

socket.on("open", () => {
  console.log("[Client] WebSocket connected to LSP server.");

  // Step 1: Send initialize request
  sendLspMessage(socket, {
    jsonrpc: "2.0",
    id: 1,
    method: "initialize",
    params: {
      capabilities: {},
      rootUri: null,
      workspaceFolders: null,
    },
  });
});

socket.on("message", (data) => {
  const text = data.toString();
  console.log("[Client] Received from server:\n", text);

  // Parse full LSP payload(s) — handle Content-Length prefix
  const messages = extractLspMessages(text);

  for (const message of messages) {
    if (message.id === 1 && message.result && !initialized) {
      initialized = true;
      console.log("[Client] Sending 'initialized' notification...");
      sendLspMessage(socket, {
        jsonrpc: "2.0",
        method: "initialized",
        params: {},
      });

      setTimeout(() => {
        if (!opened) {
          console.log("[Client] Sending 'didOpen' for virtual file...");
          sendLspMessage(socket, {
            jsonrpc: "2.0",
            method: "textDocument/didOpen",
            params: {
              textDocument: {
                uri: "file:///virtual/test.py",
                languageId: "python",
                version: 1,
                text: `import os\nprint(os.environ.get("HOME"))\n`,
              },
            },
          });
          opened = true;
        }
      }, 200); // Slight delay before didOpen
    }

    if (opened && !requestedCompletion) {
      // Send completion after another short delay
      setTimeout(() => {
        console.log("[Client] Sending 'completion' request...");
        sendLspMessage(socket, {
          jsonrpc: "2.0",
          id: 2,
          method: "textDocument/completion",
          params: {
            textDocument: {
              uri: "file:///virtual/test.py",
            },
            position: {
              line: 1,
              character: 6, // Cursor after `print(`
            },
          },
        });
        requestedCompletion = true;
      }, 400);
    }

    if (message.id === 2 && message.result) {
      console.log("[Client] Received completion items:");
      console.dir(message.result.items ?? message.result, { depth: null });
    }
  }
});

socket.on("error", (err) => {
  console.error("[Client] WebSocket error:", err.message);
});

socket.on("close", () => {
  console.log("[Client] WebSocket closed.");
});

function extractLspMessages(raw) {
  // Parse all Content-Length–framed JSON blocks in raw response
  const messages = [];
  const chunks = raw.split(/Content-Length: \d+\r\n\r\n/).filter(Boolean);
  for (const chunk of chunks) {
    try {
      const json = JSON.parse(chunk);
      messages.push(json);
    } catch (err) {
      console.warn("[Client] Could not parse chunk:", chunk);
    }
  }
  return messages;
}
