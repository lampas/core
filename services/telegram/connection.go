package telegram

import (
	"lamas/models"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"

	tele "gopkg.in/telebot.v3"
)

type StopCallback func() error

func ConnectTelegram(app *pocketbase.PocketBase) StopCallback {
	dao := app.Dao()
	config := models.GetConfiguration(dao)
	configRu := models.LocalizeConfiguration(dao, "ru")
	configUk := models.LocalizeConfiguration(dao, "uk")

	// Create a new bot
	bot, err := tele.NewBot(
		tele.Settings{
			Token: config.BotToken,
			Poller: &tele.LongPoller{
				Timeout:        10 * time.Second,
				AllowedUpdates: []string{"message", "callback_query"},
			},
		},
	)

	// Stop callback
	stop := func() error {
		if err != nil {
			app.Logger().Info("[Telegram] Bot already stopped")
			return nil
		}

		app.Logger().Info("[Telegram] Stopping bot connection")
		bot.Stop()
		return nil
	}

	// Connection error
	if err != nil {
		app.Logger().Error(
			"[Telegram] Field to create a new bot connection",
			"error", err,
		)

		return stop
	}

	// Connection success
	app.Logger().Info(
		"[Telegram] Bot connection established",
		"Username", bot.Me.Username,
	)

	// Sync bot description
	if SyncBotDescription(bot, configRu, config.AppDefaultLang == "ru") {
		app.Logger().Info("[Telegram] Bot description updated")
	}

	if SyncBotDescription(bot, configUk, config.AppDefaultLang == "uk") {
		app.Logger().Info("[Telegram] Bot description updated")
	}

	// Need to update username
	if !strings.EqualFold(bot.Me.Username, config.BotUsername) {
		config.BotUsername = bot.Me.Username
		dao.Save(config)

		app.Logger().Info(
			"[Telegram] Bot username updated",
			"Old", config.BotUsername,
			"New", bot.Me.Username,
		)

		return stop
	}

	// Commands
	RegisterCommands(app, bot)

	// Start bot
	go bot.Start()
	return stop
}

func SyncBotDescription(bot *tele.Bot, config *models.LocalizedConfiguration, isPrimaryLang bool) bool {
	updated := false
	lang := config.AppLang

	// A trick to update the l10n fallback
	if isPrimaryLang {
		updated = SyncBotDescription(bot, config, false)
		lang = ""
	}

	currentDescription, _ := bot.MyDescription(lang)
	currentShortDescription, _ := bot.MyShortDescription(lang)
	currentName, _ := bot.MyName(lang)

	if currentDescription.Description != config.BotDescription {
		updated = true
		bot.SetMyDescription(config.BotDescription, lang)
	}

	if currentShortDescription.ShortDescription != config.BotShortDescription {
		updated = true
		bot.SetMyShortDescription(config.BotShortDescription, lang)
	}

	if currentName.Name != config.BotName {
		updated = true
		bot.SetMyName(config.BotName, lang)
	}

	return updated
}
