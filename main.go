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
	sMux.HandleFunc("/healthz", handleReadiness)
	sMux.HandleFunc("/metrics", apiCfg.handleMetricsLookup)
	sMux.HandleFunc("/reset", apiCfg.handleMetricsReset)
	sMux.Handle(
		fsPref+"/",
		apiCfg.middelwareMitricsInc(fileServeHandler(fsPref, fsRoot)),
	)

	serv := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port %s\n", fsRoot, port)
	log.Fatal(serv.ListenAndServe()) // Start server
}
