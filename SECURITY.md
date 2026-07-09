# Security Policy 🔒

Vessel takes infrastructure and deployment security extremely seriously. Because `vesseld` interacts directly with your server's Docker daemon (`/var/run/docker.sock`) and manages sensitive `.env` secrets, we adhere to strict security practices.

---

## 🛡️ Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| `0.1.x` | :white_check_mark: |

---

## 🚨 Reporting a Vulnerability

If you discover a security vulnerability in Vessel (such as unauthorized container access, `.env` secret leakage, or authentication bypass), **please do not report it publicly via GitHub Issues**.

Instead, please send an email or private advisory to our maintainers at:
**<security@vessel.dev>** (or open a private GitHub Security Advisory).

### What to Include in Your Report

- Type of issue (e.g., buffer overflow, SQL injection, privilege escalation, cross-site scripting).
- Full paths of source file(s) related to the vulnerability.
- Step-by-step instructions or proof-of-concept (PoC) code to reproduce the vulnerability safely.

We commit to acknowledging all security disclosures within **48 hours** and issuing prompt patches via automated zero-downtime upgrades (`scripts/upgrade.sh`).
