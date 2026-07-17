package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"vessl.dev/vessl/internal/models"
)

type AISettingsRepository interface {
	GetAISettings(ctx context.Context) (*models.AISettings, error)
	UpdateAISettings(ctx context.Context, cfg *models.AISettings) error
}

type AISettingsRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewAISettingsRepo(db *sql.DB) *AISettingsRepo {
	return &AISettingsRepo{db: sqlx.NewDb(db, "sqlite")}
}

const aiSettingsColumns = `id, default_provider, openai_key, anthropic_key, google_key, mistral_key, groq_key, deepseek_key, xai_key, moonshot_key, created_at, updated_at`

func aiSettingsPlaceholders() string {
	columns := strings.Split(aiSettingsColumns, ",")
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}

func scanAISettings(scanner interface{ Scan(dest ...any) error }, cfg *models.AISettings) error {
	return scanner.Scan(
		&cfg.ID, &cfg.DefaultProvider,
		&cfg.OpenAIKey, &cfg.AnthropicKey, &cfg.GoogleKey, &cfg.MistralKey, &cfg.GroqKey, &cfg.DeepSeekKey, &cfg.XAIKey, &cfg.MoonshotKey,
		&cfg.CreatedAt, &cfg.UpdatedAt,
	)
}

func aiSettingsArgs(cfg *models.AISettings) []any {
	return []any{
		cfg.ID, cfg.DefaultProvider,
		cfg.OpenAIKey, cfg.AnthropicKey, cfg.GoogleKey, cfg.MistralKey, cfg.GroqKey, cfg.DeepSeekKey, cfg.XAIKey, cfg.MoonshotKey,
		cfg.CreatedAt, cfg.UpdatedAt,
	}
}

func (r *AISettingsRepo) GetAISettings(ctx context.Context) (*models.AISettings, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var cfg models.AISettings
	err := scanAISettings(r.db.QueryRowContext(ctx, `SELECT `+aiSettingsColumns+` FROM ai_settings WHERE id = 'global' LIMIT 1`), &cfg)
	if errors.Is(err, sql.ErrNoRows) {
		defaultSettings := &models.AISettings{
			ID:              "global",
			DefaultProvider: "none",
			CreatedAt:       time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
		}
		query := fmt.Sprintf(`INSERT INTO ai_settings (%s) VALUES (%s)`, aiSettingsColumns, aiSettingsPlaceholders())
		_, _ = r.db.ExecContext(ctx, query, aiSettingsArgs(defaultSettings)...)
		return defaultSettings, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ai settings: %w", err)
	}
	return &cfg, nil
}

func (r *AISettingsRepo) UpdateAISettings(ctx context.Context, cfg *models.AISettings) error {
	if cfg.ID == "" {
		cfg.ID = "global"
	}
	cfg.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	r.mu.Lock()
	defer r.mu.Unlock()
	query := fmt.Sprintf(`INSERT INTO ai_settings (%s)
	          VALUES (%s)
	          ON CONFLICT(id) DO UPDATE SET
	          default_provider = excluded.default_provider,
	          openai_key = excluded.openai_key,
	          anthropic_key = excluded.anthropic_key,
	          google_key = excluded.google_key,
	          mistral_key = excluded.mistral_key,
	          groq_key = excluded.groq_key,
	          deepseek_key = excluded.deepseek_key,
	          xai_key = excluded.xai_key,
	          moonshot_key = excluded.moonshot_key,
	          updated_at = excluded.updated_at`, aiSettingsColumns, aiSettingsPlaceholders())
	_, err := r.db.ExecContext(ctx, query, aiSettingsArgs(cfg)...)
	if err != nil {
		return fmt.Errorf("failed to update ai settings: %w", err)
	}
	return nil
}
