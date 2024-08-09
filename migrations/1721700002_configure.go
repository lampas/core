package migrations

import (
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		/** Settings */
		settings, _ := dao.FindSettings()

		settings.Meta.AppName = "Lampas"
		settings.Meta.HideControls = true

		settings.Backups.Cron = "0 0 * * *"
		settings.Backups.CronMaxKeep = 7

		settings.Logs.MaxDays = 31

		if err := dao.SaveSettings(settings); err != nil {
			return err
		}

		/** Admin */
		if os.Getenv("INITIALIZE_ADMIN_EMAIL") != "" {
			admin := &models.Admin{}
			admin.Email = os.Getenv("INITIALIZE_ADMIN_EMAIL")
			admin.SetPassword(os.Getenv("INITIALIZE_ADMIN_PASSWORD"))
			return dao.SaveAdmin(admin)
		}

		return nil
	}, nil)
}
