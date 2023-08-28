package http

import "net/http"

// CustomResponseWriter is http response writer, that can contain status of the response and body if request is not successful
type CustomResponseWriter struct {
	http.ResponseWriter

	status int
	msg    []byte
}

// WrapResponseWriter may be used for wrapping default response writer
func WrapResponseWriter(w http.ResponseWriter) *CustomResponseWriter {
	return &CustomResponseWriter{
		ResponseWriter: w,
	}
}

func (r *CustomResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Write does the same thing as default http.ResponseWriter.Write method, but also stores the response body if status code is not 2XX
func (r *CustomResponseWriter) Write(msg []byte) (int, error) {
	if r.status > 299 {
		r.msg = append(r.msg, msg...)
	}

	return r.ResponseWriter.Write(msg)
}

// Status returns the status code of the request
func (r *CustomResponseWriter) Status() int {
	return r.status
}

// StringMsg returns a string representation of the response body
func (r *CustomResponseWriter) StringMsg() string {
	return string(r.msg)
}
