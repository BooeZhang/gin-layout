package infra

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type crudTestModel struct {
	ID   int64 `gorm:"primaryKey"`
	Name string
}

func TestCRUDRepositoryFindByID(t *testing.T) {
	t.Parallel()

	db, err := gorm.Open(sqlite.Open("file:crud-find-by-id?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&crudTestModel{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	repo := NewModelCRUDRepository[crudTestModel, int64](db)
	ctx := context.Background()

	entity := &crudTestModel{Name: "found"}
	if err := repo.Create(ctx, entity); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, found, err := repo.FindByID(ctx, entity.ID)
	if err != nil {
		t.Fatalf("find existing: %v", err)
	}
	if !found {
		t.Fatal("expected existing entity to be found")
	}
	if got == nil || got.ID != entity.ID || got.Name != entity.Name {
		t.Fatalf("FindByID() = %#v, want %#v", got, entity)
	}

	got, found, err = repo.FindByID(ctx, entity.ID+1)
	if err != nil {
		t.Fatalf("find missing: %v", err)
	}
	if found {
		t.Fatal("expected missing entity not to be found")
	}
	if got != nil {
		t.Fatalf("expected nil entity for missing row, got %#v", got)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql database: %v", err)
	}
	if err := sqlDB.Close(); err != nil {
		t.Fatalf("close sql database: %v", err)
	}

	got, found, err = repo.FindByID(ctx, entity.ID)
	if err == nil {
		t.Fatal("expected database error after closing connection")
	}
	if found {
		t.Fatal("expected found=false when the lookup fails")
	}
	if got != nil {
		t.Fatalf("expected nil entity when the lookup fails, got %#v", got)
	}
}
