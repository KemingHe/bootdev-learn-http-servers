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

func handleFileServe(fsPref, fsRoot string) http.Handler {
	return http.StripPrefix("/app", http.FileServer(http.Dir(fsRoot)))
}

func main() {
	const fsPref = "/app"
	const fsRoot = "."
	const port = "8080"

	sMux := http.NewServeMux()
	sMux.HandleFunc("/healthz", handleReadiness)
	sMux.Handle(fsPref+"/", handleFileServe(fsPref, fsRoot))

	serv := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port %s\n", fsRoot, port)
	log.Fatal(serv.ListenAndServe()) // Start server
}
