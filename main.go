package main

import (
	"log"
	"net/http"
)

type apiHandler struct{}

func main() {
	const sPath = "."
	const sPrt = "8080"

	sMux := http.NewServeMux()
	sMux.Handle("/", http.FileServer(http.Dir(sPath)))

	serv := &http.Server{
		Addr:    ":" + sPrt,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port %s\n", sPath, sPrt)
	log.Fatal(serv.ListenAndServe()) // Start server
}
