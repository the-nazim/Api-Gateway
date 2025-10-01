package main

import (
	"io"
	"net/http"
	"net/url"
	"os"

	"gopkg.in/yaml.v3"
)

type Route struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Type    string `yaml:"type"`
	Backend string `yaml:"backend"`
}

type Config struct {
	Routes []Route `yaml:"routes"`
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func proxyHandler(backend string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, err := url.Parse(backend)
		if err != nil {
			http.Error(w, "Invalid backend URL", http.StatusInternalServerError)
			return
		}
		proxyReq, err := http.NewRequest(r.Method, target.String()+r.URL.Path, r.Body)
		if err != nil {
			http.Error(w, "Failed to create request", http.StatusInternalServerError)
			return
		}
		proxyReq.Header = r.Header

		resp, err := http.DefaultClient.Do(proxyReq)
		if err != nil {
			http.Error(w, "Failed to reach backend", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
