package main

import (
	"net/http"
	"os"
	"root/handler"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	godotenv.Load(".env")
	handler.GoogleOauthConfig = oauth2.Config{
		RedirectURL:  "http://localhost:3000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	http.ListenAndServe(":3000", handler.NewHandler())
}
