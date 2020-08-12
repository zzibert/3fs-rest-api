package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "3fs-rest-api", log.LstdFlags)

	// create the user handlers
	userHandler := handlers.NewUsers()

	// create the group handlers
	groupHandler := handlers.NewGroups()

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API

	// GET Subrouter
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/users", userHandler.ListAll)
	getRouter.HandleFunc("/users/{id:[0-9]+}", userHandler.ListSingle)
	getRouter.HandleFunc("/groups", groupHandler.ListAll)
	getRouter.HandleFunc("/groups/{id:[0-9]+}", groupHandler.ListSingle)

	// PUT Subrouter
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/users", userHandler.Update)
	putRouter.HandleFunc("/groups", groupHandler.Update)

	// POST Subrouter
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/users", userHandler.Create)
	postRouter.HandleFunc("/groups", groupHandler.Create)

	// DELETE Subrouter
	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/users/{id:[0-9]+}", userHandler.Delete)
	deleteRouter.HandleFunc("/groups/{id:[0-9]+}", groupHandler.Delete)

	// create a new server
	s := http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start a new server
	go func() {
		l.Println("Starting the server on port 8080")

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received
	sig := <-c
	log.Println("Got Signal ", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
