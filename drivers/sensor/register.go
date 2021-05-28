package sensor

import (
	"github.com/timoth-y/chainmetric-core/models"
)

// SensorsRegister represents pool of the multiply Sensor devices.
type SensorsRegister map[string]Sensor

// SupportedMetrics aggregates all supported by sensors models.Metric devices.
func (sr SensorsRegister) SupportedMetrics() []models.Metric {
	var (
		availableMetrics = make(map[models.Metric]int)
	)

	for _, s := range sr {
		for _, metric := range s.Metrics() {
			availableMetrics[metric]++
		}
	}

	var (
		metrics = make([]models.Metric, len(availableMetrics))
		i       = 0
	)

	for m, _ := range availableMetrics {
		metrics[i] = m
		i++
	}
	
	return metrics
}

// Union produces new SensorsRegister combining sensors from original and `sr2`.
func (sr SensorsRegister) Union(sr2 SensorsRegister) SensorsRegister {
	sru := SensorsRegister{}

	for id, s := range sr {
		sru[id] = s
	}

	for id, s := range sr2 {
		sru[id] = s
	}

	return sr
}
