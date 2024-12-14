package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getAllTasks(res http.ResponseWriter, req *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(resp)

	return
}

func addTask(res http.ResponseWriter, req *http.Request) {
	var task Task
	err := json.NewDecoder(req.Body).Decode(&task)
	if err != nil {
		http.Error(res, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if _, exists := tasks[task.ID]; exists {
		http.Error(res, "Task with this ID already exists", http.StatusConflict)
		return
	}

	tasks[task.ID] = task

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
}

func getTaskItem(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	task, exists := tasks[id]
	if !exists {
		http.Error(res, "Task not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	res.Write(resp)
}

func removeTaskItem(res http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	_, exists := tasks[id]
	if !exists {
		http.Error(res, "Task not found", http.StatusNotFound)
	}
	delete(tasks, id)
	res.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", addTask)
	r.Get("/tasks/{id}", getTaskItem)
	r.Delete("/tasks/{id}", removeTaskItem)

	fmt.Println("Сервер стартанул на хосте http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
