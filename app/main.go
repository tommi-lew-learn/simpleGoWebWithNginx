package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/api/heartbeat", WithLogging(heartbeatHandler()))
	mux.Handle("/api/time", WithLogging(timesHandler()))

	addr := "0.0.0.0:8000"
	logrus.WithField("addr", addr).Info("starting server")

	err := http.ListenAndServe(addr, mux)

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

func timesHandler() http.Handler {
	fn := times
	return http.HandlerFunc(fn)
}

func times(w http.ResponseWriter, r *http.Request) {
	tz := r.FormValue("tz")
	multipleTz := tz
	timezones := strings.Split(multipleTz, ",")

	// If the tz parameter is not provided in the URL
	if len(timezones) == 1 && timezones[0] == "" {
		// replace empty string with "UTC
		timezones[0] = "UTC"
	}

	var localTimes map[string]string = currentTimes(&timezones)

	w.Header().Add("Content-Type", "application/json")

	// TODO: Handle potential encoding issue
	json.NewEncoder(w).Encode(localTimes)
}

func currentTimes(zones *[]string) map[string]string {
	localTimesMap := make(map[string]string)

	for _, zone := range *zones {
		loc, err := time.LoadLocation(zone)

		if err != nil {
			localTimesMap[zone] = "invalid timezone"
		} else {
			localTimesMap[zone] = fmt.Sprintf("%s", time.Now().In(loc))
		}
	}

	return localTimesMap
}

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

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
			"uri":      r.RequestURI,
			"method":   r.Method,
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
