package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kicodelibrary/go-http-server-2022/2_persistence/pkg/server"
	"github.com/kicodelibrary/go-http-server-2022/2_persistence/pkg/server/database/bolt"
)

func main() {
	address := ":8080"
	r := mux.NewRouter()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	b, err := bolt.New(ctx, "./data")
	if err != nil {
		log.Fatalf("Failed to start database: %v", err)
	}
	defer b.Close(ctx)
	srv := server.New(ctx, b)

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
