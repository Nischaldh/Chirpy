package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Nischaldh/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db *database.Queries
}

func main() {
	const port = "8080"
	godotenv.Load()
	dbURL:= os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
	}
	cfg.fileserverHits.Store(0)
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInt(http.StripPrefix("/app", http.FileServer(http.Dir("")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.numberofRequests)
	mux.HandleFunc("POST /admin/reset", cfg.resetFileServerHits)
	mux.HandleFunc("POST /api/validate_chirp", validateInput)

	ser := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(ser.ListenAndServe())
}
