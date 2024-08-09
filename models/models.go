package models

import (
	"errors"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

/**
 * Prevent creation/deletion of single record models
 */
func RegisterModels(app *pocketbase.PocketBase) {
	err := errors.New("you cannot create or delete records in the current table")

	app.OnModelBeforeCreate(ConfigurationCollectionId, LampacConfigCollectionId).Add(func(e *core.ModelEvent) error {
		switch e.Model.TableName() {
		case ConfigurationTableName:
			if e.Model.GetId() != ConfigurationSingleRecordId {
				return err
			}

		case LampacConfigTableName:
			if e.Model.GetId() != LampacConfigSingleRecordId {
				return err
			}
		}

		return nil
	})

	app.OnModelBeforeDelete(ConfigurationCollectionId, LampacConfigCollectionId).Add(func(e *core.ModelEvent) error {
		return err
	})
}
