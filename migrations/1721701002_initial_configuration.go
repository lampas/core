package migrations

import (
	"crypto/rand"
	"fmt"
	"lamas/models"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		configuration := &models.Configuration{}
		configuration.Id = models.ConfigurationSingleRecordId

		configuration.BotToken = os.Getenv("INITIALIZE_TELEGRAM_TOKEN")
		configuration.BotUsername = ""
		configuration.BotNameRu = "LampaS"
		configuration.BotNameUk = "LampaS"
		configuration.BotDescriptionRu = "🐳 Закрытый клуб свободного интернета"
		configuration.BotDescriptionUk = "🐳 Закритий клуб вільного інтернету"
		configuration.BotShortDescriptionRu = "Добро пожаловать в закрытый клуб по интересам. Хочешь присоединиться? Попроси друзей отправить тебе приглашение 🐳"
		configuration.BotShortDescriptionUk = "Ласкаво просимо до закритого клубу за інтересами. Хочеш приєднатися? Попроси друзів надіслати тобі запрошення 🐳"

		configuration.AppDefaultLang = "ru"
		configuration.AppTitleRu = "Клуб свободного интернете"
		configuration.AppTitleUk = "Клуб вільного інтернету"
		configuration.AppDescriptionRu = "🐳 LampaS — это OpenSource проект, цель которого облегчить доступ к различным self-hosted сервисам. Посмотреть исходный код можно на <a href=\"https://github.com/docker-pet\" target=\"_blank\">GitHub</a>."
		configuration.AppDescriptionUk = "🐳 LampaS — це OpenSource проект, мета якого полегшити доступ до різних self-hosted сервісів. Переглянути вихідний код можна на <a href=\"https://github.com/docker-pet\" target=\"_blank\">GitHub</a>."

		// Generate encryption key
		CookieEncryptionKey := make([]byte, 16)
		rand.Read(CookieEncryptionKey)
		configuration.CookieEncryptionKey = fmt.Sprintf("%x", CookieEncryptionKey)

		return dao.Save(configuration)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		configuration := models.GetConfiguration(dao)
		if configuration.HasId() {
			return dao.Delete(configuration)
		}

		return nil
	})
}
