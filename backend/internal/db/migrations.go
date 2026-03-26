package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations reads all *.sql files from migrationsDir in lexical order
// and executes them against the PostgreSQL pool. Each file is idempotent
// (uses IF NOT EXISTS / CREATE OR REPLACE).
func RunMigrations(ctx context.Context, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("migrations: read dir %q: %w", migrationsDir, err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, filepath.Join(migrationsDir, e.Name()))
		}
	}
	sort.Strings(files)

	for _, f := range files {
		sql, err := os.ReadFile(f)
		if err != nil {
			return fmt.Errorf("migrations: read %q: %w", f, err)
		}
		if _, err := PG.Exec(ctx, string(sql)); err != nil {
			return fmt.Errorf("migrations: execute %q: %w", f, err)
		}
		log.Printf("[db] migration applied: %s", filepath.Base(f))
	}
	return nil
}
