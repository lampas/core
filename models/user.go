package models

import (
	"time"

	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*User)(nil)

type User struct {
	models.BaseModel

	Username          string         `db:"username" json:"username"`
	ExternalId        string         `db:"externalId" json:"externalId"`
	Name              string         `db:"name" json:"name"`
	TelegramUsername  string         `db:"telegramUsername" json:"telegramUsername"`
	Role              string         `db:"role" json:"role"`
	Lang              string         `db:"lang" json:"lang"`
	UserSetLang       bool           `db:"userSetLang" json:"userSetLang"`
	UsernameChangedAt types.DateTime `db:"usernameChangedAt" json:"usernameChangedAt"`
}

var (
	UserTableName          = "users"
	UserCollectionId       = "ta1zhasz244qsrx"
	UsernameChangeCooldown = time.Hour * 24 * 7
)

func (m *User) TableName() string {
	return UserTableName
}
