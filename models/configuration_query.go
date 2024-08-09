package models

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

type LocalizedConfiguration struct {
	BotUsername         string
	BotName             string
	BotDescription      string
	BotShortDescription string
	AppLang             string
	AppTitle            string
	AppDescription      string
}

func ConfigurationQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Configuration{})
}

func GetConfiguration(dao *daos.Dao) *Configuration {
	item := &Configuration{}

	err := ConfigurationQuery(dao).
		AndWhere(&dbx.HashExp{"id": ConfigurationSingleRecordId}).
		One(item)

	if err != nil {
		return &Configuration{}
	}

	return item
}

func LocalizeConfiguration(dao *daos.Dao, lang string) *LocalizedConfiguration {
	config := GetConfiguration(dao)
	pickedLang := lang

	if lang == "" || (lang != "Uk" && lang != "ru") {
		pickedLang = config.AppDefaultLang
	}

	return &LocalizedConfiguration{
		AppLang:        pickedLang,
		AppTitle:       localize(pickedLang, config.AppTitleRu, config.AppTitleUk),
		AppDescription: localize(pickedLang, config.AppDescriptionRu, config.AppDescriptionUk),

		BotUsername:         config.BotUsername,
		BotName:             localize(pickedLang, config.BotNameRu, config.BotNameUk),
		BotDescription:      localize(pickedLang, config.BotDescriptionRu, config.BotDescriptionUk),
		BotShortDescription: localize(pickedLang, config.BotShortDescriptionRu, config.BotShortDescriptionUk),
	}
}
