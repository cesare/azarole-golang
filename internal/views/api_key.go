package views

import (
	"azarole/internal/models"
	"time"
)

type ApiKeyView struct {
	Id        models.ApiKeyId `json:"id"`
	Name      string          `json:"name"`
	CreatedAt time.Time       `json:"createdAt"`
}

func FromApiKey(apiKey *models.ApiKey) *ApiKeyView {
	return &ApiKeyView{
		Id:        apiKey.Id,
		Name:      apiKey.Name,
		CreatedAt: apiKey.CreatedAt,
	}
}
