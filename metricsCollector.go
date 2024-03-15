package main

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Metric for tracking the last update timestamp of the Opsgenie data.
	lastUpdateTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "opsgenie_last_update_timestamp_seconds",
		Help: "Timestamp of the last data update in Opsgenie exporter.",
	})

	// Metric for counting Opsgenie users by various characteristics such as total, blocked, and unverified.
	opsgenieUsers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opsgenie_users",
		Help: "Counts of Opsgenie users by various characteristics.",
	}, []string{"key"})

	// Metric for tracking the verification status of Opsgenie users. Uses 0 for unverified, 1 for verified.
	opsgenieUserVerifiedStatus = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opsgenie_user_verified_status",
		Help: "Verification status of Opsgenie users. 0 for unverified, 1 for verified.",
	}, []string{"username"})

	// Metric for counting the total number of Opsgenie teams.
	opsgenieTeamsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "opsgenie_teams_total",
		Help: "Total number of Opsgenie teams.",
	})

	// Metric for displaying information about the Opsgenie account, such as user count, max user count, and yearly plan status.
	opsgenieAccount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opsgenie_account",
		Help: "Information about the Opsgenie account.",
	}, []string{"key"})

	// Metric for counting Opsgenie integrations by type.
	opsgenieIntegrationsTotal = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opsgenie_integrations_total",
		Help: "Total number of Opsgenie integrations by type.",
	}, []string{"type"})

	// Metric for counting the total number of Opsgenie heartbeats.
	opsgenieHeartbeatsTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "opsgenie_heartbeats_total",
		Help: "Total number of Opsgenie heartbeats.",
	})

	// Metric for counting the total number of enabled Opsgenie heartbeats.
	opsgenieHeartbeatsEnabledTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "opsgenie_heartbeats_enabled_total",
		Help: "Total number of enabled Opsgenie heartbeats.",
	})

	// Metric for indicating whether an Opsgenie heartbeat is expired or not, labeled by team.
	opsgenieHeartbeatsExpired = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "opsgenie_heartbeats_expired",
		Help: "Indicates whether an Opsgenie heartbeat is expired (1) or not (0), labeled by team.",
	}, []string{"team"})
)

func init() {
	// Register the metrics with Prometheus client to start collecting them.
	prometheus.MustRegister(lastUpdateTimestamp)
	prometheus.MustRegister(opsgenieUsers)
	prometheus.MustRegister(opsgenieUserVerifiedStatus)
	prometheus.MustRegister(opsgenieTeamsTotal)
	prometheus.MustRegister(opsgenieAccount)
	prometheus.MustRegister(opsgenieIntegrationsTotal)
	prometheus.MustRegister(opsgenieHeartbeatsTotal)
	prometheus.MustRegister(opsgenieHeartbeatsEnabledTotal)
	prometheus.MustRegister(opsgenieHeartbeatsExpired)
}

// Helper function to convert a bool to float64 since Prometheus metrics need to be numerical.
func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}

// updateMetrics concurrently updates various Prometheus metrics based on data fetched from Opsgenie.
func updateMetrics(client *OpsgenieClient) {
	var wg sync.WaitGroup
	wg.Add(5) // Waiting for five operations: users, teams, account, integrations, and heartbeats.

	// Fetch and update user-related metrics in a goroutine.
	go func() {
		defer wg.Done() // Ensure the WaitGroup counter decreases on goroutine completion.
		users, err := client.ListUsers()
		if err != nil {
			log.Printf("Error fetching users: %v", err)
			return
		}

		var total, blocked, unverified int
		// Iterate through users to count total, blocked, and unverified users.
		for _, user := range users {
			total++
			if user.Blocked {
				blocked++
			}
			if !user.Verified {
				unverified++
			}
		}
		// Set the Prometheus gauge values for users metrics.
		opsgenieUsers.With(prometheus.Labels{"key": "total"}).Set(float64(total))
		opsgenieUsers.With(prometheus.Labels{"key": "blocked"}).Set(float64(blocked))
		opsgenieUsers.With(prometheus.Labels{"key": "unverified"}).Set(float64(unverified))

		// Additionally, set the verification status for each user.
		for _, user := range users {
			verifiedStatus := 1
			if !user.Verified {
				verifiedStatus = 0
				opsgenieUserVerifiedStatus.WithLabelValues(user.Username).Set(float64(verifiedStatus))
			}
		}
	}()

	// Fetch and update teams-related metrics in a goroutine.
	go func() {
		defer wg.Done()
		teams, err := client.ListTeams()
		if err != nil {
			log.Printf("Error fetching teams: %v", err)
			return
		}
		// Set the total number of teams in a Prometheus gauge.
		opsgenieTeamsTotal.Set(float64(len(teams)))
	}()

	// Concurrently fetch account information and update account-related metrics.
	go func() {
		defer wg.Done()
		accountInfo, err := client.GetAccountInfo()
		if err != nil {
			log.Printf("Error fetching account info: %v", err)
			return
		}

		// Set metrics based on the fetched account information.
		opsgenieAccount.With(prometheus.Labels{"key": "userCount"}).Set(float64(accountInfo.Data.UserCount))
		opsgenieAccount.With(prometheus.Labels{"key": "maxUserCount"}).Set(float64(accountInfo.Data.Plan.MaxUserCount))
		opsgenieAccount.With(prometheus.Labels{"key": "isYearly"}).Set(boolToFloat64(accountInfo.Data.Plan.IsYearly))
	}()

	// Fetch and update integrations-related metrics in a goroutine.
	go func() {
		defer wg.Done()
		integrations, err := client.ListIntegrations()
		if err != nil {
			log.Printf("Error fetching integrations: %v", err)
			return
		}
		// Count and set the number of integrations by type.
		typeCounts := make(map[string]int)
		for _, integration := range integrations {
			typeCounts[integration.Type]++
		}
		for integrationType, count := range typeCounts {
			opsgenieIntegrationsTotal.WithLabelValues(integrationType).Set(float64(count))
		}
	}()

	// Concurrently fetch heartbeats information and update heartbeats-related metrics.
	go func() {
		defer wg.Done()
		heartbeats, err := client.ListHeartbeats()
		if err != nil {
			log.Printf("Error fetching heartbeats: %v", err)
			return
		}

		var enabledCount, totalHeartbeats int
		expiredAndEnabledCountByTeam := make(map[string]int)
		// Iterate through heartbeats to count enabled and identify expired ones by team.
		for _, hb := range heartbeats {
			totalHeartbeats++
			detail, err := client.GetHeartbeatDetail(hb.Name)
			if err != nil {
				log.Printf("Error fetching heartbeat detail for %s: %v", hb.Name, err)
				continue
			}
			if detail.Enabled {
				enabledCount++
				if detail.Expired {
					teamName := detail.OwnerTeam.Name
					if teamName == "" {
						teamName = "no_team" // Default value for heartbeats without a team.
					}
					expiredAndEnabledCountByTeam[teamName]++
				}
			}
		}

		// Update metrics for enabled and total heartbeats.
		opsgenieHeartbeatsEnabledTotal.Set(float64(enabledCount))
		// Update metrics for expired heartbeats by team.
		for team, count := range expiredAndEnabledCountByTeam {
			opsgenieHeartbeatsExpired.WithLabelValues(team).Set(float64(count))
		}
		opsgenieHeartbeatsTotal.Set(float64(totalHeartbeats))
	}()

	wg.Wait()                              // Wait for all goroutines to complete before updating the last update timestamp.
	lastUpdateTimestamp.SetToCurrentTime() // Update the timestamp metric to the current time after metrics updates.
}
