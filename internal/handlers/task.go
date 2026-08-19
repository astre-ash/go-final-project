package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"go-final-progect/internal/models"
	"go-final-progect/internal/nextdate"
)

const dateFormat = "20060102"

func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var t models.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Printf("AddTask error: failed to decode JSON: %v", err)
		sendError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		sendError(w, "Task title is required", http.StatusBadRequest)
		return
	}

	if err := setTaskDate(&t); err != nil {
		log.Printf("AddTask error: invalid date or repeat rule: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.Store.AddTask(t)
	if err != nil {
		log.Printf("AddTask error: failed to insert task into DB: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"id": strconv.FormatInt(id, 10),
	}
	sendJSON(w, response, http.StatusOK)

}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		log.Printf("GetTask error: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := h.Store.GetTaskById(id)
	if err != nil {
		log.Printf("GetTask error: task not found (ID: %d): %v", id, err)
		sendError(w, "Task not found", http.StatusNotFound)
		return
	}

	sendJSON(w, t, http.StatusOK)

}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var t models.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Printf("UpdateTask error: failed to decode JSON: %v", err)
		sendError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if t.ID == "" {
		sendError(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(t.ID)
	if err != nil || id <= 0 {
		sendError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		sendError(w, "Task title is required", http.StatusBadRequest)
		return
	}

	if err := setTaskDate(&t); err != nil {
		log.Printf("UpdateTask error: invalid date or repeat rule: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Store.UpdateTask(t)
	if err != nil {
		log.Printf("UpdateTask error: failed to update task (ID: %s): %v", t.ID, err)

		if err.Error() == "task not found" {
			sendError(w, "Task not found", http.StatusNotFound)
		} else {
			sendError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	sendJSON(w, map[string]any{}, http.StatusOK)
}

func (h *TaskHandler) TaskDone(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		log.Printf("TaskDone error: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := h.Store.GetTaskById(id)
	if err != nil {
		log.Printf("TaskDone error: task not found (ID: %d): %v", id, err)
		sendError(w, "Task not found", http.StatusNotFound)
		return
	}

	switch {
	case t.Repeat != "":

		now := time.Now()
		next, err := nextdate.NextDate(now, t.Date, t.Repeat)
		if err != nil {
			log.Printf("TaskDone error: failed to calculate next date: %v", err)
			sendError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		t.Date = next
		err = h.Store.UpdateTask(t)
		if err != nil {
			log.Printf("TaskDone error: failed to update task in DB: %v", err)
			sendError(w, "Internal server error", http.StatusInternalServerError)
			return
		}

	default:
		err = h.Store.DeleteTask(id)
		if err != nil {
			log.Printf("TaskDone error: failed to delete task in DB: %v", err)
			sendError(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	sendJSON(w, map[string]any{}, http.StatusOK)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r)
	if err != nil {
		log.Printf("DeleteTask error: %v", err)
		sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.Store.DeleteTask(id)
	if err != nil {
		log.Printf("DeleteTask error: failed to delete task (ID: %d): %v", id, err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]any{}, http.StatusOK)
}

func setTaskDate(t *models.Task) error {
	now := time.Now()
	today := now.Format(dateFormat)

	switch {
	case t.Date != "":
		_, err := time.Parse(dateFormat, t.Date)
		if err != nil {
			return fmt.Errorf("invalid date format, expected %s", dateFormat)

		}
	default:
		t.Date = today
	}

	switch {
	case t.Repeat != "":
		next, err := nextdate.NextDate(now, t.Date, t.Repeat)
		if err != nil {
			return fmt.Errorf("invalid repeat rule: %w", err)
		}

		if t.Date < today {
			t.Date = next
		}
	case t.Date < today:
		t.Date = today
	}
	return nil
}

func extractID(r *http.Request) (int, error) {
	idStr := r.FormValue("id")
	if idStr == "" {
		return 0, fmt.Errorf("task ID is required")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid task ID format")
	}

	return id, nil
}
