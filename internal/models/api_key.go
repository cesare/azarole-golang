package models

import "time"

type ApiKeyId uint32

type ApiKey struct {
	Id        ApiKeyId
	UserId    UserId
	Name      string
	Digest    string
	CreatedAt time.Time
}
