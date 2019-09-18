package gauges

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type slots struct {
	DatabaseName string  `db:"database_name"`
	SlotName     string  `db:"slot_name"`
	Active       float64 `db:"active"`
	TotalLag     float64 `db:"total_lag"`
}

// ReplicationSlotStatus returns the state of the replication slots
func (g *Gauges) ReplicationSlotStatus() *prometheus.GaugeVec {
	var gauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:        "postgresql_replication_slot_status",
			Help:        "State of the replication slots",
			ConstLabels: g.labels,
		},
		[]string{"database_name", "slot_name"},
	)
	go func() {
		for {
			gauge.Reset()
			var slots []slots
			if err := g.query(
				`
					SELECT
						"database" AS database_name,
						slot_name,
						active::int
					FROM pg_replication_slots
					WHERE slot_type = 'logical';
				`,
				&slots,
				emptyParams,
			); err == nil {
				for _, slot := range slots {
					gauge.With(prometheus.Labels{
						"database_name": lock.DatabaseName,
						"slot_name":     lock.SlotName,
					}).Set(slot.Active)
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
		[]string{"database_name", "slot_name"},
	)
	go func() {
		for {
			gauge.Reset()
			var slots []slots
			if err := g.query(
				`
					SELECT
						"database" AS database_name,
						slot_name,
						round(pg_wal_lsn_diff(pg_current_wal_lsn(), confirmed_flush_lsn) / 1048576, 0) AS total_lag
					FROM pg_replication_slots
					WHERE slot_type = 'logical';
				`,
				&slots,
				emptyParams,
			); err == nil {
				for _, slot := range slots {
					gauge.With(prometheus.Labels{
						"database_name": lock.DatabaseName,
						"slot_name":     lock.SlotName,
					}).Set(slot.TotalLag)
				}
			}
			time.Sleep(g.interval)
		}
	}()
	return gauge
}
