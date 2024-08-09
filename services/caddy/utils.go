package caddy

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"

	"lamas/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/daos"
)

func DetectUserLang(c echo.Context) string {
	acceptLanguage := c.Request().Header.Get("Accept-Language")
	if strings.Contains(acceptLanguage, "uk") {
		return "uk"
	}

	if strings.Contains(acceptLanguage, "ru") {
		return "ru"
	}

	return ""
}

func GetUserByDomain(dao *daos.Dao, domain string) (*models.User, error) {
	regexp := regexp.MustCompile(`(\.tv\..*)$`)
	username := regexp.ReplaceAllString(domain, "")

	if !regexp.MatchString(domain) {
		return nil, errors.New("invalid domain")
	}

	if username == "" {
		return nil, errors.New("domain is required")
	}

	user, err := models.GetUserByUsername(dao, username)

	if err != nil {
		return nil, errors.New("user record not found")
	}

	if user.Role == "guest" {
		return nil, errors.New("user is not allowed to access this service")
	}

	return user, nil
}

func EncryptCookieValue(value string) string {
	h := hmac.New(sha256.New, []byte(authCookieEncryptionKey))
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

func ValidateCookieValue(value string, hash string) bool {
	computedHMAC := EncryptCookieValue(value)
	return hmac.Equal([]byte(computedHMAC), []byte(hash))
}
