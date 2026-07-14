const http = require("node:http");
const fs = require("node:fs");
const path = require("node:path");

const PORT = process.env.PORT || 80;
const SCRIPTS = new Set(["/install.sh", "/upgrade.sh", "/cli"]);

const server = http.createServer((req, res) => {
  if (req.url === "/version") {
    res.writeHead(200, { "Content-Type": "application/json" });
    res.end(JSON.stringify({ version: "0.1.0", latest: true }));
    return;
  }

  let filename = path.basename(req.url);
  if (req.url === "/") filename = "install.sh";
  if (req.url === "/cli") filename = "install-cli.sh";

  const file = path.join(__dirname, filename);
  const name = path.basename(req.url === "/" ? "/install.sh" : req.url);

  if (!SCRIPTS.has(req.url === "/" ? "/install.sh" : req.url) || !fs.existsSync(file)) {
    res.writeHead(404);
    res.end("Not found");
    return;
  }

  res.writeHead(200, { "Content-Type": "text/x-shellscript" });
  fs.createReadStream(file).pipe(res);
});

server.listen(PORT, () => {
  console.log(`📦 install server running on port ${PORT}`);
});
