// create-init-request.js
const fs = require("fs");

const body = JSON.stringify({
  jsonrpc: "2.0",
  id: 3,
  method: "initialize",
  params: {
    capabilities: {},
    rootUri: null,
    workspaceFolders: null,
  },
});

const payload =
  `Content-Length: ${Buffer.byteLength(body, "utf8")}\r\n\r\n` + body;

fs.writeFileSync("initialize.bin", payload);
