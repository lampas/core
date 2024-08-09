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
		Name:     em.ConfigurationTableName,
		Type:     models.CollectionTypeBase,
		System:   false,
		ListRule: types.Pointer(""),
		Schema: schema.NewSchema(
			// BotToken
			&schema.SchemaField{
				System:      false,
				Name:        "botToken",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.TextOptions{},
			},

			// BotUsername
			&schema.SchemaField{
				System:      false,
				Name:        "botUsername",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.TextOptions{},
			},

			// BotNameRu
			&schema.SchemaField{
				System:      false,
				Name:        "botNameRu",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     nil,
					Max:     types.Pointer(64),
					Pattern: "",
				},
			},

			// BotNameUk
			&schema.SchemaField{
				System:      false,
				Name:        "botNameUk",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     nil,
					Max:     types.Pointer(64),
					Pattern: "",
				},
			},

			// BotDescriptionRu
			&schema.SchemaField{
				System:      false,
				Name:        "botDescriptionRu",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     nil,
					Max:     types.Pointer(512),
					Pattern: "",
				},
			},

			// BotDescriptionUk
			&schema.SchemaField{
				System:      false,
				Name:        "botDescriptionUk",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     nil,
					Max:     types.Pointer(512),
					Pattern: "",
				},
			},

			// BotShortDescriptionRu
			&schema.SchemaField{
				System:      false,
				Name:        "botShortDescriptionRu",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     nil,
					Max:     types.Pointer(120),
					Pattern: "",
				},
			},

			// BotShortDescriptionUk
			&schema.SchemaField{
				System:      false,
				Name:        "botShortDescriptionUk",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     nil,
					Max:     types.Pointer(120),
					Pattern: "",
				},
			},

			// AppDefaultLang
			&schema.SchemaField{
				System:      false,
				Name:        "appDefaultLang",
				Type:        schema.FieldTypeSelect,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.SelectOptions{
					MaxSelect: 1,
					Values:    []string{"Uk", "ru"},
				},
			},

			// AppTitleRu
			&schema.SchemaField{
				System:      false,
				Name:        "appTitleRu",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min: types.Pointer(1),
					Max: types.Pointer(64),
				},
			},

			// AppTitleUk
			&schema.SchemaField{
				System:      false,
				Name:        "appTitleUk",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min: types.Pointer(1),
					Max: types.Pointer(64),
				},
			},

			// AppDescriptionRu
			&schema.SchemaField{
				System:      false,
				Name:        "appDescriptionRu",
				Type:        schema.FieldTypeEditor,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.EditorOptions{},
			},

			// AppDescriptionUk
			&schema.SchemaField{
				System:      false,
				Name:        "appDescriptionUk",
				Type:        schema.FieldTypeEditor,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.EditorOptions{},
			},

			// CookieEncryptionKey
			&schema.SchemaField{
				System:      false,
				Name:        "cookieEncryptionKey",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min: types.Pointer(32),
					Max: types.Pointer(32),
				},
			},
		),
	}

	collection.Id = em.ConfigurationCollectionId

	m.Register(func(db dbx.Builder) error {
		return daos.New(db).ImportCollections([]*models.Collection{collection}, false, nil)
	}, func(db dbx.Builder) error {
		daos.New(db).DeleteCollection(collection)
		return nil
	})
}
