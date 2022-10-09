package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// 2.读取当前系统的环境变量中的 VERSION 配置，并写入 response header
// start shell 
// go run  task2.go
// go run  task2.go  -version=V1
// access url -> http://127.0.0.1:19002/back_header_with_version
func main() {
	version := flag.String("VERSION", "DEFAULT", "请输入VERSION 配置")

	//import dont forget invole this function
	flag.Parse()

	// back headers with version
	http.HandleFunc("/back_header_with_version", func(w http.ResponseWriter, r *http.Request) {

		for k, v := range r.Header {
			w.Header().Set(k, strings.Join(v, ", "))

			// local print
			fmt.Printf("Header field %q, Value %q\n", k, v)
		}

		w.Header().Set("VERSION", *version)

		io.WriteString(w, fmt.Sprintf("This is my back headers page, version %v!\n", *version))
	})

	err := http.ListenAndServe(":19002", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
