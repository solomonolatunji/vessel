package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"vessl.dev/vessl/internal/models"
)

// GetSchemas returns table and column definitions.
func (s *DatabaseService) GetSchemas(ctx context.Context, id string) ([]models.TableSchema, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	switch db.Engine {
	case "postgresql", "postgres":
		return getPostgresSchemas(db)
	case "mysql", "mariadb":
		return getMySQLSchemas(db)
	default:
		return nil, fmt.Errorf("schema introspection not supported for engine: %s", db.Engine)
	}
}

func getPostgresSchemas(db *models.Database) ([]models.TableSchema, error) {
	host := "localhost"
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, db.Port, db.Username, db.Password, db.DatabaseName)
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			tableNames = append(tableNames, t)
		}
	}

	var schemas []models.TableSchema
	for _, t := range tableNames {
		colRows, err := conn.Query("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_schema='public' AND table_name=$1", t)
		if err != nil {
			continue
		}
		var columns []models.ColumnSchema
		for colRows.Next() {
			var cName, cType, cNullable string
			if err := colRows.Scan(&cName, &cType, &cNullable); err == nil {
				columns = append(columns, models.ColumnSchema{
					Name:       cName,
					Type:       cType,
					IsNullable: cNullable == "YES",
					IsPrimary:  false, // Keeping it simple for v1
				})
			}
		}
		colRows.Close()
		schemas = append(schemas, models.TableSchema{Name: t, Columns: columns})
	}
	if schemas == nil {
		schemas = []models.TableSchema{}
	}
	return schemas, nil
}

func getMySQLSchemas(db *models.Database) ([]models.TableSchema, error) {
	host := "localhost"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.Username, db.Password, host, db.Port, db.DatabaseName)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tableNames []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err == nil {
			tableNames = append(tableNames, t)
		}
	}

	var schemas []models.TableSchema
	for _, t := range tableNames {
		colRows, err := conn.Query(fmt.Sprintf("SHOW COLUMNS FROM `%s`", t))
		if err != nil {
			continue
		}
		var columns []models.ColumnSchema
		for colRows.Next() {
			var cField, cType, cNull, cKey string
			var cDefault, cExtra sql.NullString
			if err := colRows.Scan(&cField, &cType, &cNull, &cKey, &cDefault, &cExtra); err == nil {
				columns = append(columns, models.ColumnSchema{
					Name:       cField,
					Type:       cType,
					IsNullable: cNull == "YES",
					IsPrimary:  cKey == "PRI",
				})
			}
		}
		colRows.Close()
		schemas = append(schemas, models.TableSchema{Name: t, Columns: columns})
	}
	if schemas == nil {
		schemas = []models.TableSchema{}
	}
	return schemas, nil
}

// GetTableData gets rows from a table.
func (s *DatabaseService) GetTableData(ctx context.Context, id, table string, limit, offset int) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	switch db.Engine {
	case "postgresql", "postgres":
		query := fmt.Sprintf("SELECT * FROM \"%s\" LIMIT %d OFFSET %d", table, limit, offset)
		return s.QueryDatabase(ctx, id, query)
	case "mysql", "mariadb":
		query := fmt.Sprintf("SELECT * FROM `%s` LIMIT %d OFFSET %d", table, limit, offset)
		return s.QueryDatabase(ctx, id, query)
	default:
		return nil, fmt.Errorf("data browsing not supported for engine: %s", db.Engine)
	}
}

func escapeSQLValue(v any) string {
	if v == nil {
		return "NULL"
	}
	str := fmt.Sprintf("%v", v)
	str = strings.ReplaceAll(str, "'", "''")
	return fmt.Sprintf("'%s'", str)
}

// InsertTableRow inserts a row into a table.
func (s *DatabaseService) InsertTableRow(ctx context.Context, id, table string, data map[string]any) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	var cols []string
	var vals []string
	for k, v := range data {
		if db.Engine == "postgresql" || db.Engine == "postgres" {
			cols = append(cols, fmt.Sprintf("\"%s\"", k))
		} else {
			cols = append(cols, fmt.Sprintf("`%s`", k))
		}
		vals = append(vals, escapeSQLValue(v))
	}

	var query string
	switch db.Engine {
	case "postgresql", "postgres":
		query = fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(vals, ", "))
	case "mysql", "mariadb":
		query = fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Join(vals, ", "))
	default:
		return nil, fmt.Errorf("inserts not supported for engine: %s", db.Engine)
	}

	return s.QueryDatabase(ctx, id, query)
}

// UpdateTableRow updates a row.
func (s *DatabaseService) UpdateTableRow(ctx context.Context, id, table string, keys map[string]any, data map[string]any) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	var sets []string
	for k, v := range data {
		if db.Engine == "postgresql" || db.Engine == "postgres" {
			sets = append(sets, fmt.Sprintf("\"%s\"=%s", k, escapeSQLValue(v)))
		} else {
			sets = append(sets, fmt.Sprintf("`%s`=%s", k, escapeSQLValue(v)))
		}
	}

	var wheres []string
	for k, v := range keys {
		if db.Engine == "postgresql" || db.Engine == "postgres" {
			wheres = append(wheres, fmt.Sprintf("\"%s\"=%s", k, escapeSQLValue(v)))
		} else {
			wheres = append(wheres, fmt.Sprintf("`%s`=%s", k, escapeSQLValue(v)))
		}
	}

	if len(wheres) == 0 {
		return nil, errors.New("at least one primary key is required for updates")
	}

	var query string
	switch db.Engine {
	case "postgresql", "postgres":
		query = fmt.Sprintf("UPDATE \"%s\" SET %s WHERE %s", table, strings.Join(sets, ", "), strings.Join(wheres, " AND "))
	case "mysql", "mariadb":
		query = fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", table, strings.Join(sets, ", "), strings.Join(wheres, " AND "))
	default:
		return nil, fmt.Errorf("updates not supported for engine: %s", db.Engine)
	}

	return s.QueryDatabase(ctx, id, query)
}

// DeleteTableRow deletes a row.
func (s *DatabaseService) DeleteTableRow(ctx context.Context, id, table string, keys map[string]any) (*models.DatabaseQueryResponse, error) {
	db, err := s.repo.GetByID(ctx, id)
	if err != nil || db == nil {
		return nil, errors.New("database not found")
	}

	var wheres []string
	for k, v := range keys {
		if db.Engine == "postgresql" || db.Engine == "postgres" {
			wheres = append(wheres, fmt.Sprintf("\"%s\"=%s", k, escapeSQLValue(v)))
		} else {
			wheres = append(wheres, fmt.Sprintf("`%s`=%s", k, escapeSQLValue(v)))
		}
	}

	if len(wheres) == 0 {
		return nil, errors.New("at least one primary key is required for deletes")
	}

	var query string
	switch db.Engine {
	case "postgresql", "postgres":
		query = fmt.Sprintf("DELETE FROM \"%s\" WHERE %s", table, strings.Join(wheres, " AND "))
	case "mysql", "mariadb":
		query = fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, strings.Join(wheres, " AND "))
	default:
		return nil, fmt.Errorf("deletes not supported for engine: %s", db.Engine)
	}

	return s.QueryDatabase(ctx, id, query)
}
