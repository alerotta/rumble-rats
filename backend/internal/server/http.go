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

	store := auth.NewStore(a.DB)
	svc := auth.NewService(store, a.jwtSecret)
	authHandler := auth.NewHandler(svc)

	//public routes
	mux.Handle("/auth/register", authHandler.Register())
	mux.Handle("/auth/login", authHandler.Login())
	mux.Handle("/auth/validate", authHandler.Validate())

	//protected routes
	
	requireAuth := RequireAuth(svc)
	mux.Handle("/me", requireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"ok":true}`))
	})))
	

	return devCORS(requestLogger(mux))
}



func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}