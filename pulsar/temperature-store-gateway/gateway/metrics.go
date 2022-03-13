package gateway

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var (
	ingestedUpdates = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gateway_ingested_updates",
		Help: "Number of ingested temperature updates",
	})
)

func init() {
	prometheus.MustRegister(ingestedUpdates)
}

func TrackIngestionMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(writer, request)
		ingestedUpdates.Inc()
	})
}
