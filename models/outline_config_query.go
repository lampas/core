package models

import (
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

type OutlineCommand struct {
	Command string
	Title   string
}

func OutlineConfigQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&OutlineConfig{})
}

func GetOutlineConfigs(dao *daos.Dao) (*[]OutlineConfig, error) {
	items := &[]OutlineConfig{}

	err := OutlineConfigQuery(dao).
		AndWhere(dbx.HashExp{"enabled": true}).
		All(items)

	if err != nil {
		return nil, errors.New("config records not found")
	}

	return items, nil
}

func GetOutlineCommands(dao *daos.Dao, user *User) *[]OutlineCommand {
	items, err := GetOutlineConfigs(dao)

	if err != nil {
		return &[]OutlineCommand{}
	}

	commands := []OutlineCommand{}
	for _, item := range *items {
		title := item.TitleRu
		if item.TitleRu == "" || (user.Lang == "uk" && item.TitleUk != "") {
			title = item.TitleUk
		}

		commands = append(commands, OutlineCommand{
			Command: item.Slug + "_vpn",
			Title:   title,
		})
	}

	return &commands
}
