package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type LogRW struct {
	http.ResponseWriter
	status int
	size   int
}

func newLogRW(w http.ResponseWriter) *LogRW {
	return &LogRW{ResponseWriter: w}
}

func (rw *LogRW) Status() int {
	return rw.status
}
func (rw *LogRW) Size() int {
	return rw.size
}

// Implements http.ResponseWriter
func (rw *LogRW) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *LogRW) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func (rw *LogRW) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Implements http.Hijacker
func (rw *LogRW) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("http.ResponseWriter does not support Hijacker")
}
