package gauges

import (
	"strings"

	"github.com/apex/log"
)

func (g *Gauges) version() string {
	var version string
	if err := g.db.QueryRow("show server_version").Scan(&version); err != nil {
		log.WithError(err).Error("failed to get postgresql version, assuming 9.6.0")
		return "9.6.0"
	}
	return version
}

func isPG96(version string) bool {
	return strings.HasPrefix(version, "9.6.")
}

func isPG10(version string) bool {
	return strings.HasPrefix(version, "10.")
}
