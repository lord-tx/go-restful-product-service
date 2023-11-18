package main

import (
	"context"
	"go-restful-product-service/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	ph := handlers.NewProduct(l)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetProducts)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)

	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 8080")

		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// Create a context with a timeout of 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to gracefully shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Error shutting down server: %s\n", err)
	}

}
