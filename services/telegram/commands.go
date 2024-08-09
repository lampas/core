package telegram

import (
	"lamas/models"
	"lamas/services/timeago"
	"os"
	"time"

	"github.com/pocketbase/pocketbase"

	tele "gopkg.in/telebot.v3"
)

type LocalizedMenu struct {
	Ru *tele.ReplyMarkup
	Uk *tele.ReplyMarkup
}

type LocalizedBtn struct {
	Enabled bool
	Ru      tele.Btn
	Uk      tele.Btn
}

func RegisterCommands(app *pocketbase.PocketBase, bot *tele.Bot) {
	menuUser := &LocalizedMenu{
		Ru: &tele.ReplyMarkup{ResizeKeyboard: true},
		Uk: &tele.ReplyMarkup{ResizeKeyboard: true},
	}

	menuGuest := &LocalizedMenu{
		Ru: &tele.ReplyMarkup{ResizeKeyboard: true},
		Uk: &tele.ReplyMarkup{ResizeKeyboard: true},
	}

	menuAdmin := &LocalizedMenu{
		Ru: &tele.ReplyMarkup{ResizeKeyboard: true},
		Uk: &tele.ReplyMarkup{ResizeKeyboard: true},
	}

	// Buttons
	btnRefresh := &LocalizedBtn{
		Ru: menuUser.Ru.Text("🔄 Обновить"),
		Uk: menuUser.Uk.Text("🔄 Оновити"),
	}

	btnBack := &LocalizedBtn{
		Ru: menuUser.Ru.Text("🏠 Главная"),
		Uk: menuUser.Uk.Text("🏠 Головна"),
	}

	btnAdminPanel := &LocalizedBtn{
		Ru: menuUser.Ru.Text("👑 Админ-панель"),
		Uk: menuUser.Uk.Text("👑 Адмін-панель"),
	}

	// Commands
	rulesCommand := RegisterRulesCommands(app, bot, menuUser, btnBack)
	settingsCommand := RegisterSettingsCommands(app, bot, menuUser, btnBack)
	watchCommand := RegisterWatchCommands(app, bot, menuUser, btnBack, rulesCommand)

	// Guest menu
	menuGuest.Ru.Reply(menuGuest.Ru.Row(rulesCommand.Ru))
	menuGuest.Uk.Reply(menuGuest.Uk.Row(rulesCommand.Uk))

	// User menu
	menuUser.Ru.Reply(
		menuUser.Ru.Row(watchCommand.Ru),
		menuUser.Ru.Row(rulesCommand.Ru),
		menuUser.Ru.Row(settingsCommand.Ru, btnRefresh.Ru),
	)
	menuUser.Uk.Reply(
		menuUser.Uk.Row(watchCommand.Uk),
		menuUser.Uk.Row(rulesCommand.Uk),
		menuUser.Uk.Row(settingsCommand.Uk, btnRefresh.Uk),
	)

	// Admin menu
	menuAdmin.Ru.Reply(
		menuUser.Ru.Row(watchCommand.Ru),
		menuUser.Ru.Row(btnAdminPanel.Ru, rulesCommand.Ru),
		menuUser.Ru.Row(settingsCommand.Ru, btnRefresh.Ru),
	)
	menuAdmin.Uk.Reply(
		menuUser.Uk.Row(watchCommand.Uk),
		menuUser.Uk.Row(btnAdminPanel.Uk, rulesCommand.Uk),
		menuUser.Uk.Row(settingsCommand.Uk, btnRefresh.Uk),
	)

	// Register home commands
	RegisterHomeCommands(app, bot, menuUser, menuGuest, menuAdmin, btnBack, btnRefresh)
	RegisterAdminCommands(app, bot, btnAdminPanel)
}

