package handlers

import (
	"log"
	"net/http"

	"go-final-progect/internal/models"
)

const limit = 50

func (h *TaskHandler) TasksList(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")

	tasks, err := h.Store.GetTasks(search, limit)
	if err != nil {
		log.Printf("TasksList error: failed to fetch tasks from DB: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = make([]models.Task, 0)
	}

	response := map[string]any{
		"tasks": tasks,
	}

	sendJSON(w, response, http.StatusOK)
}
