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
		// Set security headers
		w.Header().Set("Content-Type", "application/javascript")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// JavaScript code that runs immediately
		originalJS := `
            console.log("Executing script...");
            document.body.innerHTML = "";
            document.body.style.margin = "0";
            document.body.style.background = "#000";
            const div = document.createElement("div");
            div.style.cssText = "position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);font-size:24px;color:#00ff00;font-family:Arial;text-align:center;background:#000;padding:20px;border-radius:10px;box-shadow:0 0 10px rgba(0,255,0,0.3);z-index:999999;";
            div.textContent = "Controlled by PhantomCoreX";
            document.body.appendChild(div);
        `

		// Temporarily skip encryption for testing
		w.Write([]byte(originalJS))
	}
}

func isValidOrigin(origin string) bool {
	return origin == os.Getenv("ALLOWED_ORIGIN")
}
