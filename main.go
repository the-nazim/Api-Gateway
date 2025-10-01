package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalln("Error loading config:", err)
		return
	}

	mux := http.NewServeMux()
	for _, route := range config.Routes {
		fmt.Printf("Requesting route: %s->%s\n", route.Path, route.Backend)
		mux.HandleFunc(route.Path, proxyHandler(route.Backend))
	}

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
