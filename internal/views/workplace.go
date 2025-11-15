package views

import "azarole/internal/models"

type WorkplaceView struct {
	Id   models.WorkplaceId `json:"id"`
	Name string             `json:"name"`
}

func FromWorkplace(wp *models.Workplace) *WorkplaceView {
	return &WorkplaceView{
		Id:   wp.Id,
		Name: wp.Name,
	}
}
