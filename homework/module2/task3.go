package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// 3.Server 端记录访问日志包括客户端 IP，HTTP 返回码，输出到 server 端的标准输出
// start shell
// go run  task3.go
// access url -> http://127.0.0.1:19003/trace
// reference
// https://www.alexedwards.net/blog/making-and-using-middleware
// https://dev.to/julienp/logging-the-status-code-of-a-http-handler-in-go-25aa

func main() {
	mux := http.NewServeMux()

	traceHandler := http.HandlerFunc(trace)

	mux.Handle("/trace", loggerMiddleware(traceHandler))

	err := http.ListenAndServe(":19003", mux)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(recorder, r)
		log.Printf("Handling request for %s from %s, status: %d", r.URL.Path, r.RemoteAddr, recorder.Status)
	})
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func trace(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "This is my trace page!\n")
}
