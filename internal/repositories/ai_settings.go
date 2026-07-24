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

	"codedock.run/codedock/internal/models"
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

const aiSettingsColumns = `id, default_provider, openai_key, openai_model, anthropic_key, anthropic_model, google_key, google_model, mistral_key, mistral_model, groq_key, groq_model, deepseek_key, deepseek_model, xai_key, xai_model, moonshot_key, moonshot_model, created_at, updated_at`

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
		&cfg.OpenAIKey, &cfg.OpenAIModel, &cfg.AnthropicKey, &cfg.AnthropicModel, &cfg.GoogleKey, &cfg.GoogleModel, &cfg.MistralKey, &cfg.MistralModel, &cfg.GroqKey, &cfg.GroqModel, &cfg.DeepSeekKey, &cfg.DeepSeekModel, &cfg.XAIKey, &cfg.XAIModel, &cfg.MoonshotKey, &cfg.MoonshotModel,
		&cfg.CreatedAt, &cfg.UpdatedAt,
	)
}

func aiSettingsArgs(cfg *models.AISettings) []any {
	return []any{
		cfg.ID, cfg.DefaultProvider,
		cfg.OpenAIKey, cfg.OpenAIModel, cfg.AnthropicKey, cfg.AnthropicModel, cfg.GoogleKey, cfg.GoogleModel, cfg.MistralKey, cfg.MistralModel, cfg.GroqKey, cfg.GroqModel, cfg.DeepSeekKey, cfg.DeepSeekModel, cfg.XAIKey, cfg.XAIModel, cfg.MoonshotKey, cfg.MoonshotModel,
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
	          openai_model = excluded.openai_model,
	          anthropic_key = excluded.anthropic_key,
	          anthropic_model = excluded.anthropic_model,
	          google_key = excluded.google_key,
	          google_model = excluded.google_model,
	          mistral_key = excluded.mistral_key,
	          mistral_model = excluded.mistral_model,
	          groq_key = excluded.groq_key,
	          groq_model = excluded.groq_model,
	          deepseek_key = excluded.deepseek_key,
	          deepseek_model = excluded.deepseek_model,
	          xai_key = excluded.xai_key,
	          xai_model = excluded.xai_model,
	          moonshot_key = excluded.moonshot_key,
	          moonshot_model = excluded.moonshot_model,
	          updated_at = excluded.updated_at`, aiSettingsColumns, aiSettingsPlaceholders())
	_, err := r.db.ExecContext(ctx, query, aiSettingsArgs(cfg)...)
	if err != nil {
		return fmt.Errorf("failed to update ai settings: %w", err)
	}
	return nil
}
