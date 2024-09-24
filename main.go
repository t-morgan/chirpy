package main

import (
	"log"
	"net/http"
)

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /api/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	var srv http.Server
	srv.Handler = mux
	srv.Addr = ":8080"

	log.Fatal(srv.ListenAndServe())
}
