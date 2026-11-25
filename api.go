package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

type cleanedRes struct {
	CleanedBody string `json:"cleaned_body"`
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

func cleanAndCount(body string) (string, int) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	const mask = "****"
	const sep = " "
	words := strings.Split(body, sep)
	cleaned := []string{}
	count := 0
	for _, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			cleaned = append(cleaned, mask)
			count += 1
		} else {
			cleaned = append(cleaned, word)
		}
	}
	return strings.Join(cleaned, sep), count
}

func handleValidation(w http.ResponseWriter, req *http.Request) {
	// Step 1. Validate request JSON shape
	decoder := json.NewDecoder(req.Body)
	vr := validReq{}
	if err := decoder.Decode(&vr); err != nil {
		log.Printf("Error decoding request: %s", err)
		writeJSONRes(w, http.StatusBadRequest, errorRes{Error: "Something went wrong"})
		return
	}

	// Step 2. Check request body within max length
	const maxLen = 140
	bodyLen := len(vr.Body)
	if bodyLen > maxLen {
		log.Printf(
			"Received invalid bodyLen %d > maxLen %d",
			bodyLen,
			maxLen,
		)
		writeJSONRes(w, http.StatusBadRequest, errorRes{Error: "Chirp is too long"})
		return
	}
	// Onwards, bodyLen <= maxLen

	// Step 3. Detect and sensor bad words
	cleanedBody, count := cleanAndCount(vr.Body)
	log.Printf(
		"Received valid bodyLen %d <= maxLen %d; sensored %d bad word(s)",
		bodyLen,
		maxLen,
		count,
	)
	writeJSONRes(w, http.StatusOK, cleanedRes{CleanedBody: cleanedBody})
}
