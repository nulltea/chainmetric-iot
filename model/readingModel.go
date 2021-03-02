package model

type MetricReading struct {
	Source string
	Value interface{}
}

type MetricReadings map[Metric] interface{}

type MetricReadingsPipe map[Metric] chan MetricReading

type MetricReadingsResults map[Metric] interface{}
