package main

import "net/http"

func handlerReadiness(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(http.StatusText(http.StatusOK)))
}