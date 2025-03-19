package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var GoogleOauthConfig oauth2.Config

func NewHandler() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/auth/google/login", googleLoginHandler)
	mux.HandleFunc("/auth/google/callback", googleOatuhCallback)
	return mux
}

// 2. 클라이언트가 정상적으로 로그인을 하면, Athorization code를 받음
func googleOatuhCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")
	if r.FormValue("state") != oauthState.Value {
		log.Printf("invaild google oauth state cookie:%s state:%s \n", oauthState.Value, r.FormValue("state"))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	data, _ := getGoogleUserInfo(r.FormValue("code"))
	fmt.Fprint(w, string(data))
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	// 3. exchange code for toekn
	token, err := GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo: %s\n", err.Error())
	}

	// 4. token을 사용해서 api 요청
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to request google api: %s\n", err.Error())
	}

	return io.ReadAll(resp.Body)
}

// 1. token 요청
func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// verify using cookie, token(jwt)을 사용한 방식도 존재.
	state := generateStateOauthCookie(w)
	// Redirect google login url(e.g. https://accounts.google.com/o/oauth2/auth?client_id=518***&redirect_uri=http%3A%2F%2...)
	url := GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(24 * time.Hour)
	buf := make([]byte, 16)
	rand.Read(buf)
	state := base64.URLEncoding.EncodeToString(buf)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	http.ServeFile(w, r, filepath.Join("build", "index.html"))
}
