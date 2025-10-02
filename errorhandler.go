package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func respondError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func doRequestWithRetry(req *http.Request, retries int) (*http.Response, error) {
	var resp *http.Response
	var err error
	for i := 0; i <= retries; i++ {
		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			return resp, nil
		}
		log.Printf("Retry %d after error: %v", i+1, err)
		time.Sleep(100 * time.Millisecond)
	}
	return nil, err
}
