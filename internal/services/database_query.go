package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"vessel.dev/vessel/internal/models"
)

func (s *DatabaseService) QueryDatabase(ctx context.Context, id string, query string) (*models.DatabaseQueryResponse, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	if query == "" {
		return nil, errors.New("query cannot be empty")
	}

	host := "localhost"
	if db.Port == 0 {
		return nil, errors.New("database port not configured")
	}

	switch db.Engine {
	case "postgresql", "postgres":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, db.Port, db.Username, db.Password, db.DatabaseName)
		return s.querySQL("postgres", dsn, query)

	case "mysql", "mariadb":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
			db.Username, db.Password, host, db.Port, db.DatabaseName)
		return s.querySQL("mysql", dsn, query)

	case "redis":
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, db.Port),
			Password: db.Password,
			DB:       0,
		})
		defer client.Close()

		parts := strings.Fields(query)
		if len(parts) == 0 {
			return nil, errors.New("empty command")
		}

		var args []interface{}
		for _, p := range parts {
			args = append(args, p)
		}

		res, err := client.Do(ctx, args...).Result()
		if err != nil {
			return nil, err
		}
		return &models.DatabaseQueryResponse{
			Result: res,
		}, nil

	default:
		return nil, fmt.Errorf("querying not supported for engine: %s", db.Engine)
	}
}

func (s *DatabaseService) querySQL(driver, dsn, query string) (*models.DatabaseQueryResponse, error) {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query(query)
	if err != nil {
		if strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "SELECT") || strings.HasPrefix(strings.ToUpper(strings.TrimSpace(query)), "EXPLAIN") {
			return nil, err
		}
		res, execErr := conn.Exec(query)
		if execErr != nil {
			return nil, err
		}
		rowsAffected, _ := res.RowsAffected()
		return &models.DatabaseQueryResponse{
			Result: fmt.Sprintf("%d rows affected", rowsAffected),
		}, nil
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var resultRows []map[string]any
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			if b, ok := (*val).([]byte); ok {
				m[colName] = string(b)
			} else {
				m[colName] = *val
			}
		}
		resultRows = append(resultRows, m)
	}

	return &models.DatabaseQueryResponse{
		Columns: cols,
		Rows:    resultRows,
	}, nil
}
