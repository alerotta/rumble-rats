package server

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/alerotta/rumble-rats/backend/internal/auth"
)

type App struct {
	DB *sql.DB
	jwtSecret string
}

func NewApp(db *sql.DB, JWT string) *App {
	return &App{
		DB: db,
		jwtSecret: JWT,
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	authHandler := auth.NewHandler(a.DB, a.jwtSecret)
	mux.Handle("/auth/register", authHandler.Register())
	mux.Handle("/auth/login", authHandler.Login())
	//mux.HandleFunc("/health" , a.handleHealth)
	//mux.HandleFunc("/matchmaking", a.handleMatchmaking)
	//mux.HandleFunc("/health/db", a.handleHealthDB)

	return devCORS(requestLogger(mux))
}



func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}