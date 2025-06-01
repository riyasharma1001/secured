package handlers

import (
	"fmt"
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
		// Set security headers
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Get values from environment
		secretKey := os.Getenv("JS_SECRET_KEY")
		apiEndpoint := os.Getenv("API_ENDPOINT")

		// JavaScript code
		originalJS := fmt.Sprintf(`
            (function() {
                function sensitiveFunction() {
                    const secretKey = "%v";
                    const apiEndpoint = "%v";
                    
                    return {
                        doSomething: function() {
                            document.body.innerHTML = "";
                            const div = document.createElement("div");
                            div.style.cssText = "position:fixed;top:50%%;left:50%%;transform:translate(-50%%,-50%%);font-size:24px;color:#00ff00;font-family:Arial;text-align:center;background:#000;padding:20px;border-radius:10px;box-shadow:0 0 10px rgba(0,255,0,0.3);";
                            div.textContent = "Controlled by PhantomCoreX";
                            document.body.style.background = "#000";
                            document.body.appendChild(div);
                        }
                    };
                }
                
                // Auto-execute
                const instance = sensitiveFunction();
                instance.doSomething();
            })();
        `, secretKey, apiEndpoint)

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
