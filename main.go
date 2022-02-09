package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api/heartbeat", WithLogging(heartbeatHandler()))

	addr := "localhost:8000"
	logrus.WithField("addr", addr).Info("starting server")

	err := http.ListenAndServe("localhost:8000", mux)

	if err != nil {
		logrus.WithField("event", "start server").Fatal(err)
	}
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pong")
}

func heartbeatHandler() http.Handler {
	fn := heartbeat
	return http.HandlerFunc(fn)
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method

		h.ServeHTTP(w, r)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"uri":      uri,
			"method":   method,
			"duration": duration,
		}).Info()
	}
	return http.HandlerFunc(logFn)
}
