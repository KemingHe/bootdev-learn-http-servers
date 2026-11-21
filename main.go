package main

import (
	"log"
	"net/http"
)

type apiHandler struct{}

func main() {
	const sPrt = "8080"
	sMux := http.NewServeMux()
	serv := &http.Server{
		Addr:    ":" + sPrt,
		Handler: sMux,
	}

	log.Printf("Serving on port %s\n", sPrt)
	log.Fatal(serv.ListenAndServe()) // Start server
}
