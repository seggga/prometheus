package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// prometheus: declaring metrics
var (
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_response_time_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"hist"},
	)

	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_simplehandler_counter",
			Help: "Number of http requests.",
		},
		[]string{"path"},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// prometheus: declare middleware
func promMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// httpDuration
		timer := prometheus.NewTimer(httpDuration.WithLabelValues("my_metric"))
		defer timer.ObserveDuration()

		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		// requestCounter
		requestCounter.WithLabelValues("/").Inc()

	})
}

// prometheus: register metrics
func init() {
	prometheus.Register(requestCounter)
	prometheus.Register(httpDuration)
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	// business-logic time delay
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	// output
	fmt.Fprint(w, "hello, world!")
}

func main() {
	router := mux.NewRouter()
	router.Use(promMiddleware)

	// Serving root path
	router.Path("/metrics").Handler(promhttp.Handler())
	router.HandleFunc("/", simpleHandler)

	fmt.Println("Serving requests on port 9000")
	err := http.ListenAndServe(":9000", router)
	log.Fatal(err)
}
