package models

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"

	ung "github.com/dillonstreator/go-unique-name-generator"
	"github.com/dillonstreator/go-unique-name-generator/dictionaries"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	tele "gopkg.in/telebot.v3"
)

const (
	AuthProviderTelegram = "telegram"
)

var usernameGenerator = func() *ung.UniqueNameGenerator {
	dictionary := append(dictionaries.Adjectives, dictionaries.Animals...)
	var filteredDictionary []string

	for _, str := range dictionary {
		if len(str) <= 10 {
			filteredDictionary = append(filteredDictionary, str)
		}
	}

	return ung.NewUniqueNameGenerator(
		ung.WithSeparator("-"),
		ung.WithStyle(ung.Lower),
		ung.WithDictionaries([][]string{
			filteredDictionary,
		}),
	)
}()

func UserQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&User{})
}

func GetUserById(dao *daos.Dao, id string) (*User, error) {
	item := &User{}
	err := UserQuery(dao).
		AndWhere(dbx.HashExp{"id": id}).
		One(item)

	if err != nil {
		return nil, errors.New("User record not found")
	}

	return item, nil
}

func GetUserByUsername(dao *daos.Dao, username string) (*User, error) {
	item := &User{}
	err := UserQuery(dao).
		AndWhere(dbx.HashExp{"username": username}).
		One(item)

	if err != nil {
		return nil, errors.New("User record not found")
	}

	return item, nil
}

func GetUserByExternalId(dao *daos.Dao, externalId string, provider string) (*User, error) {
	if provider != AuthProviderTelegram {
		return nil, errors.New("invalid provider")
	}

	item := &User{}
	err := UserQuery(dao).
		AndWhere(dbx.HashExp{"externalId": externalId + "@" + provider}).
		One(item)

	if err != nil {
		return nil, errors.New("User record not found")
	}

	return item, nil
}

func CreateUser(dao *daos.Dao, user *User, externalId string, provider string) (*User, error) {
	// Generate and validate username (while it's unique)
	postfixMaxInt := int64(100)
	attempt := 2
	for {
		username := usernameGenerator.Generate()

		// Add random number to username
		if attempt > 0 {
			username = fmt.Sprintf("%s-%0*d", username, attempt, rand.Int63n(postfixMaxInt/10))
		}

		// Validate username
		_, err := GetUserByUsername(dao, username)
		if err != nil {
			user.Username = username
			break
		}

		postfixMaxInt *= 10
		attempt++
	}

	// Create user
	user.ExternalId = externalId + "@" + provider
	err := dao.Save(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpdateUser(dao *daos.Dao, user *User) bool {
	err := dao.Save(user)
	if err != nil {
		log.Println("WARNING: Unable to update user", err)
		return false
	}

	return true
}

func GetUserByTelegramId(dao *daos.Dao, sender *tele.User) (*User, error) {
	externalId := strconv.FormatInt(sender.ID, 10)
	user, err := GetUserByExternalId(dao, externalId, AuthProviderTelegram)
	name := sender.FirstName + " " + sender.LastName

	// User language
	lang := "ru"
	if sender.LanguageCode == "uk" {
		lang = "uk"
	}

	// Create user if not exists
	if err != nil {
		user, err = CreateUser(dao, &User{
			TelegramUsername: sender.Username,
			Name:             name,
			Role:             "guest",
			Lang:             lang,
		}, externalId, AuthProviderTelegram)
	}

	// Unable to create user
	if err != nil {
		return nil, err
	}

	// Update user
	if user.TelegramUsername != sender.Username || user.Name != name || user.Lang != lang {
		user.TelegramUsername = sender.Username
		user.Name = name

		if !user.UserSetLang {
			user.Lang = lang
		}
		UpdateUser(dao, user)
	}

	return user, nil
}

func SetUserLang(dao *daos.Dao, user *User, lang string) {
	user.Lang = lang
	user.UserSetLang = true
	UpdateUser(dao, user)
}

type UsersAnalytics struct {
	Users  int
	Guests int
	Admins int
}

func GetUsersAnalytics(dao *daos.Dao) UsersAnalytics {
	analytics := UsersAnalytics{
		Users:  0,
		Guests: 0,
		Admins: 0,
	}

	UserQuery(dao).
		Select("count(*)").
		AndWhere(dbx.HashExp{"role": "user"}).
		Row(&analytics.Users)

	UserQuery(dao).
		Select("count(*)").
		AndWhere(dbx.HashExp{"role": "guest"}).
		Row(&analytics.Guests)

	UserQuery(dao).
		Select("count(*)").
		AndWhere(dbx.HashExp{"role": "admin"}).
		Row(&analytics.Admins)

	return analytics
}
