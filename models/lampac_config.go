package models

import (
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

var _ models.Model = (*LampacConfig)(nil)

type LampacConfig struct {
	models.BaseModel

	Enabled             bool `db:"enabled" json:"enabled"`
	AdultEnabled        bool `db:"adultEnabled" json:"adultEnabled"`
	JacketEnabled       bool `db:"jacketEnabled" json:"jacketEnabled"`
	TorrServerEnabled   bool `db:"torrServerEnabled" json:"torrServerEnabled"`
	DLNAEnabled         bool `db:"dlnaEnabled" json:"dlnaEnabled"`
	TracksModifyEnabled bool `db:"tracksModifyEnabled" json:"tracksModifyEnabled"`

	CheckOnlineSearch bool `db:"checkOnlineSearch" json:"checkOnlineSearch"`
	MultiAccess       bool `db:"multiAccess" json:"multiAccess"`
	LocalIP           bool `db:"localIP" json:"localIP"`

	ListenScheme      string `db:"listenScheme" json:"listenScheme"`
	ImageProxyService string `db:"imageProxyService" json:"imageProxyService"`
	TorSocksProxy     string `db:"torSocksProxy" json:"torSocksProxy"`

	OverrideInitConfig     types.JsonRaw        `db:"overrideInitConfig" json:"overrideInitConfig"`
	OverrideModuleManifest types.JsonArray[any] `db:"overrideModuleManifest" json:"overrideModuleManifest"`
}

var (
	LampacConfigTableName      = "lampac_config"
	LampacConfigCollectionId   = "5qil9t8fqb58hg0"
	LampacConfigSingleRecordId = "r6mng99ey51dkj5"
)

func (m *LampacConfig) TableName() string {
	return LampacConfigTableName
}
