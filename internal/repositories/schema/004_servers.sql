CREATE TABLE IF NOT EXISTS servers (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    name TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    status TEXT DEFAULT 'provisioning',
    worker_token TEXT NOT NULL,
    last_seen_at DATETIME,
    metrics TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_servers_user_id ON servers(user_id);
CREATE INDEX IF NOT EXISTS idx_servers_worker_token ON servers(worker_token);

ALTER TABLE projects ADD COLUMN server_id TEXT REFERENCES servers(id) ON DELETE SET NULL;
