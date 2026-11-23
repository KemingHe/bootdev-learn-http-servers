package main

import "net/http"

func fileServeHandler(fsPref, fsRoot string) http.Handler {
	return http.StripPrefix("/app", http.FileServer(http.Dir(fsRoot)))
}
