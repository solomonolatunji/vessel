# Vessl Security & Code Quality Audit Report

**Date:** July 14, 2026  
**Scope:** Go backend (`internal/`, `cmd/`), bootstrap scripts, shell scripts

---

## Executive Summary

The Vessl codebase demonstrates solid architectural foundations with proper layered design (handlers → services → repositories). However, a few security vulnerabilities and code quality issues require attention.

**Note:** The previous audit contained several critical **False Positives** (such as SQL Injection and Missing Authorization), which have been debunked and corrected in this revised report.

**Risk Level:** MEDIUM  
**Estimated Remediation Time:** 3-5 days

---

## 1. Security Vulnerabilities

### 1.1 Critical Issues

#### 🔴 CORS Wildcard Configuration

**File:** `internal/http/setup.go`  
**Issue:** `AllowOrigins: []string{"*"}` permits any origin.  
**Risk:** Cross-origin attacks, data theft.  
**Fix:** Replace with an explicit domain whitelist matching the deployment domain.

```go
AllowOrigins: []string{"https://yourdomain.com", "https://app.yourdomain.com"}
```

#### 🔴 Argument Injection in Git Service

**File:** `internal/services/git_service.go`  
**Lines:** Branch parameter passed directly as an argument to `exec.CommandContext("git", "clone", ... "-b", branch)`.  
**Risk:** Argument injection (e.g., passing `--upload-pack=...` as a branch name) leading to potential command execution or unauthorized file access.  
**Fix:** Validate branch names with regex before executing git commands.

```go
validBranch := regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
if !validBranch.MatchString(branch) {
    return errors.New("invalid branch name")
}
```

### 1.2 Deprecated Packages

#### 🟠 Deprecated PostgreSQL Driver

**Files:** `internal/services/database_query.go`, `go.mod`  
**Issue:** Uses `github.com/lib/pq` to query user PostgreSQL databases.  
**Risk:** `lib/pq` is officially in maintenance mode and deprecated. It lacks support for newer PostgreSQL features and may not receive timely security patches.  
**Fix:** Migrate to `github.com/jackc/pgx/v5`.

---

## 2. Code Quality Issues

### 2.1 Architectural Concerns

#### 🟡 os.Exit() in Handler

**File:** `internal/handlers/system.go`  
**Issue:** `os.Exit(0)` in HTTP handler terminates the entire server abruptly.  
**Fix:** Perform a graceful shutdown using Echo's `Shutdown(ctx)` method to ensure pending requests and database connections are safely closed.

### 2.2 Shell Script Analysis

#### 🟡 Unquoted Variables in `.env` Generation

**File:** `bootstrap/install.sh`  
**Issue:** Shell injection risk with unquoted variables in the `.env` heredoc generation.  
**Risk:** If user input (like `TLS_EMAIL`) contains spaces or malicious payloads, it breaks the `.env` formatting.  
**Fix:** Quote all variables in the heredoc.

```bash
cat > "$VESSL_DIR/.env" <<EOF
VESSL_TLS_EMAIL="${TLS_EMAIL}"
VESSL_WILDCARD_DOMAIN="${WILDCARD_DOMAIN}"
EOF
```

---

## 3. False Positives (Debunked from Previous Audit)

The previous audit flagged several high-severity issues. After deep inspection of the codebase, these were proven to be **FALSE POSITIVES** and are NOT vulnerabilities:

1. **SQL Injection Vulnerabilities:**
   - _Claim:_ String concatenation in SQL queries.
   - _Reality:_ `fmt.Sprintf` is only used safely to format static column names (e.g., `serverSettingsColumns`) during DB initialization. All user-supplied data correctly uses parameterized queries (e.g., `?` bindings).
2. **Missing Authorization Checks:**
   - _Claim:_ 157 handler functions without explicit `RequireAuth` middleware.
   - _Reality:_ Vessl uses Echo Route Groups (e.g., `authGroup.Use(s.authGuard.RequireAuth())` in `routes.go`). The handlers are inherently protected by the group router.
3. **Context.Background() Overuse:**
   - _Claim:_ Used inappropriately in long-running operations.
   - _Reality:_ `context.Background()` is appropriately wrapped with `context.WithTimeout()` whenever instantiated inline, preventing resource leaks.

---

## 4. Architecture Tooling Recommendations

Based on the `AGENTS.md` mandate to use **CGO-Free SQLite (`modernc.org/sqlite`)**, here is the verdict on incorporating `sqlx`, `goose`, and `golang-migrate`:

### ✅ `sqlx` (Recommended)

Yes, you should use `sqlx`. It integrates perfectly with `modernc.org/sqlite` (simply specify `"sqlite"` as the driver name). It strictly adheres to the CGO-free rule while significantly reducing boilerplate code (e.g., eliminating manual `rows.Scan` loops).

### ✅ `goose` (Recommended)

Yes, Goose is the ideal migration tool for this stack. It natively allows you to pass your own `*sql.DB` instance or use `modernc.org/sqlite` simply by registering the driver and setting the dialect to `"sqlite3"`. It works seamlessly in CGO-free setups.

### ❌ `golang-migrate` (Not Recommended)

No, avoid `golang-migrate`. The standard CLI and toolset historically rely on the CGO-based `mattn/go-sqlite3` driver. While technically possible to use as a library with `modernc`, it involves writing custom runner wrappers and introduces unnecessary complexity compared to `goose`.

---

**Report Generated:** July 14, 2026  
**Auditor:** Automated Deep Audit & Review
