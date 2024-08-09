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
		Name:     em.TemplateTableName,
		Type:     models.CollectionTypeBase,
		System:   false,
		ListRule: types.Pointer(""),
		Schema: schema.NewSchema(
			// Slug
			&schema.SchemaField{
				System:      false,
				Name:        "slug",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      true,
				Options:     schema.TextOptions{},
			},

			// Banner
			&schema.SchemaField{
				System:      false,
				Name:        "banner",
				Type:        schema.FieldTypeFile,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.FileOptions{
					MimeTypes: []string{
						"image/png",
						"image/jpeg",
						"image/webp",
					},
					MaxSelect: 1,
					MaxSize:   52428800,
					Protected: false,
				},
			},

			// ContentRu
			&schema.SchemaField{
				System:      false,
				Name:        "contentRu",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.TextOptions{},
			},

			// ContentUk
			&schema.SchemaField{
				System:      false,
				Name:        "contentUk",
				Type:        schema.FieldTypeText,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options:     schema.TextOptions{},
			},
		),
	}

	collection.Id = em.TemplateCollectionId

	m.Register(func(db dbx.Builder) error {
		return daos.New(db).ImportCollections([]*models.Collection{collection}, false, nil)
	}, func(db dbx.Builder) error {
		daos.New(db).DeleteCollection(collection)
		return nil
	})
}