func RegisterSettingsCommands(app *pocketbase.PocketBase, bot *tele.Bot, menuUser *LocalizedMenu, btnBack *LocalizedBtn) *LocalizedBtn {
	dao := app.Dao()

	menu := LocalizedMenu{
		Ru: &tele.ReplyMarkup{},
		Uk: &tele.ReplyMarkup{},
	}

	btnLang := LocalizedBtn{
		Ru: menu.Ru.Text("🗣️ Сменить язык / UA"),
		Uk: menu.Uk.Text("🗣️ Змінити мову / RU"),
	}

	btnChangeUsername := LocalizedBtn{
		Ru: menu.Ru.Text("⚠️ Сменить никнейм"),
		Uk: menu.Uk.Text("⚠️ Змінити нікнейм"),
	}

	btnMain := LocalizedBtn{
		Enabled: true,
		Ru:      menu.Ru.Text("⚙️ Настройки"),
		Uk:      menu.Uk.Text("⚙️ Налаштування"),
	}

	menu.Ru.Reply(
		menu.Ru.Row(btnLang.Ru),
		menu.Ru.Row(btnChangeUsername.Ru),
		menu.Ru.Row(btnBack.Ru),
	)

	menu.Uk.Reply(
		menu.Uk.Row(btnLang.Uk),
		menu.Uk.Row(btnChangeUsername.Uk),
		menu.Uk.Row(btnBack.Uk),
	)

	/** Settings menu */
	handler := func(c tele.Context, user *models.User, msg string) error {
		menu := pickByLang(user.Lang, menu.Ru, menu.Uk)

		role := pickByLang(user.Lang, "гость", "гість")
		if user.Role == "user" {
			role = pickByLang(user.Lang, "пользователь", "користувач")
		} else if user.Role == "admin" {
			role = pickByLang(user.Lang, "администратор", "адміністратор")
		}

		renderParams := RenderParams{
			"Settings": map[string]interface{}{
				"CreatedAt": timeago.Format(user.Lang, user.Created.Time()),
				"Role":      role,
			},
		}

		if msg != "" {
			if err := RenderTemplateAndSend(c, dao, user, "settings", renderParams); err != nil {
				return err
			}

			return c.Send(msg, menu, tele.ModeHTML, tele.Silent)
		}

		return RenderTemplateAndSend(c, dao, user, "settings", renderParams, menu)
	}

	handlerUserInfo := func(c tele.Context, user *models.User) error {
		return handler(c, user, "")
	}

	handlerChangeLang := func(c tele.Context, user *models.User) error {
		models.SetUserLang(dao, user, pickByLang(user.Lang, "uk", "ru"))
		msg := pickByLang(user.Lang, "🗣️ Успешно! Язык внешних сервисов также был изменен.", "🗣️ Успішно! Мова зовнішніх сервісів також була змінена.")
		return handler(c, user, msg)
	}

	handlerChangeUsername := func(c tele.Context, user *models.User) error {
		if err := RenderTemplateAndSend(c, dao, user, "settings:username", RenderParams{}); err != nil {
			return err
		}

		if user.Role == "guest" {
			msg := pickByLang(
				user.Lang,
				"❌ Смена имени пользователя станет доступна после утверждения членства.",
				"❌ Зміна імені користувача стане доступною після затвердження членства.",
			)
			return c.Send(msg, tele.ModeHTML, tele.Silent)
		}

		cooldown := user.UsernameChangedAt.Time().Add(models.UsernameChangeCooldown)
		if cooldown.After(time.Now()) {
			msg := pickByLang(
				user.Lang,
				"⌛ До снятия ограничения на повторную смену имени пользователя осталась: "+timeago.Format(user.Lang, cooldown),
				"⌛ До зняття обмеження на повторну зміну імені користувача залишилось: "+timeago.Format(user.Lang, cooldown),
			)
			return c.Send(msg, tele.ModeHTML, tele.Silent)
		}

		return nil
	}

	BotHandlerWithAuth(dao, bot, &btnMain.Ru, handlerUserInfo)
	BotHandlerWithAuth(dao, bot, &btnMain.Uk, handlerUserInfo)

	BotHandlerWithAuth(dao, bot, &btnLang.Ru, handlerChangeLang)
	BotHandlerWithAuth(dao, bot, &btnLang.Uk, handlerChangeLang)

	BotHandlerWithAuth(dao, bot, &btnChangeUsername.Ru, handlerChangeUsername)
	BotHandlerWithAuth(dao, bot, &btnChangeUsername.Uk, handlerChangeUsername)

	return &btnMain
}

