package views

import "azarole/internal/models"

type UserView struct {
	Id models.UserId `json:"id"`
}

func FromUser(user *models.User) *UserView {
	return &UserView{
		Id: user.Id,
	}
}
