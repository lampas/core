package migrations

import (
	"lamas/models"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		configuration := &models.LampacConfig{}
		configuration.Id = models.LampacConfigSingleRecordId

		configuration.Enabled = true
		configuration.AdultEnabled = false
		configuration.JacketEnabled = false
		configuration.TorrServerEnabled = false
		configuration.DLNAEnabled = false
		configuration.TracksModifyEnabled = false

		configuration.CheckOnlineSearch = true
		configuration.MultiAccess = true
		configuration.LocalIP = false
		configuration.ListenScheme = "https"
		configuration.ImageProxyService = "imgproxy"

		configuration.TorSocksProxy = "socks5://tor-proxy:9050"
		configuration.OverrideInitConfig = types.JsonRaw(`{}`)
		configuration.OverrideModuleManifest = []any{}

		return dao.Save(configuration)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		configuration := models.GetConfiguration(dao)
		if configuration.HasId() {
			return dao.Delete(configuration)
		}

		return nil
	})
}
