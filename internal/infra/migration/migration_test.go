package migration

import (
	"path/filepath"
	"testing"

	"gin-layout/config"
	"gin-layout/internal/infra"
)

func TestRunCreatesApplicationTables(t *testing.T) {
	db, err := infra.NewDatabase(&config.DatabaseConfig{
		Driver: "sqlite",
		DSN:    filepath.Join(t.TempDir(), "migration.db"),
	})
	if err != nil {
		t.Fatalf("create database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close database: %v", err)
		}
	})

	if err := Run(db); err != nil {
		t.Fatalf("run migration: %v", err)
	}

	for _, table := range []string{
		"sys_user",
		"sys_user_roles",
		"roles",
		"role_menus",
		"menus",
		"token_black_lists",
	} {
		if !db.DB.Migrator().HasTable(table) {
			t.Errorf("table %q was not created", table)
		}
	}
}
