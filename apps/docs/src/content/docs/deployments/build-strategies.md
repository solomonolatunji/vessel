---
title: Build Strategies
description: Codedock auto-detects the best build strategy based on your project.
---

Codedock supports multiple build strategies and deployment workflows to get your applications online. It auto-detects the best build strategy based on your project, but you can override it per deployment.

### Dockerfile

If your repository contains a `Dockerfile` at the root, Codedock uses it by default.

```dockerfile
FROM node:22-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
CMD ["node", "index.js"]
```

No additional configuration needed — just push and deploy.

### Railpack

[Railpack](https://railpack.com) auto-detects your language and framework. Supported stacks:

- Node.js
- Go
- Python
- Rust
- PHP
- Ruby
- Static sites (HTML/CSS/JS)

Railpack generates an optimal Dockerfile for your project without you writing one.

### Nixpacks

[Nixpacks](https://nixpacks.com) uses Nix expressions to build reproducible environments. It supports the same languages as Railpack plus additional ecosystem tools.

### Buildpacks

Cloud Native Buildpacks support is available for OCI-compliant builds. Select **Buildpacks** in the deployment settings.

### Build Overrides

When using auto-detection builders like Railpack or Nixpacks, you can customize the pipeline without maintaining a Dockerfile by providing overrides in your service settings:

- **Install Command**: Override package dependency installation (`npm ci`, `pip install -r requirements.txt`, `go mod download`).
- **Build Command**: Override compilation and asset generation (`npm run build`, `go build -o app ./cmd`).
- **Start Command**: Override the container execution command (`npm start`, `./app`).

Codedock injects these flags directly into the builder CLI (`--install-cmd`, `--build-cmd`, `--start-cmd`) or synthesizes the corresponding `RUN` and `CMD` instructions inside fallback build containers automatically.
