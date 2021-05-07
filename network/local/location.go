package local

import (
	"github.com/go-ble/ble"
	"github.com/spf13/viper"

	"github.com/timoth-y/chainmetric-sensorsys/shared"
)

type LocationService struct {
	*ble.Service
}

func NewLocationService() *LocationService {
	ls := &LocationService{
		Service: ble.NewService(ble.MustParse(viper.GetString("bluetooth.location.service_uuid"))),
	}

	ls.AddCharacteristic(ls.eastCoordinate())
	ls.AddCharacteristic(ls.northCoordinate())
	ls.AddCharacteristic(ls.locationName())

	return ls
}

func (ls *LocationService) eastCoordinate() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB1))
	char.HandleWrite(ble.WriteHandlerFunc(ls.handleWriteEast))

	return
}

func (ls *LocationService) northCoordinate() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB0))
	char.HandleWrite(ble.WriteHandlerFunc(ls.handleWriteNorth))

	return
}

func (ls *LocationService) locationName() (char *ble.Characteristic) {
	char = ble.NewCharacteristic(ble.UUID16(0x2AB5))
	char.HandleWrite(ble.WriteHandlerFunc(ls.handleWriteName))

	return
}

func (ls *LocationService) handleWriteEast(req ble.Request, rsp ble.ResponseWriter) {
	shared.Logger.Debugf("East Coordinate: Wrote %s", string(req.Data()))
}

func (ls *LocationService) handleWriteNorth(req ble.Request, rsp ble.ResponseWriter) {
	shared.Logger.Debugf("North Coordinate: Wrote %s", string(req.Data()))
}

func (ls *LocationService) handleWriteName(req ble.Request, rsp ble.ResponseWriter) {
	shared.Logger.Debugf("Location name: Wrote %s", string(req.Data()))
}

