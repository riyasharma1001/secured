package main

import (
	"log"
	"net/http"
	"os"
	"secureserver/internal/handlers"
	"secureserver/internal/middleware"
	"secureserver/internal/security"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set single CORS header
		origin := r.Header.Get("Origin")
		if origin == "https://webgyans.com" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin")
		}

		// Handle preflight
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Load .env file from server root
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
		// Don't fatal here, allow default values
	}

	// Get encryption key from .env with fallback
	encryptionKey := os.Getenv("JS_SECRET_KEY")
	if encryptionKey == "" {
		encryptionKey = "2xLHEbZAJw6EAoxbPXlrdYleZJBOsXmg" // Default key
	}

	// Initialize security pipeline with env key
	pipeline := security.NewPipeline(&security.Config{
		EncryptionKey:   encryptionKey,
		EnableAntiDebug: true,
		EnableWASM:      true,
	})

	// Create rate limiter: 10 requests per minute, burst of 20
	rateLimiter := middleware.NewRateLimiter(rate.Limit(10/60.0), 20)

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	// Apply rate limiting and CORS
	http.HandleFunc("/js/protected.js",
		rateLimiter.Limit(
			enableCORS(
				handlers.ServeProtectedJS(pipeline),
			),
		),
	)

	log.Println("Starting local development server on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
