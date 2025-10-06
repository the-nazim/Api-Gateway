package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	// "honnef.co/go/tools/lintcmd/cache"
)

type responseWriter struct {
	http.ResponseWriter
	body *[]byte
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("[%s] %s %s %s", r.Method, r.RequestURI, r.RemoteAddr, duration)
	})
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			w.Header().Set("X-Request-ID", requestID)
		}

		w.Header().Set("X-Request-ID", requestID)
		log.Printf("RequestID=%s %s %s", requestID, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("Error parsing IP:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		key := "rate_limit:" + ip
		pipe := rdb.TxPipeline()
		incrCmd := pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
		_, err = pipe.Exec(ctx)
		if err != nil {
			log.Println("Redis error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		count := incrCmd.Val()

		// if count == 1 {
		// 	rdb.Expire(ctx, key, time.Minute)
		// }

		if count > 10 {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("Rate limit exceeded"))
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		cacheKey := "cache:" + r.URL.String()
		cached, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(cached))
			return
		}

		rw := &responseWriter{ResponseWriter: w, body: &[]byte{}}
		next.ServeHTTP(rw, r)

		rdb.Set(ctx, cacheKey, string(*rw.body), 30*time.Minute)
	})
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	*rw.body = append(*rw.body, b...)
	return rw.ResponseWriter.Write(b)
}
