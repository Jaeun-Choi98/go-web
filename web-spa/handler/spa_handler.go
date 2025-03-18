package handler

import (
	"net/http"
	"os"
	"path/filepath"
)

// http.Handler 인터페이스 구현 구조체
type spaHandler struct {
	staticPath string
	indexPath  string
}

func NewSpaHandler(sp, ip string) spaHandler {
	return spaHandler{staticPath: sp, indexPath: ip}
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join(h.staticPath, r.URL.Path)
	file, err := os.Stat(path)
	if err == nil && !file.IsDir() {
		http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
		return
	}
	// fallback
	http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
}
