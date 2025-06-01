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

func getAllowedOrigins() map[string]bool {
	// Create a map of allowed origins
	return map[string]bool{
		"https://webgyans.com":        true,
		"https://health.webgyans.com": true,
		"https://legal.webgyans.com":  true,
		// Add more domains as needed
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	// Pre-compile allowed origins
	allowedOrigins := getAllowedOrigins()

	return func(w http.ResponseWriter, r *http.Request) {
		// Get origin from request
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		if allowedOrigins[origin] {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Requested-With, Origin, Accept")
			w.Header().Set("Access-Control-Max-Age", "3600")
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
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Get encryption key from .env with fallback
	encryptionKey := os.Getenv("JS_SECRET_KEY")
	if encryptionKey == "" {
		encryptionKey = "2xLHEbZAJw6EAoxbPXlrdYleZJBOsXmg" // Default key
	}

	// Initialize security pipeline
	pipeline := security.NewPipeline(&security.Config{
		EncryptionKey:   encryptionKey,
		EnableAntiDebug: true,
		EnableWASM:      true,
	})

	// Create rate limiter
	rateLimiter := middleware.NewRateLimiter(rate.Limit(10/60.0), 20)

	// Setup server routes
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/js/protected.js",
		rateLimiter.Limit(
			enableCORS(
				handlers.ServeProtectedJS(pipeline),
			),
		),
	)

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
