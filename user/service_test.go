package user

import "testing"

type fakeRepo struct {
	users []User
}

func (f *fakeRepo) GetAll() ([]User, error)              { return f.users, nil }
func (f *fakeRepo) GetByID(id int) (*User, error)        { return &f.users[0], nil }
func (f *fakeRepo) Create(u User) (*User, error)         { return &u, nil }
func (f *fakeRepo) Update(id int, u User) (*User, error) { u.ID = id; return &u, nil }
func (f *fakeRepo) Delete(id int) error                  { return nil }

func TestUserService_ListUsers(t *testing.T) {
	svc := NewUserService(&fakeRepo{users: []User{{ID: 1, Name: "Ana", Email: "ana@example.com"}}})

	users, err := svc.ListUsers()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
	if users[0].ID != 1 {
		t.Fatalf("expected id=1, got %d", users[0].ID)
	}
}
