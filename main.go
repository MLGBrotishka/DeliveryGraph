package main

import (
	"graph/database"
	"graph/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	database.ConnectRedisDB()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		http.HandleFunc("/v1/path", handlers.GetV1Path)
		http.HandleFunc("/v1/path/multiple-couriers", handlers.GetV1PathMultipleCouriers)
		http.HandleFunc("/v1/point/is-available", handlers.GetV1PointIsAvailable)
		log.Println("Server started on http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
	<-stop
	database.CloseRedisDB()
	log.Println("Server stopped")
}
