// Command server runs the TodoList REST API.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ygddabong/TodoList/todo"
)

func main() {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           logRequests(todo.NewHandler(todo.NewStore())),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}
