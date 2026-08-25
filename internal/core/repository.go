package core

import (
	"context"
	"time"

	"le-grimoire/internal/device"
	"le-grimoire/internal/kavita"
)

func (a *App) InitRepositories() {
	// Initialize Kavita Repository
	kavitaClient := kavita.NewClient(kavita.Config{
		BaseURL:    a.Config.GetKavitaAPIURI(),
		APIKey:     a.Config.KavitaAPI.APIKey,
		PluginName: a.Config.KavitaAPI.PluginName,
		Timeout:    a.Config.KavitaAPI.Timeout,
	})
	kavitaRepo := kavita.NewRepository(kavitaClient)

	authCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := kavitaRepo.Authenticate(authCtx); err != nil {
		a.Logger.Warn("Initial Kavita auth failed (will retry reactively on requests)", "error", err)
	} else {
		a.Logger.Info("Successfully authenticated with Kavita server")
	}

	// Initialize Device Repository
	a.DeviceRepository = device.NewDeviceRepository(a.Database)
}
