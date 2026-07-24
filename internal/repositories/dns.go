package repositories

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"codedock.run/codedock/internal/models"
	"codedock.run/codedock/internal/utils"
)

type DNSRepository interface {
	Create(ctx context.Context, record *models.DNSRecord) error
	GetByID(ctx context.Context, id string) (*models.DNSRecord, error)
	ListByDomain(ctx context.Context, domainName string) ([]*models.DNSRecord, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, record *models.DNSRecord) error
}

type DNSRepo struct {
	db *sqlx.DB
	mu sync.Mutex
}

func NewDNSRepo(db *sql.DB) *DNSRepo {
	return &DNSRepo{db: sqlx.NewDb(db, "sqlite")}
}

func (r *DNSRepo) Create(_ context.Context, record *models.DNSRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if record.ID == "" {
		record.ID = uuid.NewString()
	}
	if record.TTL == 0 {
		record.TTL = 3600
	}
	now := time.Now()
	record.CreatedAt = now
	record.UpdatedAt = now

	_, err := r.db.Exec(`INSERT INTO dns_records (
		id, domain_name, record_type, record_name, record_value, ttl, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		record.ID, record.DomainName, record.RecordType, record.RecordName, record.RecordValue, record.TTL, record.CreatedAt, record.UpdatedAt)
	return err
}

func (r *DNSRepo) GetByID(_ context.Context, id string) (*models.DNSRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var record models.DNSRecord
	err := r.db.Get(&record, `SELECT * FROM dns_records WHERE id = ?`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.NewNotFoundError("DNSRecord", id)
	}
	return &record, err
}

func (r *DNSRepo) ListByDomain(_ context.Context, domainName string) ([]*models.DNSRecord, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var list []*models.DNSRecord
	var err error
	if domainName == "" {
		err = r.db.Select(&list, `SELECT * FROM dns_records ORDER BY created_at ASC`)
	} else {
		err = r.db.Select(&list, `SELECT * FROM dns_records WHERE domain_name = ? ORDER BY created_at ASC`, domainName)
	}
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = make([]*models.DNSRecord, 0)
	}
	return list, nil
}

func (r *DNSRepo) Update(_ context.Context, record *models.DNSRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	record.UpdatedAt = time.Now()
	_, err := r.db.Exec(`UPDATE dns_records SET record_type = ?, record_name = ?, record_value = ?, ttl = ?, updated_at = ? WHERE id = ?`,
		record.RecordType, record.RecordName, record.RecordValue, record.TTL, record.UpdatedAt, record.ID)
	return err
}

func (r *DNSRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, err := r.db.Exec(`DELETE FROM dns_records WHERE id = ?`, id)
	return err
}
