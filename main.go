package main

import (
	"net/http"
)

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("GET /metrics", apiCfg.handleMetrics)
	mux.HandleFunc("GET /reset", apiCfg.handleReset)

	var srv http.Server
	srv.Handler = mux
	srv.Addr = ":8080"
	
	srv.ListenAndServe()
}
