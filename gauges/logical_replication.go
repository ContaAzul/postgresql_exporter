package gauges

import (
	"fmt"
	"time"

	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
)

type slots struct {
	SlotName     string  `db:"slot_name"`
	IsSlotActive float64 `db:"active"`
	SlotTotalLag float64 `db:"total_lag"`
}

// ReplicationSlotStatus returns the state of the replication slots
func (g *Gauges) ReplicationSlotStatus() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_slot_status",
			Help:        "State of the replication slots",
			ConstLabels: g.labels,
		},
		[]string{"slot_name"},
	)
	go func() {
		for {
			gauge.Reset()
			var slots []slots
			if err := g.query(
				`
					SELECT
						slot_name,
						active::int
					FROM pg_replication_slots
					WHERE slot_type = 'logical'
					  AND "database" = current_database();
				`,
				&slots,
				emptyParams,
			); err == nil {
				for _, slot := range slots {
					gauge.With(prometheus.Labels{
						"slot_name": slot.SlotName,
					}).Set(slot.IsSlotActive)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

// ReplicationSlotLagInMegabytes returns the total lag from the replication slots
func (g *Gauges) ReplicationSlotLagInMegabytes() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_slot_lag",
			Help:        "Total lag of the replication slots",
			ConstLabels: g.labels,
		},
		[]string{"slot_name"},
	)
	go func() {
		for {
			gauge.Reset()
			var slots []slots
			if err := g.query(
				fmt.Sprintf(
					`
						SELECT
							slot_name,
							round(%s(%s(), confirmed_flush_lsn) / 1048576, 0) AS total_lag
						FROM pg_replication_slots
						WHERE slot_type = 'logical'
						AND "database" = current_database();
					`,
					postgres.Version(g.version()).WalLsnDiffFunctionName(),
					postgres.Version(g.version()).CurrentWalLsnFunctionName(),
				),
				&slots,
				emptyParams,
			); err == nil {
				for _, slot := range slots {
					gauge.With(prometheus.Labels{
						"slot_name": slot.SlotName,
					}).Set(slot.SlotTotalLag)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}
