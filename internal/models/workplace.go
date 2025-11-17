package models

import (
	"fmt"
	"strconv"
)

type WorkplaceId uint32

type Workplace struct {
	Id     WorkplaceId `db:"id"`
	UserId UserId      `db:"user_id"`
	Name   string      `db:"name"`
}

func FromStringToWorkplaceId(value string) (WorkplaceId, error) {
	uintValue, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		dummy := WorkplaceId(0)
		return dummy, fmt.Errorf("invalid workplace id string %s: %s", value, err)
	}

	return WorkplaceId(uintValue), nil
}
