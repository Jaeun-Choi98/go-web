package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pjt/juchoi/handler"
	"sync"
	"time"
)

func main() {
	router := handler.NewRouter()
	server := &http.Server{
		Handler: router,
		//Addr:    "127.0.0.1:8080",
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	c := make(chan os.Signal, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			wg.Done()
		}
	}()

	signal.Notify(c, os.Interrupt)
	<-c

	server.Shutdown(ctx)
	wg.Wait()
	log.Println("shutting down")
	os.Exit(0)
}
