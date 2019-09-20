package gauges

import (
	"fmt"
	"time"

	"github.com/ContaAzul/postgresql_exporter/postgres"
	"github.com/prometheus/client_golang/prometheus"
)

type slots struct {
	Name     string  `db:"slot_name"`
	Active   float64 `db:"active"`
	TotalLag float64 `db:"total_lag"`
}

// ReplicationSlotStatus returns the state of the replication slots
func (g *Gauges) ReplicationSlotStatus() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_slot_status",
			Help:        "Returns 1 if the slot is currently actively being used",
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
						"slot_name": slot.Name,
					}).Set(slot.Active)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}

// ReplicationSlotLagInBytes returns the total lag in bytes from the replication slots
func (g *Gauges) ReplicationSlotLagInBytes() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_lag_bytes",
			Help:        "Total lag in bytes of the replication slots",
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
						"slot_name": slot.Name,
					}).Set(slot.TotalLag)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}
