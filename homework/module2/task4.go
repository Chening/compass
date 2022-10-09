package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// 4.当访问 localhost/healthz 时，应返回 200
// start shell
// go run  task4.go
// reference
// https://go.dev/src/net/http/status.go
func main() {
	// back headers
	http.HandleFunc("/healthz", healthz)

	err := http.ListenAndServe(":19004", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "This is my status eq 200 page!\n")
}
