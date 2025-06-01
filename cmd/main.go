package main

import (
	"log"
	"net/http"
	"secureserver/internal/handlers"
	"secureserver/internal/middleware"
	"secureserver/internal/security"

	"golang.org/x/time/rate"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow only webgyans.com
		w.Header().Set("Access-Control-Allow-Origin", "https://webgyans.com")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin")

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Initialize security pipeline
	pipeline := security.NewPipeline(&security.Config{
		EncryptionKey:   "your-secure-key-here",
		EnableAntiDebug: true,
		EnableWASM:      true,
	})

	// Create rate limiter: 10 requests per minute, burst of 20
	rateLimiter := middleware.NewRateLimiter(rate.Limit(10/60.0), 20)

	// Apply rate limiting and CORS
	http.HandleFunc("/js/protected.js",
		rateLimiter.Limit(
			enableCORS(
				handlers.ServeProtectedJS(pipeline),
			),
		),
	)

	// Add API endpoint test route
	http.HandleFunc("/api/test",
		rateLimiter.Limit(
			enableCORS(
				handlers.NewAPIHandler(pipeline),
			),
		),
	)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
