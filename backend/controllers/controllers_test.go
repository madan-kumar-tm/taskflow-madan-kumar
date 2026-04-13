package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"taskflow/config"
	"taskflow/routes"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	config.DB = db

	db.Exec(`CREATE TABLE users (id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))), name TEXT, email TEXT UNIQUE, password TEXT, created_at TEXT);`)
	db.Exec(`CREATE TABLE projects (id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))), name TEXT, description TEXT, owner_id TEXT, created_at TEXT);`)
	db.Exec(`CREATE TABLE tasks (id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))), title TEXT, description TEXT, status TEXT, priority TEXT, project_id TEXT, assignee_id TEXT, creator_id TEXT, due_date TEXT, created_at TEXT, updated_at TEXT);`)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	routes.AuthRoutes(r)
	routes.ProjectRoutes(r)
	routes.UserRoutes(r)
	return r
}

type authResponse struct {
	Token string `json:"token"`
}

type projectResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type taskResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ProjectID string `json:"project_id"`
	Status    string `json:"status"`
	Priority  string `json:"priority"`
}

func mustRegisterUser(t *testing.T, router *gin.Engine, name, email, password string) {
	w := httptest.NewRecorder()
	reqBody := `{"name":"` + name + `","email":"` + email + `","password":"` + password + `"}`
	req, _ := http.NewRequest("POST", "/auth/register", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("register failed: expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func mustLoginUser(t *testing.T, router *gin.Engine, email, password string) string {
	w := httptest.NewRecorder()
	reqBody := `{"email":"` + email + `","password":"` + password + `"}`
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login failed: expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var resp authResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse login response: %v", err)
	}
	if resp.Token == "" {
		t.Fatal("expected token in login response")
	}
	return resp.Token
}

func authRequest(method, path, body, token string) *http.Request {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func TestRegister(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	w := httptest.NewRecorder()
	req := authRequest("POST", "/auth/register", `{"name":"test","email":"test@example.com","password":"secret"}`, "")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestLoginWrongPassword(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	mustRegisterUser(t, router, "user", "user@example.com", "secret")

	w := httptest.NewRecorder()
	req := authRequest("POST", "/auth/login", `{"email":"user@example.com","password":"bad"}`, "")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateProjectRequiresAuth(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	w := httptest.NewRecorder()
	req := authRequest("POST", "/projects", `{"name":"My Project"}`, "")

	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateProjectAndListProjects(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	mustRegisterUser(t, router, "owner", "owner@example.com", "secret")
	token := mustLoginUser(t, router, "owner@example.com", "secret")

	w := httptest.NewRecorder()
	req := authRequest("POST", "/projects", `{"name":"Test Project","description":"A project"}`, token)

	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	var project projectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &project); err != nil {
		t.Fatalf("failed to parse project response: %v", err)
	}
	if project.Name != "Test Project" {
		t.Fatalf("expected project name to be Test Project, got %q", project.Name)
	}

	w = httptest.NewRecorder()
	req = authRequest("GET", "/projects", ``, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	var projects []projectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &projects); err != nil {
		t.Fatalf("failed to parse projects list: %v", err)
	}
	if len(projects) != 1 || projects[0].ID != project.ID {
		t.Fatalf("expected one project with id %s, got %#v", project.ID, projects)
	}
}

func TestCreateTaskAndListTasks(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	mustRegisterUser(t, router, "owner", "owner@example.com", "secret")
	token := mustLoginUser(t, router, "owner@example.com", "secret")

	w := httptest.NewRecorder()
	req := authRequest("POST", "/projects", `{"name":"Task Project","description":"Testing tasks"}`, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating project, got %d, body: %s", w.Code, w.Body.String())
	}

	var project projectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &project); err != nil {
		t.Fatalf("failed to parse project response: %v", err)
	}

	w = httptest.NewRecorder()
	req = authRequest("POST", "/projects/"+project.ID+"/tasks", `{"title":"A Task","description":"Task desc","priority":"high"}`, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating task, got %d, body: %s", w.Code, w.Body.String())
	}

	var task taskResponse
	if err := json.Unmarshal(w.Body.Bytes(), &task); err != nil {
		t.Fatalf("failed to parse task response: %v", err)
	}
	if task.Title != "A Task" || task.ProjectID != project.ID {
		t.Fatalf("unexpected task response: %#v", task)
	}

	w = httptest.NewRecorder()
	req = authRequest("GET", "/projects/"+project.ID+"/tasks", ``, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 listing tasks, got %d, body: %s", w.Code, w.Body.String())
	}

	var tasks []taskResponse
	if err := json.Unmarshal(w.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("failed to parse task list: %v", err)
	}
	if len(tasks) != 1 || tasks[0].ID != task.ID {
		t.Fatalf("expected one task with id %s, got %#v", task.ID, tasks)
	}
}

func TestUpdateAndDeleteTask(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	mustRegisterUser(t, router, "owner", "owner@example.com", "secret")
	token := mustLoginUser(t, router, "owner@example.com", "secret")

	w := httptest.NewRecorder()
	req := authRequest("POST", "/projects", `{"name":"Update Project","description":"Project for update"}`, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating project, got %d, body: %s", w.Code, w.Body.String())
	}

	var project projectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &project); err != nil {
		t.Fatalf("failed to parse project response: %v", err)
	}

	w = httptest.NewRecorder()
	req = authRequest("POST", "/projects/"+project.ID+"/tasks", `{"title":"Delete Task","description":"Task to remove"}`, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating task, got %d, body: %s", w.Code, w.Body.String())
	}

	var task taskResponse
	if err := json.Unmarshal(w.Body.Bytes(), &task); err != nil {
		t.Fatalf("failed to parse task response: %v", err)
	}

	w = httptest.NewRecorder()
	req = authRequest("PATCH", "/tasks/"+task.ID, `{"status":"completed"}`, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 updating task, got %d, body: %s", w.Code, w.Body.String())
	}

	var updated taskResponse
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("failed to parse updated task response: %v", err)
	}
	if updated.Status != "done" {
		t.Fatalf("expected status done, got %q", updated.Status)
	}

	w = httptest.NewRecorder()
	req = authRequest("DELETE", "/tasks/"+task.ID, ``, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 deleting task, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateTaskInvalidAssigneeValidation(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()
	mustRegisterUser(t, router, "owner", "owner@example.com", "secret")
	token := mustLoginUser(t, router, "owner@example.com", "secret")

	w := httptest.NewRecorder()
	req := authRequest("POST", "/projects", `{"name":"Validation Project","description":"Project for assignee validation"}`, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201 creating project, got %d, body: %s", w.Code, w.Body.String())
	}

	var project projectResponse
	if err := json.Unmarshal(w.Body.Bytes(), &project); err != nil {
		t.Fatalf("failed to parse project response: %v", err)
	}

	w = httptest.NewRecorder()
	req = authRequest("POST", "/projects/"+project.ID+"/tasks", `{"title":"Bad assignee","assignee_id":"1"}`, token)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid assignee_id, got %d, body: %s", w.Code, w.Body.String())
	}

	var payload map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse validation response: %v", err)
	}
	if payload["error"] != "validation failed" {
		t.Fatalf("expected validation failed error, got %#v", payload["error"])
	}
	fields, ok := payload["fields"].(map[string]any)
	if !ok {
		t.Fatalf("expected fields object, got %#v", payload["fields"])
	}
	if fields["assignee_id"] != "must be a valid uuid" {
		t.Fatalf("expected assignee_id validation message, got %#v", fields["assignee_id"])
	}
}
