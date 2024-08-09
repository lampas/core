package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*OutlineConfig)(nil)

type OutlineConfig struct {
	models.BaseModel

	Enabled  bool   `db:"enabled" json:"enabled"`
	Slug     string `db:"slug" json:"slug"`
	Template string `db:"template" json:"template"`

	TitleRu       string `db:"titleRu" json:"titleRu"`
	TitleUk       string `db:"titleUk" json:"titleUk"`
	DescriptionRu string `db:"descriptionRu" json:"descriptionRu"`
	DescriptionUk string `db:"descriptionUk" json:"descriptionUk"`

	OutlineConfig types.JsonRaw `db:"outlineConfig" json:"outlineConfig"`
	OutlineTitle  string        `db:"outlineTitle" json:"outlineTitle"`
	OutlineDomain string        `db:"outlineDomain" json:"outlineDomain"`
}

var (
	OutlineConfigTableName    = "outline_configs"
	OutlineConfigCollectionId = "osb1prd163k4nhm"
)

func (m *OutlineConfig) TableName() string {
	return OutlineConfigTableName
}
