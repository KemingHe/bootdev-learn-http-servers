package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const fsPref = "/app"
	const fsRoot = "."
	const port = "8080"

	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}

	sMux := http.NewServeMux()
	sMux.HandleFunc("GET /healthz", handleReadiness)
	sMux.HandleFunc("GET /metrics", apiCfg.handleMetricsLookup)
	sMux.HandleFunc("POST /reset", apiCfg.handleMetricsReset)
	sMux.Handle(
		fsPref+"/",
		apiCfg.middelwareMitricsInc(fileServeHandler(fsPref, fsRoot)),
	)

	serv := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port %s", fsRoot, port)
	log.Fatal(serv.ListenAndServe()) // Start server
}
