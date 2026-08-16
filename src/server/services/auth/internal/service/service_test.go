package service

import (
	"errors"
	"testing"

	"github.com/nova/auth/internal/model"
	"github.com/nova/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	users map[string]*model.User
	err   error
	next  uint
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{users: map[string]*model.User{}, next: 1}
}

func (r *fakeUserRepository) Create(user *model.User) error {
	if r.err != nil {
		return r.err
	}
	if user.ID == 0 {
		user.ID = r.next
		r.next++
	}
	copied := *user
	r.users[user.Username] = &copied
	return nil
}

func (r *fakeUserRepository) GetByUsername(username string) (*model.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	user, ok := r.users[username]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copied := *user
	return &copied, nil
}

func (r *fakeUserRepository) GetByID(id uint) (*model.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	for _, user := range r.users {
		if user.ID == id {
			copied := *user
			return &copied, nil
		}
	}
	return nil, repository.ErrNotFound
}

func TestAuthenticateReturnsTokenResponse(t *testing.T) {
	repo := newFakeUserRepository()
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(&model.User{Username: "admin", Password: string(hash), Role: model.RoleAdmin}); err != nil {
		t.Fatal(err)
	}

	svc := NewService(repo, "test-secret", "admin", "secret")
	auth, err := svc.Authenticate(" admin ", "secret")
	if err != nil {
		t.Fatal(err)
	}
	if auth.Token == "" {
		t.Fatal("expected token")
	}
	if auth.TokenType != "Bearer" {
		t.Fatalf("expected Bearer token type, got %q", auth.TokenType)
	}
	if auth.User.Username != "admin" || auth.User.Role != model.RoleAdmin {
		t.Fatalf("unexpected user in auth response: %+v", auth.User)
	}

	claims, err := svc.ValidateToken(auth.Token)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Username != "admin" || claims.Role != model.RoleAdmin {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestCreateUserMapsValidationAndDuplicateErrors(t *testing.T) {
	repo := newFakeUserRepository()
	svc := NewService(repo, "test-secret", "admin", "secret")

	if _, err := svc.CreateUser(" ", "secret"); !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}

	if _, err := svc.CreateUser("alice", "secret"); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateUser("alice", "secret"); !errors.Is(err, ErrUserExists) {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestCreateUserReturnsRepositoryErrors(t *testing.T) {
	repo := newFakeUserRepository()
	repo.err = errors.New("database unavailable")
	svc := NewService(repo, "test-secret", "admin", "secret")

	if _, err := svc.CreateUser("alice", "secret"); err == nil || errors.Is(err, ErrUserExists) {
		t.Fatalf("expected repository error, got %v", err)
	}
}
