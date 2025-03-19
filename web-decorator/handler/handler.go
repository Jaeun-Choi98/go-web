package handler

import (
	"fmt"
	"log"
	"net/http"
	"root/decorator"
	"time"

	"github.com/gorilla/mux"
)

func NewHandler() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/", indexHandler)
	loggerHandler := decorator.NewDecoratorHandler(logger, mux)
	//loggerHadnler2 := decorator.NewDecoratorHandler(logger2, loggerHandler)
	return loggerHandler
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "hellow world!")
}

func logger(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER] Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER] Completed time: ", time.Since(start).Milliseconds())
}

func logger2(w http.ResponseWriter, r *http.Request, h http.Handler) {
	start := time.Now()
	log.Println("[LOGGER2] Started")
	h.ServeHTTP(w, r)
	log.Println("[LOGGER2] Completed time: ", time.Since(start).Milliseconds())
}
