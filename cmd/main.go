package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Yangiboev/request-counter/config"
	"github.com/Yangiboev/request-counter/server"
)

func main() {
	server := server.NewServer(config.LoadConfig())
	server.Routes()

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)

	go func() {
		signal := <-quit
		log.Printf("Server received signal '%v'.", signal)
		if err := server.PersistState(); err != nil {
			log.Fatalf("Could not store the state: %v\n", err)
		}

		log.Println("Shutting down...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Error on gracefully shutdowning the server: %v\n", err)
		}
		close(done)
	}()

	log.Println("Server listens to the port:", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}

	<-done
}
