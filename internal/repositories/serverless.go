package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"vessel.dev/vessel/internal/models"
)

type ServerlessRepository interface {
	SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*models.ServerlessFunctionCode, error)
	GetCodeByServiceID(ctx context.Context, serviceID string) (*models.ServerlessFunctionCode, error)
}

type sqliteServerlessRepo struct {
	db *sql.DB
}

func NewServerlessRepository(db *sql.DB) ServerlessRepository {
	return &sqliteServerlessRepo{db: db}
}

func (r *sqliteServerlessRepo) SaveCode(ctx context.Context, serviceID, runtime, codeContent string) (*models.ServerlessFunctionCode, error) {
	// check if exists
	existing, err := r.GetCodeByServiceID(ctx, serviceID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	now := time.Now()
	if existing != nil {
		query := `UPDATE serverless_functions_code SET runtime = ?, code_content = ?, updated_at = ? WHERE service_id = ?`
		_, err := r.db.ExecContext(ctx, query, runtime, codeContent, now, serviceID)
		if err != nil {
			return nil, err
		}
		existing.Runtime = runtime
		existing.CodeContent = codeContent
		existing.UpdatedAt = now
		return existing, nil
	}

	id := uuid.New().String()
	query := `INSERT INTO serverless_functions_code (id, service_id, runtime, code_content, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err = r.db.ExecContext(ctx, query, id, serviceID, runtime, codeContent, now, now)
	if err != nil {
		return nil, err
	}

	return &models.ServerlessFunctionCode{
		ID:          id,
		ServiceID:   serviceID,
		Runtime:     runtime,
		CodeContent: codeContent,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *sqliteServerlessRepo) GetCodeByServiceID(ctx context.Context, serviceID string) (*models.ServerlessFunctionCode, error) {
	query := `SELECT id, service_id, runtime, code_content, created_at, updated_at FROM serverless_functions_code WHERE service_id = ?`
	row := r.db.QueryRowContext(ctx, query, serviceID)

	var code models.ServerlessFunctionCode
	err := row.Scan(&code.ID, &code.ServiceID, &code.Runtime, &code.CodeContent, &code.CreatedAt, &code.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &code, nil
}
