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

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		lrw := loggingResponseWriter{
			ResponseWriter: w, // compose original http.ResponseWriter
			responseData:   responseData,
		}

		h.ServeHTTP(&lrw, r)

		duration := time.Since(start)

		logrus.WithFields(logrus.Fields{
			"uri":      uri,
			"method":   method,
			"status":   responseData.status,
			"duration": duration,
			"size":     responseData.size,
		}).Info()
	}
	return http.HandlerFunc(logFn)
}

type (
	// struct for holding response details
	responseData struct {
		status int
		size   int
	}

	// our http.ResponseWriter implementation
	loggingResponseWriter struct {
		http.ResponseWriter // compose original http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b) // write response using original http.ResponseWriter
	r.responseData.size += size            // capture size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode) // write status code using original http.ResponseWriter
	r.responseData.status = statusCode       // capture status code
}
