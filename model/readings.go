package model

import "github.com/timoth-y/iot-blockchain-contracts/models"

type MetricReading struct {
	Source string
	Value float64
}

type MetricReadings map[models.Metric] interface{}

type MetricReadingsPipe map[models.Metric] chan MetricReading

type MetricReadingsResults map[models.Metric] interface{}
