package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) numberofRequests(w http.ResponseWriter, r *http.Request) {
	numberofHits := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html")
	html := fmt.Sprintf(`
	<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, numberofHits)
	fmt.Fprintf(w, "%s", html)
}
