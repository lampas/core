package models

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

func TemplateQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Template{})
}

func GetAllTemplates(dao *daos.Dao) ([]*Template, error) {
	items := []*Template{}

	err := TemplateQuery(dao).
		All(&items)

	if err != nil {
		return nil, errors.New("templates not found")
	}

	return items, nil
}

func GetTemplate(dao *daos.Dao, slug string) (*Template, error) {
	item := &Template{}

	err := TemplateQuery(dao).
		AndWhere(dbx.HashExp{"slug": slug}).
		One(item)

	if err != nil {
		return nil, errors.New("Template \"" + slug + "\" not found")
	}

	return item, nil
}

func FormatBannerUrl(template *Template) string {
	if template.Banner == "" {
		return ""
	}

	return fmt.Sprintf(
		"/api/files/%s/%s/%s",
		TemplateCollectionId,
		template.Id,
		template.Banner,
	)
}

func FormatBannerPath(daos *daos.Dao, template *Template) string {
	if template.Banner == "" {
		return ""
	}

	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.Join(pwd, "pb_data", "storage", TemplateCollectionId, template.Id, template.Banner)
}
