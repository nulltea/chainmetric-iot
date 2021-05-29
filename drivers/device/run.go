package device

// Start performs startup of the Device and setting up all registered modules.
func (d *Device) Start() {
	d.modulesReg.Setup(d)
	d.modulesReg.Start(d.ctx)
}

func (d *Device) Close() error {
	d.active = false
	d.cancelDevice()

	return nil
}
