package models

import (
	"github.com/pocketbase/pocketbase/models"
)

var _ models.Model = (*Template)(nil)

type Template struct {
	models.BaseModel

	Slug      string `db:"slug" json:"slug"`
	Banner    string `db:"banner" json:"banner"`
	ContentRu string `db:"contentRu" json:"contentRu"`
	ContentUk string `db:"contentUk" json:"contentUk"`
}

var (
	TemplateTableName    = "templates"
	TemplateCollectionId = "elbfv78c4rbph8v"
)

func (m *Template) TableName() string {
	return TemplateTableName
}
