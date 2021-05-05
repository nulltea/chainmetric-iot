package shared

// InitCore performs core dependencies initialization sequence.
func InitCore() {
	initLogger()
	initConfig()
	initLevelDB()
	initPeriphery()
}

// CloseCore performs core dependencies close sequence.
func CloseCore() {
	LevelDB.Close()
	BluetoothDevice.Stop()
}
