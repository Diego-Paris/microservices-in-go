package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Diego-Paris/microservices-in-go/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	//? We can inject our logger, and change this dependency from here
	l := log.New(os.Stdout, "product-api ", log.LstdFlags)

	ph := handlers.NewProducts(l) // products handler

	sm := mux.NewRouter()

	// [GET]
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/{id:[0-9]+}", ph.GetProduct)
	getRouter.HandleFunc("/", ph.GetProducts)

	// [PUT]
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProducts) // Updates the whole product
	putRouter.Use(ph.MiddlewareProductValidation)

	// [POST]
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.AddProduct)
	postRouter.Use(ph.MiddlewareProductValidation)

	// [DELETE]
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/{id:[0-9]+}", ph.DeleteProduct)

	// [DOCS]
	ops := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(ops, nil)
	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// [SERVER CONFIG]
	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	// [RUNS SERVER ON GOROUTINE]
	go func() {
		l.Printf("Starting server on port %s\n", s.Addr)
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	//? this context allows requests 30 seconds
	//? to attempt gracefully shutting down
	//? if handlers are still working after
	//? 30 seconds forcefully shutdown

	// [GRACEFUL SHUTDOWN ON KILL OR INTERRUPT]
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("[EXIT] Received terminate, graceful shutdown", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc) // waits until the requests are completed, then shutdown
	cancel()
}
