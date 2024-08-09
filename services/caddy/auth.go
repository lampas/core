package caddy

import (
	"fmt"
	"lamas/models"
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/daos"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

type AuthCache struct {
	userId   string
	httpCode int
}

type AuthState struct {
	authorized bool
	userId     string
	httpCode   int
}

type AuthKeys struct {
	Keys      []string `json:"keys"`
	Challenge string   `json:"challenge"`
}

var authCookieName = "lamas_auth"
var authCache = expirable.NewLRU[string, AuthCache](1000, nil, time.Minute*10)
var authCookieEncryptionKey = ""

func ParseAuthState(dao *daos.Dao, c echo.Context) *AuthState {
	domain := c.Request().Header.Get("Auth-Domain")
	service := c.Request().Header.Get("Auth-Service")
	userIp := c.Request().Header.Get("Auth-User-Ip")
	cacheKey := domain + ":" + userIp
	authKey := service + ":" + domain

	// Get from cache or check if user exists
	cachedAuth, cacheValid := authCache.Get(cacheKey)

	// CACHE: Auth pass
	if cacheValid && cachedAuth.httpCode == http.StatusOK {
		return &AuthState{
			authorized: true,
			userId:     cachedAuth.userId,
			httpCode:   http.StatusOK,
		}
	}

	// CACHE: User not found
	if cacheValid && cachedAuth.userId == "" {
		return &AuthState{
			authorized: false,
			userId:     "",
			httpCode:   http.StatusNotFound,
		}
	}

	// Get user
	user, err := GetUserByDomain(dao, domain)
	userFound := err == nil
	if !userFound {
		authCache.Add(cacheKey, AuthCache{
			userId:   "",
			httpCode: http.StatusNotFound,
		})

		return &AuthState{
			authorized: false,
			userId:     "",
			httpCode:   http.StatusNotFound,
		}
	}

	// Validate auth token
	authToken, err := c.Cookie(authCookieName)
	if err == nil && authToken.Value != "" && ValidateCookieValue(authKey, authToken.Value) {
		authCache.Add(cacheKey, AuthCache{
			userId:   user.Id,
			httpCode: http.StatusOK,
		})

		return &AuthState{
			authorized: true,
			userId:     user.Id,
			httpCode:   http.StatusOK,
		}
	}

	// Unauthorized
	if !cacheValid {
		authCache.Add(cacheKey, AuthCache{
			userId:   user.Id,
			httpCode: http.StatusUnauthorized,
		})
	}

	return &AuthState{
		authorized: false,
		userId:     "",
		httpCode:   http.StatusUnauthorized,
	}
}

func GenerateAuthKeys() *AuthKeys {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)
	numbers := make(map[int]struct{})

	for len(numbers) < 4 {
		num := rng.Intn(99) + 1
		numbers[num] = struct{}{}
	}

	result := make([]string, 0, len(numbers))
	for num := range numbers {
		result = append(result, fmt.Sprintf("%02d", num))
	}

	return &AuthKeys{
		Keys:      result,
		Challenge: result[rand.Intn(4)],
	}
}

func RequestAuthChallenge(dao *daos.Dao, userId string) (*AuthKeys, error) {
	_, err := models.GetUserById(dao, userId)

	if err != nil {
		return nil, err
	}

	authKey := GenerateAuthKeys()

	return authKey, nil
}

func ReloadEncryptionKey(app *pocketbase.PocketBase) {
	config := models.GetConfiguration(app.Dao())
	if config.CookieEncryptionKey == authCookieEncryptionKey {
		return
	}

	app.Logger().Info("Reloading cookie encryption key")
	authCookieEncryptionKey = config.CookieEncryptionKey
}
