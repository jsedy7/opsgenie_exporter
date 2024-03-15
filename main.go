package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Define command-line flags with specified names and defaults.
	httpPort := flag.Int("http.port", 8080, "The port to run the Prometheus metrics server on")
	refreshInterval := flag.Int("refresh", 600, "Interval between metric updates in seconds")
	flag.Parse()

	client := NewOpsgenieClient(os.Getenv("OPSGENIE_API_KEY"))

	// Convert refreshInterval to time.Duration for use with time.Sleep
	updateInterval := time.Duration(*refreshInterval) * time.Second

	// Start a separate goroutine for periodically updating metrics
	go func() {
		for {
			log.Println("Starting metrics update")
			updateMetrics(client)
			log.Println("Metrics update completed")
			time.Sleep(updateInterval) // Wait for the specified interval before updating again
		}
	}()

	// Set up the HTTP server to expose Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting Prometheus metrics server on port %d\n", *httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), nil))
}
