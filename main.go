package main

import (
	"encoding/json" // for JSON encode/decode
	"fmt"           // for printing logs to terminal
	"net/http"      // for HTTP server & handlers
	"strconv"       // for string -> int conversion
	"sync"          // for mutex (concurrency safety)
)

// Todo represents a single todo item (response structure)
type Todo struct {
	ID    int    `json:"id"`    // unique identifier
	Title string `json:"title"` // task description
	Done  bool   `json:"done"`  // completion status
}

// CreateTodoRequest represents input body for creating todo
type CreateTodoRequest struct {
	Title string `json:"title"`
}

// shared in-memory storage
var todos = make(map[int]Todo) // stores todos as id -> Todo
var mu sync.Mutex              // mutex to protect todos map
var nextID = 1                 // auto-incrementing id


// get all todos
func getTodosHandler(w http.ResponseWriter, r *http.Request) {

	// tell client that response is JSON
	w.Header().Set("Content-Type", "application/json")

	// lock before accessing shared map
	mu.Lock()
	defer mu.Unlock()

	// encode todos map as JSON and send response
	json.NewEncoder(w).Encode(todos)
}


// get
func createTodoHandler(w http.ResponseWriter, r *http.Request) {

	// if request is not POST, return 405
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// since we returning JSON
	w.Header().Set("Content-Type", "application/json")

	// err handling for decoding request body (bad input)
	var req CreateTodoRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// lock coz concurrent access to shared resource
	mu.Lock()
	defer mu.Unlock()

	// create new todo object
	todo := Todo{
		ID:    nextID,
		Title: req.Title,
		Done:  false,
	}

	// store todo in map and increment ID
	todos[nextID] = todo
	nextID++

	// convert todo to JSON and send response
	json.NewEncoder(w).Encode(todo)
}


// put update
func updateTodoHandler(w http.ResponseWriter, r *http.Request) {

	// allow only PUT method
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// response will be JSON
	w.Header().Set("Content-Type", "application/json")

	// read id from query param (?id=1)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert id from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// lock shared data before modifying
	mu.Lock()
	defer mu.Unlock()

	// check if todo exists
	todo, exists := todos[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// update todo status
	todo.Done = true
	todos[id] = todo

	// return updated todo
	json.NewEncoder(w).Encode(todo)
}


// delete
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {

	// allow only DELETE method
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// read id from query param
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// convert id to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// lock before deleting from map
	mu.Lock()
	defer mu.Unlock()

	// check existence
	if _, exists := todos[id]; !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// delete todo
	delete(todos, id)

	// 204 = success with no response body
	w.WriteHeader(http.StatusNoContent)
}


func main() {

	// route registrations
	http.HandleFunc("/todos", getTodosHandler)
	http.HandleFunc("/todos/create", createTodoHandler)
	http.HandleFunc("/todos/update", updateTodoHandler)
	http.HandleFunc("/todos/delete", deleteTodoHandler)

	fmt.Println("Server started on port 8080")

	// start HTTP server using default router
	http.ListenAndServe(":8080", nil)
}
