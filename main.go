package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/t-morgan/chirpy/internal/database"
)

func main() {
	godotenv.Load()
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: 0,
		dbQueries:      dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handleHealthz)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("PUT /api/users", apiCfg.handleUpdateUser)

	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlePolkaWebhook)

	mux.HandleFunc("POST /api/chirps", apiCfg.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleDeleteChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	var srv http.Server
	srv.Handler = mux
	srv.Addr = ":8080"

	log.Fatal(srv.ListenAndServe())
}
