package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type validReq struct {
	Body string `json:"body"`
}

type errorRes struct {
	Error string `json:"error"`
}

type validRes struct {
	Valid bool `json:"valid"`
}

func writeJSONRes(
	w http.ResponseWriter,
	code int,
	payload any,
) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		w.Write(dat)
	}
}

func handleValidation(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	vr := validReq{}
	if err := decoder.Decode(&vr); err != nil {
		log.Printf("Error decoding request: %s", err)
		writeJSONRes(w, http.StatusBadRequest, errorRes{Error: "Something went wrong"})
		return
	}

	const maxBodyLen = 140
	bodyLen := len(vr.Body)
	if bodyLen > maxBodyLen {
		log.Printf(
			"Received invalid chirp body of len %d, exceeding maxBodyLen of %d",
			bodyLen,
			maxBodyLen,
		)
		writeJSONRes(w, http.StatusBadRequest, errorRes{Error: "Chirp is too long"})
		return
	}

	// Onwards, bodyLen <= maxBodyLen
	log.Printf(
		"Received valid chirp body of len %d, within maxBodyLen of %d",
		bodyLen,
		maxBodyLen,
	)
	writeJSONRes(w, http.StatusOK, validRes{Valid: true})
}
