package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"todo/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

//go:embed migrations/*
var migrationsFS embed.FS

// getMigrations returns the migrations subdirectory from the embedded FS
func getMigrations() fs.FS {
	sub, _ := fs.Sub(migrationsFS, "migrations")
	return sub
}

// App holds the application dependencies
type App struct {
	DB *gorm.DB
}

// NewApp creates a new App with database connection
func NewApp(db *gorm.DB) *App {
	return &App{DB: db}
}

// --- Todo Handlers ---

// ListTodos returns all todos
func (a *App) ListTodos(w http.ResponseWriter, r *http.Request) {
	var todos []models.Todo
	a.DB.Preload("Category").Find(&todos)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo creates a new todo
func (a *App) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.DB.Create(&todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// CompleteTodo marks a todo as completed
// BUG: This function has a bug that prevents todos from being marked complete
func (a *App) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var todo models.Todo
	if err := a.DB.First(&todo, id).Error; err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	a.DB.Model(&todo).Update("completed", true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo deletes a todo
func (a *App) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := a.DB.Delete(&models.Todo{}, id).Error; err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Category Handlers ---

// ListCategories returns all categories
func (a *App) ListCategories(w http.ResponseWriter, r *http.Request) {
	var categories []models.Category
	a.DB.Find(&categories)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// GetCategory returns a single category by ID
func (a *App) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := a.DB.First(&category, id).Error; err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// CreateCategory creates a new category
func (a *App) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	a.DB.Create(&category)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(category)
}

// UpdateCategory updates an existing category
func (a *App) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var category models.Category
	if err := a.DB.First(&category, id).Error; err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	var input models.Category
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	a.DB.Model(&category).Update("name", input.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

// DeleteCategory deletes a category
func (a *App) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := a.DB.Delete(&models.Category{}, id).Error; err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SetupRouter configures the HTTP routes
func (a *App) SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Todo routes
	r.HandleFunc("/todos", a.ListTodos).Methods("GET")
	r.HandleFunc("/todos", a.CreateTodo).Methods("POST")
	r.HandleFunc("/todos/{id}/complete", a.CompleteTodo).Methods("PUT")
	r.HandleFunc("/todos/{id}", a.DeleteTodo).Methods("DELETE")

	// Category routes
	r.HandleFunc("/categories", a.ListCategories).Methods("GET")
	r.HandleFunc("/categories", a.CreateCategory).Methods("POST")
	r.HandleFunc("/categories/{id}", a.GetCategory).Methods("GET")
	r.HandleFunc("/categories/{id}", a.UpdateCategory).Methods("PUT")
	r.HandleFunc("/categories/{id}", a.DeleteCategory).Methods("DELETE")

	return r
}

func main() {
	db, err := models.InitDB("todos.db", getMigrations())
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	app := NewApp(db)
	router := app.SetupRouter()

	log.Println("Server starting on :8080")
	http.ListenAndServe(":8080", router)
}
