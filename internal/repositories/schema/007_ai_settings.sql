CREATE TABLE IF NOT EXISTS ai_settings (
    id TEXT PRIMARY KEY,
    default_provider TEXT DEFAULT 'none',
    openai_key TEXT,
    anthropic_key TEXT,
    google_key TEXT,
    mistral_key TEXT,
    groq_key TEXT,
    deepseek_key TEXT,
    xai_key TEXT,
    moonshot_key TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
