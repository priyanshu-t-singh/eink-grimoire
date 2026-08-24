package core

import (
	"le-grimoire/internal/device"
	"le-grimoire/internal/kavita"
)

func (a *App) InitRepositories() {
	// Initialize Kavita Repository
	a.KavitaRepository = &kavita.Repository{}

	// Initialize Device Repository
	a.DeviceRepository = device.NewDeviceRepository(a.Database)
}
