package handlers

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkplaceView struct {
	Id   models.WorkplaceId `json:"id"`
	Name string             `json:"name"`
}

func fromWorkplace(wp *models.Workplace) *WorkplaceView {
	return &WorkplaceView{
		Id:   wp.Id,
		Name: wp.Name,
	}
}

func RegisterWorkplacesHandlers(group *gin.RouterGroup, app *core.App) {
	group.GET("", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)
		workplaces, err := findWorkplaces(app, &currentUser)
		if err != nil {
			slog.Debug("failed to find workplaces", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		views := []WorkplaceView{}
		for _, wp := range workplaces {
			views = append(views, *fromWorkplace(&wp))
		}

		c.JSON(http.StatusOK, gin.H{
			"workplaces": views,
		})
	})
}

func findWorkplaces(app *core.App, user *models.User) ([]models.Workplace, error) {
	statement, err := app.Database().Prepare("select id, user_id, name from workplaces where user_id = $1 order by id")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findWorkplaces: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(user.Id)
	if err != nil {
		return nil, fmt.Errorf("query for findWorkplaces failed: %s", err)
	}

	workplaces := []models.Workplace{}
	for rows.Next() {
		var wp models.Workplace
		err = rows.Scan(&wp.Id, &wp.UserId, &wp.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to map row into workplace: %s", err)
		}

		workplaces = append(workplaces, wp)
	}

	return workplaces, nil
}
