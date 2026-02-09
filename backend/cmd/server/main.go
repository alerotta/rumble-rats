package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/alerotta/rumble-rats/backend/internal/db"
	"github.com/alerotta/rumble-rats/backend/internal/server"
)

func main (){

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
	log.Fatal("DATABASE_URL is not set")
	}

	s := os.Getenv("JWT_SECRET")
	if s == ""{
	log.Fatal("JWT_SECRET is not set")
}
	dbConn, err := db.Open(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	app := server.NewApp(dbConn, s)
	srv := &http.Server{
		Addr: ":8080",
		Handler: app.Routes(),
	}

	go func (){
		log.Printf("HTTP server listening on http://localhost%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make (chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Printf("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	log.Printf("bye")
}
