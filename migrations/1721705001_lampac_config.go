package migrations

import (
	em "lamas/models"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	collection := &models.Collection{
		Name:     em.LampacConfigTableName,
		Type:     models.CollectionTypeBase,
		System:   false,
		ListRule: types.Pointer(""),
		Schema: schema.NewSchema(
			// Enabled
			&schema.SchemaField{
				System:      false,
				Name:        "enabled",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      true,
				Options:     schema.BoolOptions{},
			},

			// AdultEnabled
			&schema.SchemaField{
				System:      false,
				Name:        "adultEnabled",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      true,
				Options:     schema.BoolOptions{},
			},

			// JacketEnabled
			&schema.SchemaField{
				System:      false,
				Name:        "jacketEnabled",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// TorrServerEnabled
			&schema.SchemaField{
				System:      false,
				Name:        "torrServerEnabled",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// DLNAEnabled
			&schema.SchemaField{
				System:      false,
				Name:        "dlnaEnabled",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// TracksModifyEnabled
			&schema.SchemaField{
				System:      false,
				Name:        "tracksModifyEnabled",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// CheckOnlineSearch
			&schema.SchemaField{
				System:      false,
				Name:        "checkOnlineSearch",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// MultiAccess
			&schema.SchemaField{
				System:      false,
				Name:        "multiAccess",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// LocalIP
			&schema.SchemaField{
				System:      false,
				Name:        "localIP",
				Type:        schema.FieldTypeBool,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// ListenScheme
			&schema.SchemaField{
				System:      false,
				Name:        "listenScheme",
				Type:        schema.FieldTypeSelect,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.SelectOptions{
					MaxSelect: 1,
					Values:    []string{"https", "http"},
				},
			},

			// ImageProxyService
			&schema.SchemaField{
				System:      false,
				Name:        "imageProxyService",
				Type:        schema.FieldTypeSelect,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.SelectOptions{
					MaxSelect: 1,
					Values:    []string{"lampac", "imgproxy"},
				},
			},

			// TorSocksProxy
			&schema.SchemaField{
				System:      false,
				Name:        "torSocksProxy",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Pattern: "socks5://[0-9a-zA-Z.]+:[0-9]+",
				},
			},

			// OverrideInitConfig
			&schema.SchemaField{
				System:      false,
				Name:        "overrideInitConfig",
				Type:        schema.FieldTypeJson,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.JsonOptions{},
			},

			// OverrideModuleManifest
			&schema.SchemaField{
				System:      false,
				Name:        "overrideModuleManifest",
				Type:        schema.FieldTypeJson,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.JsonOptions{},
			},
		),
	}

	collection.Id = em.LampacConfigCollectionId

	m.Register(func(db dbx.Builder) error {
		return daos.New(db).ImportCollections([]*models.Collection{collection}, false, nil)
	}, func(db dbx.Builder) error {
		daos.New(db).DeleteCollection(collection)
		return nil
	})
}
