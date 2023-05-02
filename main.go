package main

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareOne")
		next.ServeHTTP(w, r)
		log.Print("Executing middlewareOne again")
	})
}

func middlewareTwo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing middlewareTwo")
		if r.URL.Path == "/foo" {
			return
		}

		next.ServeHTTP(w, r)
		log.Print("Executing middlewareTwo again")
	})
}

func newLoggingHandler(dst io.Writer) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(dst, h)
	}
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func finalbas(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func enforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}

			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	logFile, err := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(finalbas)
	//authHandler := httpauth.SimpleBasicAuth("alice", "pa$$word")

	//mux.Handle("/", middlewareOne(middlewareTwo(finalHandler)))
	//mux.Handle("/", enforceJSONHandler(finalHandler))
	//mux.Handle("/", authHandler(finalHandler))
	mux.Handle("/", handlers.LoggingHandler(logFile, finalHandler))

	log.Print("Listening on :3000...")
	error := http.ListenAndServe(":3000", mux)
	log.Fatal(error)
}
