//go:build integration

package user

import (
	"os"
	"testing"
)

func TestPostgresUserRepository_CRUD(t *testing.T) {
	if os.Getenv("DATABASE_URL") == "" && os.Getenv("DB_HOST") == "" {
		t.Skip("set DATABASE_URL (or DB_* env vars) to run integration tests")
	}

	repo, err := NewPostgresUserRepositoryFromEnv()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	created, err := repo.Create(User{Name: "Joao", Email: "joao@example.com"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == 0 {
		t.Fatalf("expected non-zero id")
	}

	got, err := repo.GetByID(created.ID)
	if err != nil {
		t.Fatalf("get by id: %v", err)
	}
	if got.Email != created.Email {
		t.Fatalf("expected email %q, got %q", created.Email, got.Email)
	}

	updated, err := repo.Update(created.ID, User{Name: "Joao Silva", Email: "js@example.com"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Name != "Joao Silva" {
		t.Fatalf("expected updated name")
	}

	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("get all: %v", err)
	}
	if len(all) == 0 {
		t.Fatalf("expected at least 1 user")
	}

	if err := repo.Delete(created.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	_, err = repo.GetByID(created.ID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
}
