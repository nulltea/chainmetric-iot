package readings

import (
	"time"

	"sensorsys/model"
)

type ReceiverFunc func(model.MetricReadings)

type Receiver struct {
	Handler ReceiverFunc
	Metrics []model.Metric
	Period  time.Duration
}

func (r *Receiver) Request(ctx *Context) {
	for {


		time.Sleep(r.Period)
	}
}

func (r *Request) formRequest(ctx *Context) *Request {
	return &Request{
		Context: ctx,
		Metrics: r.Metrics,
	}
}

