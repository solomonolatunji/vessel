CREATE TABLE IF NOT EXISTS dns_records (
    id TEXT PRIMARY KEY,
    domain_name TEXT NOT NULL,
    record_type TEXT NOT NULL,
    record_name TEXT NOT NULL,
    record_value TEXT NOT NULL,
    ttl INTEGER DEFAULT 3600,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    action TEXT NOT NULL,
    resource TEXT NOT NULL,
    details TEXT,
    ip_address TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ai_settings (
    id TEXT PRIMARY KEY,
    default_provider TEXT DEFAULT 'none',
    openai_key TEXT DEFAULT '',
    openai_model TEXT DEFAULT '',
    anthropic_key TEXT DEFAULT '',
    anthropic_model TEXT DEFAULT '',
    google_key TEXT DEFAULT '',
    google_model TEXT DEFAULT '',
    mistral_key TEXT DEFAULT '',
    mistral_model TEXT DEFAULT '',
    groq_key TEXT DEFAULT '',
    groq_model TEXT DEFAULT '',
    deepseek_key TEXT DEFAULT '',
    deepseek_model TEXT DEFAULT '',
    xai_key TEXT DEFAULT '',
    xai_model TEXT DEFAULT '',
    moonshot_key TEXT DEFAULT '',
    moonshot_model TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notification_settings (
    id TEXT PRIMARY KEY,
    discord_webhook_url TEXT,
    discord_ping_enabled BOOLEAN DEFAULT FALSE,
    discord_enabled BOOLEAN DEFAULT FALSE,
    slack_webhook_url TEXT,
    slack_enabled BOOLEAN DEFAULT FALSE,
    telegram_bot_token TEXT,
    telegram_chat_id TEXT,
    telegram_enabled BOOLEAN DEFAULT FALSE,
    smtp_host TEXT,
    smtp_port INTEGER DEFAULT 587,
    smtp_user TEXT,
    smtp_password TEXT,
    smtp_from_name TEXT,
    smtp_from_address TEXT,
    smtp_enabled BOOLEAN DEFAULT FALSE,
    resend_api_key TEXT,
    resend_enabled BOOLEAN DEFAULT FALSE,
    pushover_user_key TEXT,
    pushover_api_token TEXT,
    pushover_enabled BOOLEAN DEFAULT FALSE,
    generic_webhook_url TEXT,
    generic_webhook_enabled BOOLEAN DEFAULT FALSE,
    notification_alerts BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_dns_records_domain_name ON dns_records(domain_name);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id_created_at ON audit_logs(user_id, created_at DESC);
