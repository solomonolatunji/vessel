package project

import "context"

// Repository defines data operations for the project domain.
type Repository interface {
	List(ctx context.Context) ([]ProjectConfig, error)
	Get(ctx context.Context, id string) (*ProjectConfig, error)
	Create(ctx context.Context, p *ProjectConfig) error
	Delete(ctx context.Context, id string) error
}
