package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

func LampacConfigQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&LampacConfig{})
}

func GetLampacConfig(dao *daos.Dao) *LampacConfig {
	item := &LampacConfig{}

	err := LampacConfigQuery(dao).
		AndWhere(&dbx.HashExp{"id": LampacConfigSingleRecordId}).
		One(item)

	if err != nil {
		return &LampacConfig{
			Enabled: false,
		}
	}

	return item
}
