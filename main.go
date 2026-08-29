package main

import (
	
	"log"
	"sync/atomic"

	"net/http"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}




func main() {
	const port = "8080"
	cfg := &apiConfig{}
	cfg.fileserverHits.Store(0)
	mux := http.NewServeMux()
	mux.Handle("/app/", cfg.middlewareMetricsInt(http.StripPrefix("/app", http.FileServer(http.Dir("")))))
	mux.HandleFunc("GET /api/healthz",handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.numberofRequests)
	mux.HandleFunc("POST /admin/reset", cfg.resetFileServerHits)
	mux.HandleFunc("POST /api/validate_chirp", validateInput)

	ser := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Fatal(ser.ListenAndServe())
}
