package models

type WorkplaceId uint32

type Workplace struct {
	Id     WorkplaceId `db:"id"`
	UserId UserId      `db:"user_id"`
	Name   string      `db:"name"`
}
