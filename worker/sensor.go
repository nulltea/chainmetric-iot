package worker

type Sensor interface {
	ID() string
	Init() error
	Harvest(ctx *Context)
	Close() error
}
