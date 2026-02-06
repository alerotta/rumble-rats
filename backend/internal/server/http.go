package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type App struct {
	DB *sql.DB
}

func NewApp(db *sql.DB) *App {
	return &App{
		DB: db,
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health" , a.handleHealth)
	mux.HandleFunc("/matchmaking", a.handleMatchmaking)
	mux.HandleFunc("/health/db", a.handleHealthDB)

	return devCORS(requestLogger(mux))
}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, JSON{"error": "method not allowed"})
		return
	}
	writeJSON(w, http.StatusOK, JSON{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (a *App) handleMatchmaking(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, JSON{"error": "method not allowed"})
		return
	}

	var body struct {
		Mode string `json:"mode"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	writeJSON(w, http.StatusOK, JSON{
		"status": "queued",
		"mode":   body.Mode,
	})

	
}

func (a *App) handleHealthDB( w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet {
	writeJSON(w, http.StatusMethodNotAllowed, JSON{"error": "method not allowed"})
	return
	}


	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var now time.Time
	err := a.DB.QueryRowContext(ctx, "SELECT now()").Scan(&now)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, JSON{
			"status": "error",
			"error":  "database not reachable",
		})
		return
	}

	writeJSON(w, http.StatusOK, JSON{
		"status": "ok",
		"db":     "connected",
		"time":   now.UTC().Format(time.RFC3339),
	})

}

// helper functions

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}