package local

import (
	"context"
	"encoding/binary"
	"math"

	"github.com/go-ble/ble"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/timoth-y/chainmetric-core/models"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

// LocationTopic defines Bluetooth service for tethering geo-location data.
type LocationTopic struct {
	*ble.Service

	ready chan models.Location
}

// NewLocationTopic creates instance of the LocationTopic.
func NewLocationTopic() *LocationTopic {
	lt := &LocationTopic{
		Service: ble.NewService(ble.MustParse(viper.GetString("bluetooth.location.service_uuid"))),
		ready: make(chan models.Location, 1),
	}

	lt.AddCharacteristic(lt.eastCoordinate())
	lt.AddCharacteristic(lt.northCoordinate())
	lt.AddCharacteristic(lt.locationName())

	return lt
}

// Join subscribes to the messages on local network related to "location" topic.
func (lt *LocationTopic) Join(ctx context.Context, handler func(location models.Location) error) error {
	for {
		select {
		case location := <-lt.ready:
			go func(l models.Location) {
				if err := handler(l); err != nil {
					shared.Logger.Error(err)
				}
			}(location)
		case <- ctx.Done():
			switch ctx.Err() {
			case context.DeadlineExceeded:
				return errors.New("timeout waiting for location message")
			default:
				shared.Logger.Debug("Local network 'location' topic listener ended.")
				return nil
			}
		}
	}
}

func (lt *LocationTopic) eastCoordinate() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB1))
	char.HandleWrite(ble.WriteHandlerFunc(lt.handleWriteEast))

	return
}

func (lt *LocationTopic) northCoordinate() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB0))
	char.HandleWrite(ble.WriteHandlerFunc(lt.handleWriteNorth))

	return
}

func (lt *LocationTopic) locationName() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB5))
	char.HandleWrite(ble.WriteHandlerFunc(lt.handleWriteName))

	return
}

func (lt *LocationTopic) handleWriteEast(req ble.Request, rsp ble.ResponseWriter) {
	shared.Logger.Debugf("East Coordinate: Wrote %v", bytesToFloat64(req.Data()))
}

func (lt *LocationTopic) handleWriteNorth(req ble.Request, rsp ble.ResponseWriter) {
	shared.Logger.Debugf("North Coordinate: Wrote %v", bytesToFloat64(req.Data()))
}

func (lt *LocationTopic) handleWriteName(req ble.Request, rsp ble.ResponseWriter) {
	shared.Logger.Debugf("Location name: Wrote %s", string(req.Data()))
}

func bytesToFloat64(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}
