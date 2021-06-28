package localnet

import (
	"context"
	"encoding/binary"
	"math"
	"sync"

	"github.com/go-ble/ble"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"
	"github.com/timoth-y/chainmetric-core/utils"

	"github.com/timoth-y/chainmetric-iot/drivers/periphery"
	"github.com/timoth-y/chainmetric-iot/shared"
)

var (
	geoOnce = sync.Once{}
)

// GeoLocationChannel defines Bluetooth service for tethering geo-location data.
type (
	GeoLocationChannel struct {
		service *ble.Service
		uuid ble.UUID

		receive chan func(*geoPayload)
		release chan models.Location
	}

	geoPayload struct {
		sync.Mutex
		lat *float64
		lng *float64
		name *string
	}
)

// newGeoLocationChannel creates instance of the GeoLocationChannel.
func newGeoLocationChannel() *GeoLocationChannel {
	return &GeoLocationChannel{
		receive: make(chan func(*geoPayload)),
		release: make(chan models.Location, 1),
	}
}

func (gc *GeoLocationChannel) init() {
	var (
		uuid = ble.MustParse(viper.GetString("bluetooth.location.service_uuid"))
	)

	gc.uuid = uuid
	gc.service = ble.NewService(uuid)
	gc.service.AddCharacteristic(gc.eastCoordinate())
	gc.service.AddCharacteristic(gc.northCoordinate())
	gc.service.AddCharacteristic(gc.locationName())
}

// Subscribe subscribes to the messages related to "geo" topic.
func (gc *GeoLocationChannel) Subscribe(ctx context.Context, handler func(location models.Location) error) error {
	// Run aggregator goroutine
	geoOnce.Do(func() {
		go func() {
			var payload = &geoPayload{}

			for {
				select {
				case setter := <- gc.receive:
					setter(payload)
					payload.tryRelease(gc.release)
				case <- ctx.Done():
					return
				}
			}
		}()
	})

	for {
		select {
		case location := <-gc.release:
			go func(l models.Location) {
				if err := handler(l); err != nil {
					shared.Logger.Error(err)
				}
			}(location)
		case <- ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return errors.New("timeout waiting for 'location' message")
			default:
				shared.Logger.Debug("Local network 'location' topic listener ended")
				return nil
			}
		}
	}
}

func (gc *GeoLocationChannel) expose(dev *periphery.Bluetooth) error {
	if err := dev.AddService(gc.service); err != nil {
		return errors.Wrap(err, "error adding location service to Bluetooth device")
	}

	return nil
}

func (gc *GeoLocationChannel) eastCoordinate() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB1))
	char.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		gc.receive <- func(p *geoPayload) {
			p.lat = utils.Float64Pointer(bytesToFloat64(req.Data()))
		}
	}))

	return
}

func (gc *GeoLocationChannel) northCoordinate() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB0))
	char.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		gc.receive <- func(p *geoPayload) {
			p.lng = utils.Float64Pointer(bytesToFloat64(req.Data()))
		}
	}))

	return
}

func (gc *GeoLocationChannel) locationName() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB5))
	char.HandleWrite(ble.WriteHandlerFunc(func(req ble.Request, rsp ble.ResponseWriter) {
		gc.receive <- func(p *geoPayload) {
			p.name = utils.StringPointer(string(req.Data()))
		}
	}))

	return
}

func (gp *geoPayload) complete() bool {
	if gp.lat != nil && gp.lng != nil && gp.name != nil {
		return true
	}

	return false
}

func (gp *geoPayload) tryRelease(ch chan models.Location) {
	if gp.complete() {
		ch <- models.Location{
			Latitude: *gp.lat,
			Longitude: *gp.lng,
			Name: *gp.name,
		}

		*gp = geoPayload{}
	}

	return
}

func bytesToFloat64(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
