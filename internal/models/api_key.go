package models

import (
	"fmt"
	"strconv"
	"time"
)

type ApiKeyId uint32

type ApiKey struct {
	Id        ApiKeyId
	UserId    UserId
	Name      string
	Digest    string
	CreatedAt time.Time
}

func FromStringToApiKeyId(value string) (ApiKeyId, error) {
	uintValue, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		dummy := ApiKeyId(0)
		return dummy, fmt.Errorf("invalid api key id string %s: %s", value, err)
	}

	return ApiKeyId(uintValue), nil
}
