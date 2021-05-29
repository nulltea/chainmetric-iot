package sensor

import (
	"github.com/timoth-y/chainmetric-core/models"
)

// SensorsRegister represents pool of the multiply Sensor devices.
type SensorsRegister map[string]Sensor

// SupportedMetrics aggregates all supported by sensors models.Metric devices.
func (sr SensorsRegister) SupportedMetrics() models.Metrics {
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

// ToList returns slice of all Sensor devices presented in SensorsRegister.
func (sr SensorsRegister) ToList() []Sensor {
	var (
		l = make([]Sensor, len(sr))
		i = 0
	)

	for id := range sr {
		l[i] = sr[id]
		i++
	}

	return l
}

// NotEmpty determines whether SensorsRegister contains at least one Sensor.
func (sr SensorsRegister) NotEmpty() bool {
	return len(sr) > 0
}

// Exists determines whether the Sensor exists in SensorsRegister by given `id`.
func (sr SensorsRegister) Exists(id string) bool {
	_, is := sr[id]
	return is
}
