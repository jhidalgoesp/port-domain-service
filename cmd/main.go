package main

import (
	"github.com/jhidalgoesp/ports/internal/database"
	"github.com/jhidalgoesp/ports/internal/domain"
	"github.com/jhidalgoesp/ports/internal/file_reader"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	fileName = "ports.json"
	port     = ":8080"
)

func main() {
	log.Println("starting service")

	defer log.Println("shutdown complete")

	databaseAdapter := database.NewDatabase()

	fileAdapter, err := file_reader.NewFileReader(fileName)
	if err != nil {
		log.Panicln("fileAdapter could not be started:", err)
	}

	portService, err := domain.NewPortService(fileAdapter, databaseAdapter)
	if err != nil {
		log.Panicln("portService could not be started:", err)
	}

	// Handle requests with a handler function
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle the request and call portService.CreateOrUpdatePorts() here
		err := portService.CreateOrUpdatePorts()
		if err != nil {
			log.Println("Error handling request:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Respond to the request
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Ports updated successfully"))
		if err != nil {
			log.Println("Error handling request:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	// Create a channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the HTTP server in a goroutine
	go func() {
		log.Println("Server listening on", port)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Println("Failed to start server:", err)
			stop <- os.Interrupt
		}
	}()

	// Wait for a termination signal
	<-stop

	log.Println("Shutting down serve gracefully...")
}
