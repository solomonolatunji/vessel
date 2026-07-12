# 🛰️ Vessl Development Roadmap & Tasks (`TODO.md`)

---

## 🤖 Phase 7: AI Agent Protocol (MCP) & API Ecosystem

> The MCP server and API Ecosystem is a core feature built into the `vessld` Go daemon directly so self-hosters can use it for free. Cloud-specific proxying lives in the `vessl-cloud` repository.

- [x] **REST API to MCP Bridge**:
  - Expose Vessl's REST API as an MCP server (`@modelcontextprotocol/sdk`) so AI agents (Claude Code, Cursor, etc.) can deploy apps, manage databases, and query logs programmatically.
  - Implement Local stdio transport for the CLI daemon.
- [ ] **SDKs**:
  - Publish an official Vessl API client SDK for Node.js and Go.
