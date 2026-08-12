package handlers

import (
	"go-final-progect/internal/nextdate"
	"log"
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeatStr := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(nextdate.DateFormat, nowStr)
		if err != nil {
			log.Printf("NextDateHandler error: invalid 'now' parameter format ('%s'): %v", nowStr, err)
			http.Error(w, "Invalid format for 'now' parameter", http.StatusBadRequest)
			return
		}
	}
	next, err := nextdate.NextDate(now, dateStr, repeatStr)
	if err != nil {
		log.Printf("NextDateHandler error: failed to calculate next date: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}
