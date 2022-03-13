package gateway

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jponge/playground-go-microservices/pulsar/temperature-store-gateway/data"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func (service Service) Start() {

	service.initPulsarClient()

	address := fmt.Sprintf("%s:%d", service.Host, service.Port)
	srv := &http.Server{
		Addr:         address,
		Handler:      service.router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Println("🚀 Start to listen on", address)
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)

	// Close Pulsar client connection
	service.closePulsarClient()

	log.Println("👋 Bye!")
}

func (service Service) router() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Method(http.MethodGet, "/metrics", promhttp.Handler())
	router.Route("/record", func(r chi.Router) {
		r.Use(TrackIngestionMetrics)
		r.Post("/", service.record)
	})
	return router
}

func (service Service) record(writer http.ResponseWriter, request *http.Request) {
	defer request.Body.Close()
	payload, err := data.PayloadFromReader(request.Body)
	if err != nil {
		log.Println("JSON decoding failed")
		send500(writer, err.Error())
		return
	}
	err = service.pushToPulsar(request.Context(), *payload)
	if err != nil {
		log.Println("Pulsar operation failure")
		send500(writer, err.Error())
		return
	}
	writer.WriteHeader(200)
}

func send500(writer http.ResponseWriter, err string) {
	writer.WriteHeader(500)
	writer.Write([]byte(err))
}
