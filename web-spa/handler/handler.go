package handler

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	/*
		아래의 코드는 /manifest.json, /favicon.ico 와 같은 요청도 index.html로 가게 됨.
		staticDir := http.Dir("build/")
		staticFileHandler := http.FileServer(staticDir)
		// /static/ 이하의 요청은 정적 파일로 서빙 ( "bulid/{url-path}")
		router.PathPrefix("/static/").Handler(http.StripPrefix("/", staticFileHandler))
		// 나머지 모든 요청은 index.html 서빙
		router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "build/index.html") })
	*/

	// http.Handler 인터페이스를 직접 구현
	spa := NewSpaHandler("build", "index.html")
	router.PathPrefix("/").Handler(spa)
	router.Use(CORSMiddleware)
	return router
}
