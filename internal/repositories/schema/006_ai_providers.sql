ALTER TABLE server_settings ADD COLUMN default_google_key TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN default_mistral_key TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN default_groq_key TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN default_deepseek_key TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN default_xai_key TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN default_moonshot_key TEXT DEFAULT '';
ALTER TABLE server_settings ADD COLUMN default_ai_provider TEXT DEFAULT 'openai';
