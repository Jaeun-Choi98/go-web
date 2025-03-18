package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"turker.web/model"
)

func FileHandler(w http.ResponseWriter, r *http.Request) {

	uploadFile, header, err := r.FormFile("upload_file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		fmt.Fprint(w, err)
		return
	}
	defer uploadFile.Close()

	dirName := "./resource"
	os.MkdirAll(dirName, 0777)
	filePath := fmt.Sprintf("%s/%s", dirName, header.Filename)
	file, err := os.Create(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}
	defer file.Close()
	io.Copy(file, uploadFile)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, filePath)
}

func JsonHandler(w http.ResponseWriter, r *http.Request) {
	user := model.NewUser()
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user.CreatedAt = time.Now()
	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(user)
}

func BarHandler(w http.ResponseWriter, r *http.Request) {
	/*
		// handleFunc("/bar/{id}", barHandler)
		r.PathValue("id")
	*/
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	w.Write([]byte("hello " + name))
}

type FooHandler struct{}

func (fh FooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello foo"))
}
