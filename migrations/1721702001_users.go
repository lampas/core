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
		Name:     em.UserTableName,
		Type:     models.CollectionTypeBase,
		System:   false,
		ListRule: types.Pointer(""),
		Schema: schema.NewSchema(
			// Username
			&schema.SchemaField{
				System:      false,
				Name:        "username",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      true,
				Options:     schema.TextOptions{},
			},

			// External ID
			&schema.SchemaField{
				System:      false,
				Name:        "externalId",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      true,
				Options:     schema.TextOptions{},
			},

			// Name
			&schema.SchemaField{
				System:      false,
				Name:        "name",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.TextOptions{},
			},

			// TelegramUsername
			&schema.SchemaField{
				System:      false,
				Name:        "telegramUsername",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.TextOptions{},
			},

			// Role
			&schema.SchemaField{
				System:      false,
				Name:        "role",
				Type:        schema.FieldTypeSelect,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.SelectOptions{
					MaxSelect: 1,
					Values:    []string{"guest", "user", "admin"},
				},
			},

			// Lang
			&schema.SchemaField{
				System:      false,
				Name:        "lang",
				Type:        schema.FieldTypeSelect,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.SelectOptions{
					MaxSelect: 1,
					Values:    []string{"uk", "ru"},
				},
			},

			// UserSetLang
			&schema.SchemaField{
				System:      false,
				Name:        "userSetLang",
				Type:        schema.FieldTypeBool,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.BoolOptions{},
			},

			// UsernameChangedAt
			&schema.SchemaField{
				System:      false,
				Name:        "usernameChangedAt",
				Type:        schema.FieldTypeDate,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.DateOptions{},
			},
		),
	}

	collection.Id = em.UserCollectionId

	m.Register(func(db dbx.Builder) error {
		return daos.New(db).ImportCollections([]*models.Collection{collection}, false, nil)
	}, func(db dbx.Builder) error {
		daos.New(db).DeleteCollection(collection)
		return nil
	})
}
