package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*OutlineToken)(nil)

type OutlineToken struct {
	models.BaseModel

	User             string        `db:"user" json:"user"`
	OutlineConfig    string        `db:"outlineConfig" json:"outlineConfig"`
	OutlineKey       types.JsonRaw `db:"outlineKey" json:"outlineKey"`
	BytesTransferred int64         `db:"bytesTransferred" json:"bytesTransferred"`
}

var (
	OutlineTokenTableName    = "outline_tokens"
	OutlineTokenCollectionId = "j39lu9hd47djf7l"
)

func (m *OutlineToken) TableName() string {
	return OutlineTokenTableName
}
