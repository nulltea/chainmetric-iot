package device

// Start performs startup of the Device and setting up all registered modules.
func (d *Device) Start() {
	d.modulesReg.Setup(d)
	d.modulesReg.Start(d.ctx)
}

// Close stops all working device.Module and frees allocated resources.
func (d *Device) Close() error {
	d.active = false
	d.cancelDevice()
	d.modulesReg.Close()

	return nil
}
