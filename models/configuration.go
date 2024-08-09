package models

import (
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*Configuration)(nil)

type Configuration struct {
	models.BaseModel

	BotToken    string `db:"botToken" json:"botToken"`
	BotUsername string `db:"botUsername" json:"botUsername"`

	BotNameRu             string `db:"botNameRu" json:"botNameRu"`
	BotNameUk             string `db:"botNameUk" json:"botNameUk"`
	BotDescriptionRu      string `db:"botDescriptionRu" json:"botDescriptionRu"`
	BotDescriptionUk      string `db:"botDescriptionUk" json:"botDescriptionUk"`
	BotShortDescriptionRu string `db:"botShortDescriptionRu" json:"botShortDescriptionRu"`
	BotShortDescriptionUk string `db:"botShortDescriptionUk" json:"botShortDescriptionUk"`

	AppDefaultLang      string `db:"appDefaultLang" json:"appDefaultLang"`
	AppTitleRu          string `db:"appTitleRu" json:"appTitleRu"`
	AppTitleUk          string `db:"appTitleUk" json:"appTitleUk"`
	AppDescriptionRu    string `db:"appDescriptionRu" json:"appDescriptionRu"`
	AppDescriptionUk    string `db:"appDescriptionUk" json:"appDescriptionUk"`
	CookieEncryptionKey string `db:"cookieEncryptionKey" json:"cookieEncryptionKey"`
}

var (
	ConfigurationTableName      = "configuration"
	ConfigurationCollectionId   = "zrwlsj67megi0f9"
	ConfigurationSingleRecordId = "r6mng99ey51dkj5"
)

func (m *Configuration) TableName() string {
	return ConfigurationTableName
}
