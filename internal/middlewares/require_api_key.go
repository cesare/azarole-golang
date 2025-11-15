package middlewares

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"azarole/internal/resources"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireApiKey(app *core.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			slog.Debug("failed to extract token", "error", err)
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		digest := digestToken(app, token)
		apiKey, err := findApiKey(app, digest)
		if err != nil {
			slog.Debug("failed to find api key", "error", err)
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		rs := resources.NewUserResources(app)
		user, err := rs.Find(apiKey.UserId)
		if err != nil {
			slog.Debug("failed to find user", "error", err)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}
		if user == nil {
			slog.Debug("user missing")
			c.Status(http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("currentUser", *user)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	header := c.GetHeader("Authorization")
	if header == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	values := strings.SplitN(header, " ", 2)
	if len(values) != 2 {
		return "", fmt.Errorf("invaled Authorization header value: %s", header)
	}

	if values[0] != "Bearer" {
		return "", fmt.Errorf("invaled Authorization header value: %s", header)
	}

	return values[1], nil
}

func digestToken(app *core.App, token string) string {
	secretKey := app.Secrets.ApiKey.DigestingSecretKey.Bytes()
	h := hmac.New(sha256.New, secretKey)
	h.Write([]byte(token))
	digest := h.Sum(nil)

	return hex.EncodeToString(digest)
}

func findApiKey(app *core.App, digest string) (*models.ApiKey, error) {
	statement, err := app.Database().Prepare("select id, user_id, name, digest, created_at from api_keys where digest = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findApiKey: %s", err)
	}
	defer statement.Close()

	var apiKey models.ApiKey
	err = statement.QueryRow(digest).Scan(&apiKey.Id, &apiKey.UserId, &apiKey.Name, &apiKey.Digest, &apiKey.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to find api key: %s", err)
	}

	return &apiKey, nil
}
