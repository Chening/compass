package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// 1.接收客户端 request，并将 request 中带的 header 写入 response header
// reference
// https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go
// start shell
// go run  task1.go
func main() {
	// back headers
	http.HandleFunc("/back_header", backHeader)

	err := http.ListenAndServe(":19001", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}

func backHeader(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		// unknow ??
		// 放开此处打印 response中的header则 并非for循环中设置的数值 原因未知
		// fmt.Fprintf(w, "Header field %q, Value %q\n", k, v)

		w.Header().Set(k, strings.Join(v, ", "))

		// local print
		fmt.Printf("Header field %q, Value %q\n", k, v)
	}

	io.WriteString(w, "This is my back headers page!\n")
}
