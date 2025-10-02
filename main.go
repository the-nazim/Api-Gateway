package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalln("Error loading config:", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, route := range config.Routes {
			if strings.HasPrefix(r.URL.Path, route.Path) {
				// target := route.Backend + r.URL.Path
				target, err := url.Parse(route.Backend)
				if err != nil {
					http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
					return
				}
				log.Printf("Proxying request: %s -> %s", r.URL.Path, target)

				// proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: strings.TrimPrefix(route.Backend, "http://")})
				proxy := httputil.NewSingleHostReverseProxy(target)
				// r.URL.Path = strings.TrimPrefix(r.URL.Path, route.Path) // optional clean-up
				proxy.ServeHTTP(w, r)
				return
			}
		}
		http.NotFound(w, r)
	})

	handler := requestIDMiddleware(loggingMiddleware(mux))

	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
