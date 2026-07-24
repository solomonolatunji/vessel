# Codedock Multi-Server & Cloud Roadmap

## 1. How Dokploy Achieved Multi-Server & Cloud (and How We Compare)

**Dokploy:**
Dokploy achieves multi-server by having a primary Node.js Next.js instance (the control plane) that manages multiple separate Docker daemons across different servers. They use a unified TypeScript codebase (tRPC) for perfect type safety between the frontend and the backend.

**Codedock (Are we doing the same, better, or worse?):**
**Better in Architecture:** Codedock uses an outbound WebSocket Worker Daemon approach. Instead of the control plane connecting *to* the workers (which requires open ports/firewall exceptions on every worker node), our workers dial *out* to the control plane. This is essential for NAT traversal and SaaS deployments.
**Better in Performance:** We use Go for both the control plane and workers, ensuring minimal resource footprint on the worker nodes compared to running Node.js.

## 2. Features Dokploy Has That Codedock Doesn't (Yet)

- [ ] **PR Previews:** Dokploy automatically spins up ephemeral environments when a Pull Request is opened on GitHub, and tears it down when closed.
- [ ] **More Git Providers:** Dokploy supports GitLab, Bitbucket, and Gitea.
- [ ] **Private Docker Registries:** Linking AWS ECR, Google GCR, or private DockerHub accounts.
- [ ] **Organizations & Teams (Global):** Cross-project organization RBAC system.
- [ ] **Volume Backups:** S3-backed persistent Docker Volume backups.

## 3. A Feature We Both Had, But WE Did Totally Wrong

- [x] **Type Safety between Backend and Workers (The NATS / Type Safety issue):**
  - *Dokploy's win:* Because Dokploy is 100% TypeScript, they use tRPC. If a variable is renamed in the backend, the frontend instantly throws a compiler error.
  - *Our historical flaw:* Our Go backend was sending raw JSON payloads over NATS to workers without a shared source of truth. If a Go struct changed, the worker crashed in production!
	- *The fix (DONE):* We transitioned from NATS to a centralized WebSocket `WorkerHub`. The control plane and workers now import the exact same shared Go schemas (`internal/models/worker.go`). We now have complete compile-time type safety across our multi-server environment!

## 4. Major Architectural Wins over Dokploy

- [x] **Decentralized Builds (Worker-side Building):** 
  - *Dokploy's flow:* The control plane handles git cloning, building the Docker image, and sending it to the remote server. This puts an immense CPU and RAM load on the control plane, leading to OOM (Out Of Memory) crashes if multiple projects deploy simultaneously.
  - *Codedock's flow (DONE):* We moved the entire Git Clone + Docker Build + Container Start pipeline to the **worker node itself**. The control plane simply sends a `WorkerDeployAppPayload` command over WebSocket. The worker parses the config, clones the repo natively on the remote machine, builds the Nixpacks/Docker image, and deploys it. This allows the control plane to scale infinitely without resource exhaustion.
