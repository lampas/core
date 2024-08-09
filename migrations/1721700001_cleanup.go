package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		// Find users collection
		collection, err := daos.New(db).FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		// Validate collection type
		if collection.Type == models.CollectionTypeAuth {
			return daos.New(db).DeleteCollection(collection)
		}

		return nil
	}, func(db dbx.Builder) error {
		return nil
	})
}
