package main

import (
	"log"
	"net/http"
)

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	const fileRoot = "."
	const port = "8080"

	sMux := http.NewServeMux()
	sMux.HandleFunc("/healthz", handleReadiness)
	sMux.Handle("/app/", http.StripPrefix( // Serve root files from /app/ path
		"/app",
		http.FileServer(http.Dir(fileRoot)),
	))

	serv := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port %s\n", fileRoot, port)
	log.Fatal(serv.ListenAndServe()) // Start server
}
