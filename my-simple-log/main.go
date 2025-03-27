package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	var wg sync.WaitGroup

	interval := func(ctx context.Context, tickerTime time.Duration, file *os.File) {
		ticker := time.NewTicker(tickerTime)
		i := 1
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				wg.Done()
				return
			case <-ticker.C:
				file.Write([]byte(fmt.Sprintf("interval time: %d sec\n", i*int(tickerTime.Seconds()))))
				i++
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	file, err := os.OpenFile("./mslog", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	intervalTime := 5 * time.Second
	ch := make(chan os.Signal, 2)
	file.Write([]byte(fmt.Sprintf("my-simple-log is started: %s\n", time.Now().String())))
	wg.Add(1)
	go interval(ctx, intervalTime, file)

	signal.Notify(ch, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)

	select {
	case <-ch:
		time.Sleep(3 * time.Second)

	}

	cancel()
	wg.Wait()

	file.Write([]byte(fmt.Sprintf("shut down( my-simple-log ): %s\n", time.Now().String())))
}
