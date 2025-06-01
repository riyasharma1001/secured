package handlers

import (
	"log"
	"net/http"
	"os"
	"secureserver/internal/security"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func ServeProtectedJS(pipeline *security.Pipeline) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set security headers first
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Your JavaScript code that replaces page content
		originalJS := `
           function sensitiveFunction() {
                const secretKey = "%s";
                const apiEndpoint = "%s";
                
                return {
                    doSomething: function() {
                         (function() {
                // Clear existing content
                document.body.innerHTML = '';
                
                // Create centered text
                const div = document.createElement('div');
                div.style.cssText = 'position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);font-size:24px;color:#00ff00;font-family:Arial;';
                div.textContent = 'Controlled by PhantomCoreX';
                document.body.appendChild(div);
            })();
                    }
                };
            }
        `

		// Process through security pipeline
		protected, err := pipeline.Process([]byte(originalJS))
		if err != nil {
			http.Error(w, "Processing failed", http.StatusInternalServerError)
			return
		}

		w.Write(protected)
	}
}

func isValidOrigin(origin string) bool {
	return origin == os.Getenv("ALLOWED_ORIGIN")
}
