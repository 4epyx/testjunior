package http

import "net/http"

type CustomResponseWriter struct {
	http.ResponseWriter

	status int
	msg    []byte
}

func WrapResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		ResponseWriter: w,
	}
}

func (r *CustomResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *CustomResponseWriter) Write(msg []byte) (int, error) {
	if r.status > 299 {
		r.msg = append(r.msg, msg...)
	}

	return r.ResponseWriter.Write(msg)
}

func (r *CustomResponseWriter) Status() int {
	return r.status
}

func (r *CustomResponseWriter) StringMsg() string {
	return string(r.msg)
}
