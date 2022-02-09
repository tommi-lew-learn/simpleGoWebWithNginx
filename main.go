package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/heartbeat", heartbeat)

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
