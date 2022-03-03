package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jponge/playground-go-microservices/db-temperature-store/controller"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Start(host string, port int, controller controller.Controller) {
	router := AppRouter(controller)
	address := fmt.Sprintf("%s:%d", host, port)
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Println("🚀 Start to listen on", address)
		err := server.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("👋 Bye!")
}

func AppRouter(controller controller.Controller) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/record", controller.Record)
	router.Get("/data/{id}", controller.FetchOne)
	router.Get("/data", controller.FetchAll)
	return router
}
