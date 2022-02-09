package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/heartbeat", heartbeat)

	http.ListenAndServe("localhost:8000", mux)
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Pong")
}
