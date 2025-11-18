package main

import "net/http"

func main() {
	sm := http.NewServeMux()
	server := http.Server{
		Handler: sm,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}
