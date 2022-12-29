package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

type MyResponseWriter struct {
	StatusCode int
	Body       []byte
	http.ResponseWriter
}

func NewResponseWriter(w http.ResponseWriter) *MyResponseWriter {
	return &MyResponseWriter{ResponseWriter: w}
}

func (mw *MyResponseWriter) Header() http.Header {
	return mw.ResponseWriter.Header()
}

func (mw *MyResponseWriter) Write(b []byte) (int, error) {
	mw.Body = b
	return mw.ResponseWriter.Write(b)
}

func (mw *MyResponseWriter) WriteHeader(statusCode int) {
	mw.StatusCode = statusCode
	mw.ResponseWriter.WriteHeader(statusCode)
}

func logRequestHandler(h http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var mw *MyResponseWriter = NewResponseWriter(w)

		// 简单获取ip
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = ""
		}
		log.Printf("get request from %s\n", ip)

		h.ServeHTTP(mw, r)

		if mw.StatusCode == 0 {
			mw.StatusCode = http.StatusOK
		}

		log.Printf("response status code=%d\tbody=%s", mw.StatusCode, string(mw.Body))
	}

	return fn
}

func index(w http.ResponseWriter, r *http.Request) {
	ver := os.Getenv("VERSION")
	fmt.Printf("set VERSION: %s\n", ver)
	w.Header().Set("VERSION", ver)

	for key, val := range r.Header {
		fmt.Printf("set res header -- key: %s value: %s\n", key, strings.Join(val, ", "))
		w.Header().Set(key, strings.Join(val, ", "))
	}
	io.WriteString(w, "ok")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "200")
}

func main() {
	const port = 80
	http.HandleFunc("/", logRequestHandler(index))
	http.HandleFunc("/healthz", logRequestHandler(healthz))

	log.Printf("Starting HTTP server at port: %d\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	if err != nil {
		log.Fatal(err)
	}
}
