package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
    "math/rand"
    "time"

	"github.com/sirupsen/logrus"

    "github.com/gorilla/mux"

    "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)


// https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
var requestDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
    Name:    "http_request_duration_seconds",
    Help:    "A histogram of the HTTP request durations in seconds.",
    Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
})

func init() {
	prometheus.Register(requestDurations)
}



func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        now := time.Now()

		next.ServeHTTP(w, r)

        requestDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
            time.Since(now).Seconds(), prometheus.Labels{"api": "healthz"},
        )
	})
}


func main() {
    // Create non-global registry.
	registry := prometheus.NewRegistry()

    // Add go runtime metrics and process collectors.
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		requestDurations,
	)

	router := mux.NewRouter()
    router.Use(prometheusMiddleware)

    // Prometheus endpoint
	router.Path("/metrics").Handler(promhttp.Handler())

	// back headers
    router.Path("/healthz").HandlerFunc(healthz)
	http.HandleFunc("/healthz", healthz)


	go func(iRouter *mux.Router) {
		err := http.ListenAndServe(":19004", iRouter)

		if errors.Is(err, http.ErrServerClosed) {
			
		} else if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			  }).Warn("error starting server")
			os.Exit(1)
		}
	}(router)
	

	quit := make(chan os.Signal ,1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("我被优雅终止了")

}

func healthz(w http.ResponseWriter, r *http.Request) {
    now := time.Now()

	logrus.Info("info accessed healthz")
    logrus.Error("error accessed healthz")

    //为 HTTPServer 添加 0-2 秒的随机延时
    rand.Seed(now.UnixNano())
    randTime := rand.Intn(2000)
    fmt.Printf("content length %d, sleep %d ms \n ", r.ContentLength, randTime)

    time.Sleep(time.Duration(randTime) * time.Millisecond)

    //
    // requestDurations.(prometheus.ExemplarObserver).ObserveWithExemplar(
    //     time.Since(now).Seconds(), prometheus.Labels{"api": "healthz"},
    // )

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "This is my status eq 200 page!\n")
}

