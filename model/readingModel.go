package model

type MetricReading struct {
	Source string
	Value interface{}
}

type MetricReadings map[Metric] MetricReading

type MetricReadingsPipe map[Metric] chan MetricReading
