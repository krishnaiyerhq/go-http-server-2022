package main

import (
	"log"
	"net/http"
	"time"

	"github.com/KrishnaIyer/goexamples/1_http/pkg/server"
	"github.com/gorilla/mux"
)

func main() {
	address := ":8080"
	r := mux.NewRouter()
	srv := server.New()
	r.HandleFunc("/", srv.HandleIndex)
	r.HandleFunc("/users/create", srv.HandleCreateUsers)
	r.HandleFunc("/users/{name}", srv.HandleUsers)
	s := &http.Server{
		Addr:           address,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Start server: %v", address)
	log.Fatal(s.ListenAndServe())
}
