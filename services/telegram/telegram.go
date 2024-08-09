package telegram

import (
	"lamas/models"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterService(app *pocketbase.PocketBase) {
	stop := func() error {
		return nil
	}

	restart := func() error {
		stop()
		stop = ConnectTelegram(app)
		return nil
	}

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return restart()
	})

	app.OnTerminate().Add(func(e *core.TerminateEvent) error {
		return stop()
	})

	app.OnModelAfterUpdate(models.ConfigurationCollectionId, models.LampacConfigCollectionId).
		Add(func(e *core.ModelEvent) error {
			return restart()
		})
}
