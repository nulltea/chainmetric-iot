package model

import "github.com/timoth-y/iot-blockchain-contracts/models"

type MetricReading struct {
	Source string
	Value interface{}
}

type MetricReadings map[models.Metric] interface{}

type MetricReadingsPipe map[models.Metric] chan MetricReading

type MetricReadingsResults map[models.Metric] interface{}
