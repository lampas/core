package templates

import (
	"bytes"
	"errors"
	"html/template"
	"lamas/models"
	"log"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/store"
)

type Renderer struct {
	Template   *template.Template
	BannerUrl  string
	BannerPath string
	ParseError error
}

var fallbackTemplate, _ = template.New("").Parse("Template not found")
var fallbackRenderer = &Renderer{
	Template:   fallbackTemplate,
	BannerUrl:  "",
	BannerPath: "",
	ParseError: errors.New("template not found"),
}

var templateStory = store.New[*Renderer](nil)

func RegisterService(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		return fillTemplateStory(app)
	})

	app.OnModelAfterCreate(models.TemplateCollectionId).Add(func(e *core.ModelEvent) error {
		return fillTemplateStory(app)
	})

	app.OnModelAfterDelete(models.TemplateCollectionId).Add(func(e *core.ModelEvent) error {
		return fillTemplateStory(app)
	})

	app.OnModelAfterUpdate(models.TemplateCollectionId).Add(func(e *core.ModelEvent) error {
		return fillTemplateStory(app)
	})
}

func fillTemplateStory(app *pocketbase.PocketBase) error {
	dao := app.Dao()
	templates, err := models.GetAllTemplates(dao)
	if err != nil {
		log.Println("ERROR: Unable to get templates", err)
		return nil
	}

	data := make(map[string]*Renderer)
	for _, item := range templates {
		for _, lang := range []string{"ru", "uk"} {
			// Content
			content := ""
			if lang == "uk" || (lang == "ru" && item.ContentRu == "") {
				content = item.ContentUk
			}

			if lang == "ru" || (lang == "uk" && item.ContentUk == "") {
				content = item.ContentRu
			}

			// Format template
			tpl, err := template.New("").Parse(content)
			data[lang+"@"+item.Slug] = &Renderer{
				Template:   tpl,
				BannerUrl:  models.FormatBannerUrl(item),
				BannerPath: models.FormatBannerPath(dao, item),
				ParseError: err,
			}
		}
	}

	templateStory.Reset(data)
	return nil
}

func GetRenderer(slug string, lang string) *Renderer {
	fallbackLang := "ru"
	if lang == "ru" {
		fallbackLang = "uk"
	}

	if renderer := templateStory.Get(lang + "@" + slug); renderer != nil {
		return renderer
	}

	// Fallback
	if renderer := templateStory.Get(fallbackLang + "@" + slug); renderer != nil {
		return renderer
	}

	return fallbackRenderer
}

func RenderTemplate(slug string, lang string, data any) (string, string, string) {
	renderer := GetRenderer(slug, lang)

	if renderer.ParseError != nil {
		return "ERROR: unable to parse template", renderer.BannerUrl, renderer.BannerPath
	}

	if renderer.Template == nil {
		return "ERROR: invalid or nil template", renderer.BannerUrl, renderer.BannerPath
	}

	buf := new(bytes.Buffer)

	if err := renderer.Template.Execute(buf, data); err != nil {
		return "ERROR: " + err.Error(), renderer.BannerUrl, renderer.BannerPath
	}

	return buf.String(), renderer.BannerUrl, renderer.BannerPath
}
