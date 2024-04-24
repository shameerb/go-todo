package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	port = flag.String("port", "8000", "http server port")
)

type TodoServer struct {
	port string
	db   *gorm.DB
}

type Todo struct {
	gorm.Model
	Description string
	Completed   bool
}

type TodoCreateRequest struct {
	Description string
}

func NewTodoServer(port string) *TodoServer {
	return &TodoServer{
		port: port,
	}
}

// Repository
func (t *TodoServer) setupDb() error {
	dbName := os.Getenv("DB_FILE")
	if len(dbName) == 0 {
		dbName = "test.db"
	}
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Println("failed to connec to database sqlite")
		return err
	}
	t.db = db
	return t.db.Debug().AutoMigrate(&Todo{})
}

func (t *TodoServer) setupHttp() error {
	router := mux.NewRouter()
	router.HandleFunc("/health", t.checkHealth).Methods("GET")
	router.HandleFunc("/todo-completed", t.getCompleted).Methods("GET")
	router.HandleFunc("/todo-pending", t.getPending).Methods("GET")
	router.HandleFunc("/todo", t.createTodo).Methods("PUT")
	router.HandleFunc("/todo/{id}", t.updateTodo).Methods("POST")
	router.HandleFunc("/todo/{id}", t.deleteTodo).Methods("DELETE")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler(router)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", t.port), handler); err != nil {
		log.Panicf("failed to create http server: %s", err)
	}
	return nil
}

func (t *TodoServer) getTodoItemsQuery(completed bool) []Todo {
	var todos []Todo
	t.db.Where("Completed = ?", completed).Find(&todos)
	return todos
}

func (t *TodoServer) createTodoQuery(todo *Todo) error {
	result := t.db.Create(&todo)
	return result.Error
}

func (t *TodoServer) getTodoItem(id uint) (*Todo, error) {
	todo := &Todo{}
	result := t.db.First(&todo)
	if result.Error != nil {
		log.Warnf("todo item not found in database: %d", id)
		return nil, result.Error
	}
	return todo, nil
}

func (t *TodoServer) updateTodoQuery(todo *Todo) error {
	result := t.db.Save(&todo)
	return result.Error
}

func (t *TodoServer) deleteTodoQuery(todo *Todo) error {
	result := t.db.Delete(&todo)
	return result.Error
}

// Services
func (t *TodoServer) createTodo(w http.ResponseWriter, r *http.Request) {
	var todoRequest TodoCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&todoRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo := &Todo{Description: todoRequest.Description, Completed: false}
	if err := t.createTodoQuery(todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (t *TodoServer) getCompleted(w http.ResponseWriter, r *http.Request) {
	completedItems := t.getTodoItemsQuery(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedItems)
}

func (t *TodoServer) getPending(w http.ResponseWriter, r *http.Request) {
	pendingItems := t.getTodoItemsQuery(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pendingItems)
}

func (t *TodoServer) updateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	todo, err := t.getTodoItem(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.Completed = !todo.Completed
	if err := t.updateTodoQuery(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (t *TodoServer) deleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	todo, err := t.getTodoItem(uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := t.deleteTodoQuery(todo); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{'deleted'}: true"))
}

func (t *TodoServer) checkHealth(w http.ResponseWriter, r *http.Request) {
	log.Info("Health is OK")
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{'alive': true}"))
}

func (t *TodoServer) Start() error {
	if err := t.setupDb(); err != nil {
		return nil
	}
	return t.setupHttp()
}

func main() {
	flag.Parse()
	t := NewTodoServer(*port)
	if err := t.Start(); err != nil {
		log.Fatal(err)
	}
}
