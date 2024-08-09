package migrations

import (
	"lamas/models"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	guestStartTemplate := &models.Template{
		Slug:      "start:guest",
		ContentRu: "Привет! Это закрытое сообщество.\nЕсли хотите присоединиться, попросите другого участника подтвердить членство.\n Для этого друг должен отправить боту команду: <code>/invite {{.User.Username}}</code>",
		ContentUk: "Привіт! Це закрита спільнота.\nЯкщо хочете приєднатися, попросіть іншого учасника підтвердити членство.\nДля цього другу потрібно відправити боту команду: <code>/invite {{.User.Username}}</code>",
	}

	userStartTemplate := &models.Template{
		Slug:      "start:user",
		ContentRu: "С возвращением!",
		ContentUk: "",
	}

	adminStartTemplate := &models.Template{
		Slug:      "start:admin",
		ContentRu: "<strong>Администраторов:</strong> {{.Analytics.Admins}}\n<strong>Пользователей:</strong> {{.Analytics.Users}}\n<strong>Гостей:</strong> {{.Analytics.Guests}}",
		ContentUk: "<strong>Адміністраторів:</strong> {{.Analytics.Admins}}\n<strong>Користувачів:</strong> {{.Analytics.Users}}\n<strong>Гостей:</strong> {{.Analytics.Guests}}",
	}

	rulesTemplate := &models.Template{
		Slug:      "rules",
		ContentRu: "<strong>Правила использования:</strong>\nПользуясь ботом, вы соглашаетесь с правилами:\n\n1. <strong>Запрет на совместное использование.</strong> Персональный доступ к сервисам нельзя передавать другим. При злоупотреблении мы ограничим доступ вам, приглашенным вами и приведшему вас пользователю.\n\n2. <strong>Запрет на злоупотребление.</strong> Нельзя использовать сервисы для атак или нарушения законов стран вашего гражданства, проживания или местоположения VPN.\n\n3. <strong>Отказ от ответственности.</strong> Мы агрегируем публичные интернет-ресурсы и проксируем трафик, не неся ответственности за контент. Стабильность работы и корректность данных не гарантируются.\n\n4. <strong>Хранение данных.</strong> Мы храним минимум данных: имя пользователя, идентификатор Telegram, IP-адрес и общую информацию о потреблении трафика. Мы не злоупотребляем сбором данных, но не исключаем попадание дополнительных данных в логи.\n\n5. <strong>Взаимоуважение.</strong> Правила не являются юридическим документом, но мы ожидаем взаимного уважения. При нарушении правил доступ может быть ограничен.",
		ContentUk: "<strong>Правила використання:</strong>\nКористуючись ботом, ви погоджуєтесь з правилами:\n\n1. <strong>Заборона на спільне використання.</strong> Персональний доступ до сервісів не можна передавати іншим. У разі зловживання ми обмежимо доступ вам, запрошеним вами та користувачу, який вас запросив.\n\n2. <strong>Заборона на зловживання.</strong> Заборонено використовувати сервіси для атак або порушення законів країн вашого громадянства, проживання чи місцезнаходження VPN.\n\n3. <strong>Відмова від відповідальності.</strong> Ми агрегуємо публічні інтернет-ресурси та проксіруємо трафік, не несячи відповідальності за контент. Стабільність роботи та коректність даних не гарантуються.\n\n4. <strong>Зберігання даних.</strong> Ми зберігаємо мінімум даних: ім'я користувача, ідентифікатор Telegram, IP-адресу та загальну інформацію про споживання трафіку. Ми не зловживаємо збором даних, але не виключаємо потрапляння додаткових даних у логи.\n\n5. <strong>Взаємоповага.</strong> Правила не є юридичним документом, але ми очікуємо взаємної поваги. У разі порушення правил доступ може бути обмежений.",
	}

	settingsTemplate := &models.Template{
		Slug:      "settings",
		ContentRu: "<strong>Имя пользователя:</strong> {{.User.Username}}\n<strong>Регистрация:</strong> {{.Settings.CreatedAt}}\n<strong>Язык:</strong> русский\n<strong>Статус:</strong> {{.Settings.Role}}",
		ContentUk: "<strong>Ім'я користувача:</strong> {{.User.Username}}\n<strong>Реєстрація:</strong> {{.Settings.CreatedAt}}\n<strong>Мова:</strong> українська\n<strong>Статус:</strong> {{.Settings.Role}}",
	}

	settingsUsernameTemplate := &models.Template{
		Slug:      "settings:username",
		ContentRu: "<strong>Изменение имени пользователя</strong>\n\nИзменение имени пользователя приведет к безвозвратной блокировке всех активных ключей достпупа и замене их на новые. ⚠️ <strong>Рекомендуем использовать этот метод только в случае подозрения на компрометацию ваших данных.</strong>\n\nНовое имя пользователя будет сгенерировано автоматически.",
		ContentUk: "<strong>Зміна імені користувача</strong>\n\nЗміна імені користувача призведе до безповоротного блокування всіх активних ключів доступу та їх заміни на нові. ⚠️ <strong>Рекомендуємо використовувати цей метод лише в разі підозри на компрометацію ваших даних.</strong>\n\nНове ім'я користувача буде згенеровано автоматично.",
	}

	vpnDefaultTemplate := &models.Template{
		Slug:      "vpn:default",
		ContentRu: "VPN доступы",
		ContentUk: "VPN доступы",
	}

	vpnHelpTemplate := &models.Template{
		Slug:      "vpn:help",
		ContentRu: "Настройки VPN",
		ContentUk: "Настройки VPN",
	}

	watchTemplate := &models.Template{
		Slug:      "watch",
		ContentRu: "Не знаешь, что посмотреть этим вечером? Удобный каталог фильмов, сериалов и аниме поможет тебе сделать выбор!\n\n🍿 <strong>Ваш персональный адрес:</strong> <a href=\"https://{{.Watch.Domain}}\">{{.Watch.Domain}}</a>",
		ContentUk: "Не знаєш, що подивитися сьогодні ввечері? Зручний каталог фільмів, серіалів та аніме допоможе тобі зробити вибір!\n\n🍿 <strong>Ваша персональна адреса:</strong> <a href=\"https://{{.Watch.Domain}}\">{{.Watch.Domain}}</a>",
	}

	items := []*models.Template{
		guestStartTemplate,
		userStartTemplate,
		adminStartTemplate,
		rulesTemplate,
		settingsTemplate,
		settingsUsernameTemplate,
		vpnDefaultTemplate,
		vpnHelpTemplate,
		watchTemplate,
	}

	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		for _, item := range items {
			if err := dao.Save(item); err != nil {
				return err
			}
		}

		return nil
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		for _, item := range items {
			if item, err := models.GetTemplate(dao, item.Slug); err == nil {
				if err := dao.Delete(item); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
