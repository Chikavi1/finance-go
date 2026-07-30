package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type migration struct {
	version string
	name    string
	sql     string
}

func RunMigrations(pool *pgxpool.Pool, migrationsPath string) error {
	ctx := context.Background()

	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("resolve migrations path: %w", err)
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("migrations directory not found: %s", absPath)
		}
		return fmt.Errorf("read migrations directory: %w", err)
	}

	re := regexp.MustCompile(`^(\d+)_(.+)\.up\.sql$`)

	var migrations []migration
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		matches := re.FindStringSubmatch(e.Name())
		if matches == nil {
			continue
		}

		sql, err := os.ReadFile(filepath.Join(absPath, e.Name()))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", e.Name(), err)
		}

		migrations = append(migrations, migration{
			version: matches[1],
			name:    matches[2],
			sql:     string(sql),
		})
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].version < migrations[j].version
	})

	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`); err != nil {
		return fmt.Errorf("ensure schema_migrations table: %w", err)
	}

	rows, err := pool.Query(ctx, `SELECT version FROM schema_migrations ORDER BY version`)
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	applied := make(map[string]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			rows.Close()
			return fmt.Errorf("scan migration version: %w", err)
		}
		applied[v] = true
	}
	rows.Close()

	for _, m := range migrations {
		if applied[m.version] {
			continue
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin tx for %s_%s: %w", m.version, m.name, err)
		}

		for _, stmt := range splitSQL(m.sql) {
			if _, err := tx.Exec(ctx, stmt); err != nil {
				tx.Rollback(ctx)
				return fmt.Errorf("execute %s_%s: %w", m.version, m.name, err)
			}
		}

		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, m.version); err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("record migration %s_%s: %w", m.version, m.name, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit %s_%s: %w", m.version, m.name, err)
		}
	}

	return nil
}

func splitSQL(sql string) []string {
	parts := strings.Split(sql, ";")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		result = append(result, p)
	}
	return result
}
