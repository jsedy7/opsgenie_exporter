package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	client := NewOpsgenieClient(os.Getenv("OPSGENIE_API_KEY"))

	// Spuštění pravidelné aktualizace metrik v samostatné gorutině
	go func() {
		for {
			updateMetrics(client)
			time.Sleep(10 * time.Minute) // Aktualizace každých 10 minut, upravte podle potřeby
		}
	}()

	// Nastavení HTTP serveru pro vystavení metrik Prometheus
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