func RegisterAdminCommands(app *pocketbase.PocketBase, bot *tele.Bot, btnAdminPanel *LocalizedBtn) {
	dao := app.Dao()
	adminPanelBaseUrl := "https://" + os.Getenv("BASE_DOMAIN") + "/_/?#/collections?collectionId="

	menuInlineAdmin := LocalizedMenu{
		Ru: &tele.ReplyMarkup{},
		Uk: &tele.ReplyMarkup{},
	}

	link := adminPanelBaseUrl + models.UserCollectionId
	btnUsers := LocalizedBtn{
		Ru: menuInlineAdmin.Ru.URL("👥 Управление пользователями", link),
		Uk: menuInlineAdmin.Uk.URL("👥 Керування користувачами", link),
	}

	link = adminPanelBaseUrl + models.OutlineConfigCollectionId
	btnOutline := LocalizedBtn{
		Ru: menuInlineAdmin.Ru.URL("🛡️ Конфигурация Outline VPN", link),
		Uk: menuInlineAdmin.Uk.URL("🛡️ Конфігурація Outline VPN", link),
	}

	link = adminPanelBaseUrl + models.LampacConfigCollectionId + "&recordId=" + models.LampacConfigSingleRecordId
	btnLampac := LocalizedBtn{
		Ru: menuInlineAdmin.Ru.URL("🍿 Конфигурация Lampac", link),
		Uk: menuInlineAdmin.Uk.URL("🍿 Конфігурація Lampac", link),
	}

	link = adminPanelBaseUrl + models.ConfigurationCollectionId + "&recordId=" + models.ConfigurationSingleRecordId
	btnOther := LocalizedBtn{
		Ru: menuInlineAdmin.Ru.URL("🤖 Остальные настройки", link),
		Uk: menuInlineAdmin.Uk.URL("🤖 Інші налаштування", link),
	}

	menuInlineAdmin.Ru.Inline(
		menuInlineAdmin.Ru.Row(btnUsers.Ru),
		menuInlineAdmin.Ru.Row(btnOutline.Ru),
		menuInlineAdmin.Ru.Row(btnLampac.Ru),
		menuInlineAdmin.Ru.Row(btnOther.Ru),
	)

	menuInlineAdmin.Uk.Inline(
		menuInlineAdmin.Uk.Row(btnUsers.Uk),
		menuInlineAdmin.Uk.Row(btnOutline.Uk),
		menuInlineAdmin.Uk.Row(btnLampac.Uk),
		menuInlineAdmin.Uk.Row(btnOther.Uk),
	)

	handler := func(c tele.Context, user *models.User) error {
		if user.Role != "admin" {
			return nil
		}

		err := RenderTemplateAndSend(
			c, dao, user, "start:admin",
			RenderParams{
				"Analytics": models.GetUsersAnalytics(dao),
			},
		)

		if err != nil {
			return err
		}

		return c.Send(
			pickByLang(user.Lang, "<strong>Админ-панель:</strong>", "<strong>Адмін-панель:</strong>"),
			pickByLang(user.Lang, menuInlineAdmin.Ru, menuInlineAdmin.Uk),
			tele.ModeHTML,
		)
	}

	BotHandlerWithAuth(dao, bot, &btnAdminPanel.Ru, handler)
	BotHandlerWithAuth(dao, bot, &btnAdminPanel.Uk, handler)
}

