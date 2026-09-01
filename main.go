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
	platform string
	jwtSecret string
	polkaKey string
}

func main() {
	const port = "8080"
	godotenv.Load()
	dbURL:= os.Getenv("DB_URL")
	polkaApiKey := os.Getenv("POLKA_KEY")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	cfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db: dbQueries,
		platform: platform,
		jwtSecret: jwtSecret,
		polkaKey: polkaApiKey,
	}
	cfg.fileserverHits.Store(0)
	mux := http.NewServeMux()

	// app endpoints
	mux.Handle("/app/", cfg.middlewareMetricsInt(http.StripPrefix("/app", http.FileServer(http.Dir("")))))



	// api endpoints

	//admin endpoints
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.numberofRequests)
	mux.HandleFunc("POST /admin/reset", cfg.resetFileServerHits)

	//users endpont
	mux.HandleFunc("POST /api/users", cfg.createUsers)
	mux.HandleFunc("PUT /api/users", cfg.updateUsers)
	mux.HandleFunc("POST /api/login", cfg.loginUsers)

	//chirps endpoint
	mux.HandleFunc("POST /api/chirps", cfg.createChirps)
	mux.HandleFunc("GET /api/chirps", cfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)

	//refresh token endpoint
	mux.HandleFunc("POST /api/refresh", cfg.refresh)
	mux.HandleFunc("POST /api/revoke", cfg.revoke)

	//webhook
	mux.HandleFunc("POST /api/polka/webhooks", cfg.updateUserStatus)

	ser := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(ser.ListenAndServe())
}
