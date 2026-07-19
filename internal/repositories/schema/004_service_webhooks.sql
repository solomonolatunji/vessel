CREATE TABLE IF NOT EXISTS service_webhooks (
    id TEXT PRIMARY KEY,
    service_id TEXT NOT NULL REFERENCES app_services(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    event_types TEXT DEFAULT '',
    include_pr_environments BOOLEAN DEFAULT FALSE,
    created_at DATETIME,
    updated_at DATETIME
);

CREATE INDEX IF NOT EXISTS idx_service_webhooks_service_id_created_at ON service_webhooks(service_id, created_at DESC);

-- Ensure the old table exists to prevent SQLite compilation errors on fresh installs
CREATE TABLE IF NOT EXISTS project_webhooks (
    id TEXT PRIMARY KEY,
    project_id TEXT,
    url TEXT,
    event_types TEXT,
    include_pr_environments BOOLEAN,
    created_at DATETIME,
    updated_at DATETIME
);

-- Migrate existing data
INSERT INTO service_webhooks (id, service_id, url, event_types, include_pr_environments, created_at, updated_at)
SELECT w.id, s.id, w.url, w.event_types, w.include_pr_environments, w.created_at, w.updated_at
FROM project_webhooks w
JOIN app_services s ON w.project_id = s.project_id
WHERE NOT EXISTS (
    SELECT 1 FROM service_webhooks sw WHERE sw.id = w.id
);

-- Safely drop the old table now that it's no longer used
DROP TABLE project_webhooks;