func RegisterHomeCommands(app *pocketbase.PocketBase, bot *tele.Bot, menuUser *LocalizedMenu, menuGuest *LocalizedMenu, menuAdmin *LocalizedMenu, btnBack *LocalizedBtn, btnRefresh *LocalizedBtn) {
	dao := app.Dao()

	handler := func(c tele.Context, user *models.User) error {

		// Guest
		if user.Role == "guest" {
			return RenderTemplateAndSend(
				c, dao, user, "start:guest",
				RenderParams{},
				pickByLang(user.Lang, menuGuest.Ru, menuGuest.Uk),
			)
		}

		menu := pickByLang(user.Lang, menuUser.Ru, menuUser.Uk)
		if user.Role == "admin" {
			menu = pickByLang(user.Lang, menuAdmin.Ru, menuAdmin.Uk)
		}

		// User
		return RenderTemplateAndSend(
			c, dao, user, "start:user",
			RenderParams{
				"WatchEnabled":    models.GetLampacConfig(dao).Enabled,
				"OutlineCommands": models.GetOutlineCommands(dao, user),
			},
			menu,
		)
	}

	// Watch
	BotHandlerWithAuth(dao, bot, "/start", handler)
	BotHandlerWithAuth(dao, bot, &btnBack.Ru, handler)
	BotHandlerWithAuth(dao, bot, &btnBack.Uk, handler)
	BotHandlerWithAuth(dao, bot, &btnRefresh.Ru, handler)
	BotHandlerWithAuth(dao, bot, &btnRefresh.Uk, handler)
}

func RegisterRulesCommands(app *pocketbase.PocketBase, bot *tele.Bot, menuUser *LocalizedMenu, btnBack *LocalizedBtn) *LocalizedBtn {
	dao := app.Dao()

	menu := LocalizedMenu{
		Ru: &tele.ReplyMarkup{ResizeKeyboard: true},
		Uk: &tele.ReplyMarkup{ResizeKeyboard: true},
	}

	btnMain := LocalizedBtn{
		Enabled: true,
		Ru:      menu.Uk.Text("📜 Правила использования"),
		Uk:      menu.Uk.Text("📜 Правила користування"),
	}

	menu.Ru.Reply(menu.Ru.Row(btnBack.Ru))
	menu.Uk.Reply(menu.Uk.Row(btnBack.Uk))

	handler := func(c tele.Context, user *models.User) error {
		menu := pickByLang(user.Lang, menu.Ru, menu.Uk)
		return RenderTemplateAndSend(c, dao, user, "rules", RenderParams{}, menu)
	}

	BotHandlerWithAuth(dao, bot, &btnMain.Ru, handler)
	BotHandlerWithAuth(dao, bot, &btnMain.Uk, handler)

	return &btnMain
}

func RegisterWatchCommands(app *pocketbase.PocketBase, bot *tele.Bot, menuUser *LocalizedMenu, btnBack *LocalizedBtn, btnRules *LocalizedBtn) *LocalizedBtn {
	dao := app.Dao()
	config := models.GetLampacConfig(dao)

	if !config.Enabled {
		app.Logger().Info("[Telegram] Lampac service is disabled, watch commands are not registered")
		return &LocalizedBtn{Enabled: false}
	}

	menu := LocalizedMenu{
		Ru: &tele.ReplyMarkup{ResizeKeyboard: true},
		Uk: &tele.ReplyMarkup{ResizeKeyboard: true},
	}

	menu.Ru.Reply(menu.Ru.Row(btnRules.Ru, btnBack.Ru))
	menu.Uk.Reply(menu.Uk.Row(btnRules.Uk, btnBack.Uk))

	btnMain := LocalizedBtn{
		Enabled: true,
		Ru:      menu.Uk.Text("🍿 Кино"),
		Uk:      menu.Uk.Text("🍿 Кіно"),
	}

	// Menu callback
	handler := func(c tele.Context, user *models.User) error {
		if user.Role == "guest" {
			return nil
		}
		domain := user.Username + ".tv." + os.Getenv("BASE_DOMAIN")

		// Menu
		menu := pickByLang(user.Lang, menu.Ru, menu.Uk)

		// Send
		renderParams := RenderParams{
			"Watch": map[string]string{
				"Domain": domain,
			},
		}

		return RenderTemplateAndSend(c, dao, user, "watch", renderParams, menu)
	}

	// Watch
	BotHandlerWithAuth(dao, bot, &btnMain.Ru, handler)
	BotHandlerWithAuth(dao, bot, &btnMain.Uk, handler)

	return &btnMain
}
