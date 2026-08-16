package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nova/auth/internal/model"
	"github.com/nova/auth/internal/repository"
	"github.com/nova/auth/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	users map[string]*model.User
	next  uint
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{users: map[string]*model.User{}, next: 1}
}

func (r *fakeUserRepository) Create(user *model.User) error {
	if user.ID == 0 {
		user.ID = r.next
		r.next++
	}
	copied := *user
	r.users[user.Username] = &copied
	return nil
}

func (r *fakeUserRepository) GetByUsername(username string) (*model.User, error) {
	user, ok := r.users[username]
	if !ok {
		return nil, repository.ErrNotFound
	}
	copied := *user
	return &copied, nil
}

func (r *fakeUserRepository) GetByID(id uint) (*model.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			copied := *user
			return &copied, nil
		}
	}
	return nil, repository.ErrNotFound
}

func TestLoginUsesResponseEnvelope(t *testing.T) {
	app, _ := testApp(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"username":"admin","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["message"] != "login successful" {
		t.Fatalf("unexpected message: %v", body["message"])
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object: %+v", body)
	}
	if data["token"] == "" || data["token_type"] != "Bearer" {
		t.Fatalf("unexpected token data: %+v", data)
	}
}

func TestCreateUserStatusCodes(t *testing.T) {
	app, adminToken := testApp(t)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/users", strings.NewReader(`{"username":"","password":"secret"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected blank user request to return 400, got %d", resp.StatusCode)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/users", strings.NewReader(`{"username":"admin","password":"secret"}`))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("expected duplicate user request to return 409, got %d", resp.StatusCode)
	}
}

func TestValidateRequiresAuthorization(t *testing.T) {
	app, _ := testApp(t)

	resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/v1/auth/validate", nil))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}

	var body map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["message"] != "authentication required" || body["error"] != "missing authorization header" {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func testApp(t *testing.T) (*fiber.App, string) {
	t.Helper()

	repo := newFakeUserRepository()
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(&model.User{Username: "admin", Password: string(hash), Role: model.RoleAdmin}); err != nil {
		t.Fatal(err)
	}

	svc := service.NewService(repo, "test-secret", "admin", "secret")
	auth, err := svc.Authenticate("admin", "secret")
	if err != nil {
		t.Fatal(err)
	}

	app := fiber.New()
	NewHandler(svc, "test-secret").RegisterRoutes(app)
	return app, auth.Token
}
