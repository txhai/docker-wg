package main

import (
	"docker_wg/api"
	"docker_wg/loggin"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	EnvHost      = "HOST"
	EnvPort      = "PORT"
	EnvSubnet    = "INTERNAL_SUBNET"
	EnvKeepAlive = "KEEP_ALIVE"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(contentTypeHeader, contentTypeJSON)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// get env variables
	host := os.Getenv(EnvHost)
	port := os.Getenv(EnvPort)
	subnet := os.Getenv(EnvSubnet)
	keepAlive := os.Getenv(EnvKeepAlive)

	if keepAlive == "" {
		keepAlive = "0"
	}

	// create logger
	logger := loggin.NewLogger(loggin.LevelVerbose)

	var err error

	// create api
	handler, err := api.NewApi(subnet, keepAlive, logger)
	if err != nil {
		log.Panicf("create api handler %v", err)
		return
	}

	// create http server
	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.UseEncodedPath()
	router.HandleFunc("/{itf}/server_key", handler.ServerKeyHandler).Methods("GET")
	router.HandleFunc("/{itf}/add_peer", handler.AddPeerHandler).Methods("POST")
	router.HandleFunc("/{itf}/remove_peer", handler.RemovePeerHandler).Methods("DELETE").Queries("key", "{key}")
	router.HandleFunc("/{itf}/list_peer", handler.ListPeerHandler).Methods("GET")
	router.HandleFunc("/{itf}/healthcheck", handler.HealthCheckHandler).Methods("GET")

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%s", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Panicf("listen and serve %v", err)
		return
	}
}
