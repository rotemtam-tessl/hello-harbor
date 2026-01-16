package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todo/models"
)

func setupTestApp(t *testing.T) *App {
	// Create a temp file for test database
	tmpFile, err := os.CreateTemp("", "test-*.db")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	dbPath := tmpFile.Name()
	t.Cleanup(func() { os.Remove(dbPath) })

	// Run Atlas migrations using embedded migrations from main package
	db, err := models.InitDB(dbPath, getMigrations())
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}

	return NewApp(db)
}

// --- Todo Tests ---

func TestCreateTodo(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	body := []byte(`{"title":"Buy groceries"}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var todo models.Todo
	json.NewDecoder(w.Body).Decode(&todo)
	if todo.Title != "Buy groceries" {
		t.Errorf("expected title %q, got %q", "Buy groceries", todo.Title)
	}
	if todo.Completed != false {
		t.Errorf("expected completed to be false, got true")
	}
}

func TestListTodos(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Create a todo first
	app.DB.Create(&models.Todo{Title: "Test todo"})

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var todos []models.Todo
	json.NewDecoder(w.Body).Decode(&todos)
	if len(todos) != 1 {
		t.Errorf("expected 1 todo, got %d", len(todos))
	}
}

func TestCompleteTodo(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Create a todo
	todo := models.Todo{Title: "Complete me"}
	app.DB.Create(&todo)

	req := httptest.NewRequest(http.MethodPut, "/todos/1/complete", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify the todo is marked as completed in the database
	var updatedTodo models.Todo
	app.DB.First(&updatedTodo, 1)
	if !updatedTodo.Completed {
		t.Errorf("expected todo to be completed, but it was not")
	}
}

func TestDeleteTodo(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Create a todo
	app.DB.Create(&models.Todo{Title: "Delete me"})

	req := httptest.NewRequest(http.MethodDelete, "/todos/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify the todo is deleted
	var count int64
	app.DB.Model(&models.Todo{}).Count(&count)
	if count != 0 {
		t.Errorf("expected 0 todos, got %d", count)
	}
}

// --- Category Tests ---

func TestListCategories_SeededData(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var categories []models.Category
	json.NewDecoder(w.Body).Decode(&categories)

	// Should have 3 seeded categories: Home, Work, Research
	if len(categories) != 3 {
		t.Errorf("expected 3 seeded categories, got %d", len(categories))
	}

	expectedNames := map[string]bool{"Home": true, "Work": true, "Research": true}
	for _, cat := range categories {
		if !expectedNames[cat.Name] {
			t.Errorf("unexpected category name: %s", cat.Name)
		}
	}
}

func TestCreateCategory(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	body := []byte(`{"name":"Personal"}`)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var category models.Category
	json.NewDecoder(w.Body).Decode(&category)
	if category.Name != "Personal" {
		t.Errorf("expected name %q, got %q", "Personal", category.Name)
	}
}

func TestGetCategory(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Get one of the seeded categories (ID 1 = Home)
	req := httptest.NewRequest(http.MethodGet, "/categories/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var category models.Category
	json.NewDecoder(w.Body).Decode(&category)
	if category.Name != "Home" {
		t.Errorf("expected name %q, got %q", "Home", category.Name)
	}
}

func TestUpdateCategory(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Update the first seeded category
	body := []byte(`{"name":"Updated Home"}`)
	req := httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify the category is updated in the database
	var updatedCategory models.Category
	app.DB.First(&updatedCategory, 1)
	if updatedCategory.Name != "Updated Home" {
		t.Errorf("expected name %q, got %q", "Updated Home", updatedCategory.Name)
	}
}

func TestDeleteCategory(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Create a new category to delete (don't delete seeded ones)
	app.DB.Create(&models.Category{Name: "ToDelete"})

	// Get count before delete
	var countBefore int64
	app.DB.Model(&models.Category{}).Count(&countBefore)

	req := httptest.NewRequest(http.MethodDelete, "/categories/4", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
	}

	// Verify the category is deleted
	var countAfter int64
	app.DB.Model(&models.Category{}).Count(&countAfter)
	if countAfter != countBefore-1 {
		t.Errorf("expected %d categories after delete, got %d", countBefore-1, countAfter)
	}
}

func TestCreateTodoWithCategory(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Create a todo with category_id (Work = 2)
	categoryID := uint(2)
	body := []byte(`{"title":"Finish report","category_id":2}`)
	req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var todo models.Todo
	json.NewDecoder(w.Body).Decode(&todo)
	if todo.CategoryID == nil || *todo.CategoryID != categoryID {
		t.Errorf("expected category_id %d, got %v", categoryID, todo.CategoryID)
	}
}

func TestCategoryWithDescription(t *testing.T) {
	app := setupTestApp(t)
	router := app.SetupRouter()

	// Create a category with description
	body := []byte(`{"name":"Shopping","description":"Things to buy"}`)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var category models.Category
	json.NewDecoder(w.Body).Decode(&category)

	if category.Name != "Shopping" {
		t.Errorf("expected name %q, got %q", "Shopping", category.Name)
	}
	if category.Description != "Things to buy" {
		t.Errorf("expected description %q, got %q", "Things to buy", category.Description)
	}

	// Verify it persists by fetching it
	var savedCategory models.Category
	app.DB.First(&savedCategory, category.ID)
	if savedCategory.Description != "Things to buy" {
		t.Errorf("expected saved description %q, got %q", "Things to buy", savedCategory.Description)
	}
}
