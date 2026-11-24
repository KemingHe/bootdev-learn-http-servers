package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const contentPref = "/app"
	const contentRoot = "."
	const adminPref = "/admin"
	const apiPref = "/api"
	const port = "8080"

	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}

	sMux := http.NewServeMux()

	// Content
	sMux.Handle(
		contentPref+"/",
		apiCfg.middlewareMitricsInc(fileServeHandler(contentPref, contentRoot)),
	)

	// Admin
	sMux.HandleFunc("GET "+adminPref+"/metrics", apiCfg.handleMetricsLookup)
	sMux.HandleFunc("POST "+adminPref+"/reset", apiCfg.handleMetricsReset)

	// API
	sMux.HandleFunc("GET "+apiPref+"/healthz", handleReadiness)
	sMux.HandleFunc("POST "+apiPref+"/validate_chirp", handleValidation)

	serv := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port %s", contentRoot, port)
	log.Fatal(serv.ListenAndServe()) // Start server
}
