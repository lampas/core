package caddy

import (
	origTemplate "html/template"
	"net/http"
	"net/url"

	"lamas/models"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func RegisterService(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		registry := template.NewRegistry()
		dao := app.Dao()

		ReloadEncryptionKey(app)

		/** TLS Domain Validation */
		e.Router.GET("/caddy/tls-domain-validate", func(c echo.Context) error {
			domain := c.QueryParams().Get("domain")
			user, err := GetUserByDomain(dao, domain)

			if err != nil {
				return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
			}

			return c.JSON(http.StatusOK, map[string]string{"username": user.Username})
		})

		/** Auth validation */
		e.Router.GET("/caddy/auth-verify", func(c echo.Context) error {
			// Allowed paths
			authPath := c.Request().Header.Get("Auth-Path")
			if authPath == "/auth" {
				return c.NoContent(http.StatusOK)
			}

			// Auth state
			authState := ParseAuthState(dao, c)
			if authState.httpCode == http.StatusUnauthorized {
				redirectUrl := url.QueryEscape(c.Request().Header.Get("Auth-Uri"))
				return c.Redirect(http.StatusUnauthorized, "/auth?redirectUrl="+redirectUrl)
			}

			return c.NoContent(authState.httpCode)
		})

		/** Home page */
		e.Router.GET("/", func(c echo.Context) error {
			// Get configuration
			config := models.LocalizeConfiguration(dao, DetectUserLang(c))

			// Render HTML
			html, err := registry.LoadFiles(
				"views/layout.html",
			).Render(map[string]any{
				"AppLang":        config.AppLang,
				"AppTitle":       config.AppTitle,
				"AppDescription": origTemplate.HTML(config.AppDescription),
				"BotUsername":    config.BotUsername,
			})

			if err != nil {
				return apis.NewNotFoundError("", err)
			}

			return c.HTML(http.StatusOK, html)
		})

		/** Auth */
		e.Router.GET("/auth", func(c echo.Context) error {
			// Get configuration
			config := models.LocalizeConfiguration(dao, DetectUserLang(c))

			// Render HTML
			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/auth.html",
			).Render(map[string]any{
				"AppLang":        config.AppLang,
				"AppTitle":       config.AppTitle,
				"AppDescription": origTemplate.HTML(config.AppDescription),
				"BotUsername":    config.BotUsername,
			})

			if err != nil {
				return apis.NewNotFoundError("", err)
			}

			return c.HTML(http.StatusOK, html)
		})

		e.Router.POST("/auth", func(c echo.Context) error {
			authState := ParseAuthState(dao, c)

			if authState.authorized {
				return c.NoContent(http.StatusOK)
			}

			if authState.httpCode == http.StatusNotFound {
				return c.NoContent(http.StatusNotFound)
			}

			authKey, err := RequestAuthChallenge(dao, authState.userId)
			if err != nil {
				return c.NoContent(http.StatusForbidden)
			}

			return c.JSON(http.StatusUnauthorized, map[string]string{"code": authKey.Challenge})
		})

		return nil
	})

	app.OnModelAfterUpdate(models.ConfigurationCollectionId).Add(func(e *core.ModelEvent) error {
		ReloadEncryptionKey(app)
		return nil
	})
}
