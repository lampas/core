package telegram

import (
	"lamas/models"
	"lamas/services/templates"
	"log"

	tele "gopkg.in/telebot.v3"

	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/tools/store"
)

type TelegramHandler func(c tele.Context, user *models.User) error
type RenderParams map[string]interface{}

var PhotoStory = store.New[*tele.Photo](nil)

func BotHandlerWithAuth(dao *daos.Dao, bot *tele.Bot, endpoint interface{}, handler TelegramHandler) {
	bot.Handle(endpoint, func(c tele.Context) error {
		user, err := models.GetUserByTelegramId(dao, c.Sender())
		if err != nil {
			log.Println("WARNING: Unable to get user", err)
			return err
		}

		return handler(c, user)
	})
}

func GetCachedPhoto(path string) *tele.Photo {
	if photo := PhotoStory.Get(path); photo != nil {
		return photo
	}

	photo := &tele.Photo{File: tele.FromDisk(path)}
	PhotoStory.Set(path, photo)

	return photo

}

func RenderTemplate(dao *daos.Dao, user *models.User, templateSlug string, params RenderParams) (string, tele.Album) {
	// Merge params
	mergedParams := RenderParams{
		"User": user,
	}

	for k, v := range params {
		mergedParams[k] = v
	}

	// Render
	template, _, bannerPath := templates.RenderTemplate(templateSlug, user.Lang, mergedParams)

	var banner tele.Album
	if bannerPath != "" {
		banner = tele.Album{GetCachedPhoto(bannerPath)}
	}

	return template, banner
}

func RenderTemplateAndSend(
	c tele.Context,
	dao *daos.Dao,
	user *models.User,
	templateSlug string,
	params RenderParams,
	opts ...interface{},
) error {
	content, banner := RenderTemplate(dao, user, templateSlug, params)
	mergedOpts := append(opts, tele.Silent, tele.NoPreview, tele.ModeHTML)

	if banner != nil {
		banner.SetCaption(content)

		// Send message
		if err := c.SendAlbum(banner, mergedOpts...); err != nil {
			return err
		}

		// Split reply markup
		if err := SendReplyMarkup(c, user, opts...); err != nil {
			return err
		}

		return nil
	}

	return c.Send(content, mergedOpts...)
}

func pickByLang[T any](lang string, ru T, uk T) T {
	if lang == "uk" {
		return uk
	}

	return ru
}

func SendReplyMarkup(c tele.Context, user *models.User, opts ...interface{}) error {
	for _, opt := range opts {
		switch opt.(type) {
		case *tele.ReplyMarkup:
			return c.Send(pickByLang(user.Lang, "ℹ️ Для навигации используйте клавиатурное меню", "ℹ️ Для навігації використовуйте клавіатурне меню"), tele.Silent, opt)
		}
	}

	return nil
}
