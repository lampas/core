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
		Name:     em.OutlineTokenTableName,
		Type:     models.CollectionTypeBase,
		System:   false,
		ListRule: types.Pointer(""),
		Schema: schema.NewSchema(
			// User
			&schema.SchemaField{
				System:      false,
				Name:        "user",
				Type:        schema.FieldTypeRelation,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.RelationOptions{
					CollectionId:  em.UserCollectionId,
					CascadeDelete: true,
					MinSelect:     nil,
					MaxSelect:     types.Pointer(1),
					DisplayFields: nil,
				},
			},

			// OutlineConfig
			&schema.SchemaField{
				System:      false,
				Name:        "outlineConfig",
				Type:        schema.FieldTypeRelation,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.RelationOptions{
					CollectionId:  em.OutlineConfigCollectionId,
					CascadeDelete: true,
					MinSelect:     nil,
					MaxSelect:     types.Pointer(1),
					DisplayFields: nil,
				},
			},

			// OutlineKey
			&schema.SchemaField{
				System:      false,
				Name:        "outlineKey",
				Type:        schema.FieldTypeJson,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options:     schema.JsonOptions{},
			},

			// BytesTransferred
			&schema.SchemaField{
				System:      false,
				Name:        "bytesTransferred",
				Type:        schema.FieldTypeNumber,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.NumberOptions{
					NoDecimal: true,
				},
			},
		),
	}

	collection.Id = em.OutlineTokenCollectionId

	m.Register(func(db dbx.Builder) error {
		return daos.New(db).ImportCollections([]*models.Collection{collection}, false, nil)
	}, func(db dbx.Builder) error {
		return daos.New(db).DeleteCollection(collection)
	})
}
