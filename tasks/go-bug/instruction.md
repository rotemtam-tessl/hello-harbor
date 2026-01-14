# Task: Fix the Todo App Bug

The `/workspace` directory contains a Go module with a Todo application using GORM and SQLite.

The app has CRUD endpoints for managing todos:
- `GET /todos` - List all todos
- `POST /todos` - Create a new todo
- `PUT /todos/{id}/complete` - Mark a todo as completed
- `DELETE /todos/{id}` - Delete a todo

There is a bug in the application that causes one of the tests to fail. Your task is to find and fix the bug so that all tests pass.

## Requirements

- All tests must pass when running `go test ./...`
- Do not modify the test file