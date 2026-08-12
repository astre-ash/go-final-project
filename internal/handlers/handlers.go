package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go-final-progect/internal/storage/sqlite"
)

type TaskHandler struct {
	Store *sqlite.Storage
}

func RegisterRoutes(mux *http.ServeMux, store *sqlite.Storage) {
	h := &TaskHandler{
		Store: store,
	}
	mux.HandleFunc("POST /api/signin", h.Signin)

	mux.HandleFunc("GET /api/nextdate", NextDateHandler)
	mux.HandleFunc("GET /api/tasks", AuthMiddleware(h.TasksList))

	mux.HandleFunc("POST /api/task", AuthMiddleware(h.AddTask))
	mux.HandleFunc("GET /api/task", AuthMiddleware(h.GetTask))
	mux.HandleFunc("PUT /api/task", AuthMiddleware(h.UpdateTask))
	mux.HandleFunc("DELETE /api/task", AuthMiddleware(h.DeleteTask))

	mux.HandleFunc("POST /api/task/done", AuthMiddleware(h.TaskDone))

}

func sendJSON(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendError(w http.ResponseWriter, msg string, status int) {
	log.Printf("ERROR [%d]: %s", status, msg)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
