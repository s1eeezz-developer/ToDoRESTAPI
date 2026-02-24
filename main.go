package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

var (
	tasks  = []Task{}
	nextID = 1
	mu     sync.Mutex
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен. Используйте POST", http.StatusMethodNotAllowed)
		return
	}

	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Ошибка в формате JSON", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)
	mu.Unlock()

	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTask)
}

func printHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен. Используйте GET", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, tasks)

}

func main() {
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/tasks", printHandler)

	fmt.Println("Севрер запущен на портах :8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err)
	}
}
