// filepath: internal/handlers/api.go
package handlers

import (
	"fmt"
	"net/http"
	"os"
	"secureserver/internal/security"
)

func NewAPIHandler(pipeline *security.Pipeline) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiEndpoint := os.Getenv("API_ENDPOINT")

		// Validate API endpoint
		if apiEndpoint == "" {
			http.Error(w, "API endpoint not configured", http.StatusInternalServerError)
			return
		}

		// Test API connection
		resp, err := http.Get(apiEndpoint)
		if err != nil {
			http.Error(w, fmt.Sprintf("API connection failed: %v", err), http.StatusServiceUnavailable)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "API returned non-200 status", resp.StatusCode)
			return
		}

		// API is working
		w.WriteHeader(http.StatusOK)
	}
}

func NewProtectedJSHandler(pipeline *security.Pipeline) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Original JS code
		jsCode := `function secretFunction() { return "sensitive data"; }`

		// Process through security pipeline
		processed, err := pipeline.Process([]byte(jsCode))
		if err != nil {
			http.Error(w, "Processing error", http.StatusInternalServerError)
			return
		}

		// Set headers
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "no-store")
		w.Write(processed)
	}
}
