package main

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Definice metriky pro počet uživatelů
	opsgenieUsers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opsgenie_users",
		Help: "Counts of Opsgenie users by various characteristics.",
	}, []string{"key"})
	// Definice metriky pro počet týmů
	opsgenieTeamsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "opsgenie_teams_total",
		Help: "Total number of Opsgenie teams.",
	})
	// Definice metriky pro účet
	opsgenieAccount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opsgenie_account",
		Help: "Information about the Opsgenie account.",
	}, []string{"key"})
)

func init() {
	// Registrace metriky, aby Prometheus klient věděl, že ji má sbírat
	prometheus.MustRegister(opsgenieUsers)
	prometheus.MustRegister(opsgenieTeamsTotal)
	prometheus.MustRegister(opsgenieAccount)
}

// Pomocná funkce pro konverzi bool na float64, protože Prometheus metriky musí být číselné
func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// Funkce pro aktualizaci metriky uživatelů
func updateMetrics(client *OpsgenieClient) {
	var wg sync.WaitGroup
	wg.Add(3) // Čekáme na dvě operace

	go func() {
		defer wg.Done()
		users, err := client.ListUsers()
		if err != nil {
			log.Printf("Error fetching users: %v", err)
			return
		}

		var total, blocked, unverified int
		for _, user := range users {
			total++
			if user.Blocked {
				blocked++
			}
			if !user.Verified {
				unverified++
			}
		}
		opsgenieUsers.With(prometheus.Labels{"key": "total"}).Set(float64(total))
		opsgenieUsers.With(prometheus.Labels{"key": "blocked"}).Set(float64(blocked))
		opsgenieUsers.With(prometheus.Labels{"key": "unverified"}).Set(float64(unverified))
	}()

	go func() {
		defer wg.Done()
		teams, err := client.ListTeams()
		if err != nil {
			log.Printf("Error fetching teams: %v", err)
			return
		}
		opsgenieTeamsTotal.Set(float64(len(teams)))
	}()

	// Paralelní získávání informací o účtu
	go func() {
		defer wg.Done()
		accountInfo, err := client.GetAccountInfo()
		if err != nil {
			log.Printf("Error fetching account info: %v", err)
			return
		}

		// Nastavení metrik na základě získaných informací o účtu
		opsgenieAccount.With(prometheus.Labels{"key": "userCount"}).Set(float64(accountInfo.Data.UserCount))
		opsgenieAccount.With(prometheus.Labels{"key": "maxUserCount"}).Set(float64(accountInfo.Data.Plan.MaxUserCount))
		opsgenieAccount.With(prometheus.Labels{"key": "isYearly"}).Set(boolToFloat64(accountInfo.Data.Plan.IsYearly))
	}()

	wg.Wait() // Počkáme, až obě Go rutiny dokončí
}
