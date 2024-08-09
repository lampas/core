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

func makeOutlineConfigsCollection(daos *daos.Dao) (*models.Collection, error) {
	collection := &models.Collection{
		Name:     em.OutlineConfigTableName,
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

			// Slug
			&schema.SchemaField{
				System:      false,
				Name:        "slug",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      true,
				Options: schema.TextOptions{
					Min:     types.Pointer(1),
					Max:     types.Pointer(32),
					Pattern: "[0-9a-zA-Z.]",
				},
			},

			// Template
			&schema.SchemaField{
				System:      false,
				Name:        "template",
				Type:        schema.FieldTypeRelation,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.RelationOptions{
					CollectionId:  em.TemplateCollectionId,
					CascadeDelete: false,
					MinSelect:     nil,
					MaxSelect:     types.Pointer(1),
					DisplayFields: nil,
				},
			},

			// TitleRu
			&schema.SchemaField{
				System:      false,
				Name:        "titleRu",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     types.Pointer(1),
					Max:     types.Pointer(64),
					Pattern: "",
				},
			},

			// TitleUk
			&schema.SchemaField{
				System:      false,
				Name:        "titleUk",
				Type:        schema.FieldTypeText,
				Required:    true,
				Presentable: false,
				Unique:      false,
				Options: schema.TextOptions{
					Min:     types.Pointer(1),
					Max:     types.Pointer(64),
					Pattern: "",
				},
			},

			// DescriptionRu
			&schema.SchemaField{
				System:      false,
				Name:        "descriptionRu",
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

			// DescriptionUk
			&schema.SchemaField{
				System:      false,
				Name:        "descriptionUk",
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

			// Outline Config
			&schema.SchemaField{
				System:      false,
				Name:        "outlineConfig",
				Type:        schema.FieldTypeJson,
				Required:    false,
				Presentable: false,
				Unique:      false,
				Options: schema.JsonOptions{
					MaxSize: 65536,
				},
			},

			// Outline title
			&schema.SchemaField{
				System:      false,
				Name:        "outlineTitle",
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

			// Outline domain
			&schema.SchemaField{
				System:      false,
				Name:        "outlineDomain",
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
		),
	}

	collection.Id = em.OutlineConfigCollectionId

	return collection, nil
}

func init() {
	m.Register(func(db dbx.Builder) error {
		daos := daos.New(db)
		collection, err := makeOutlineConfigsCollection(daos)
		if err != nil {
			return err
		}

		return daos.ImportCollections([]*models.Collection{collection}, false, nil)
	}, func(db dbx.Builder) error {
		daos := daos.New(db)
		collection, err := makeOutlineConfigsCollection(daos)
		if err != nil {
			return err
		}

		daos.DeleteCollection(collection)
		return nil
	})
}
