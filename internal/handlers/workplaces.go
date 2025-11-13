package handlers

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

type createWorkplaceParams struct {
	Name string `form:"name" binding:"required"`
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

	group.POST("", func(c *gin.Context) {
		currentUser := c.MustGet("currentUser").(models.User)

		var params createWorkplaceParams
		err := c.ShouldBind(&params)
		if err != nil {
			slog.Debug("failed to bind parameters", "error", err)
			c.Status(http.StatusBadRequest)
			return
		}

		workplace, err := createWorkplace(app, &currentUser, params.Name)
		if err != nil {
			slog.Debug("failed to create workplace", "error", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		view := fromWorkplace(workplace)
		c.JSON(http.StatusCreated, gin.H{
			"workplace": view,
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
	defer rows.Close()

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

func createWorkplace(app *core.App, user *models.User, name string) (*models.Workplace, error) {
	statement, err := app.Database().Prepare("insert into workplaces (user_id, name, created_at, updated_at) values ($1, $2, $3, $4) returning id, user_id, name")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for createWorkplace: %s", err)
	}
	defer statement.Close()

	now := time.Now().UTC()
	rows, err := statement.Query(user.Id, name, now, now)
	if err != nil {
		return nil, fmt.Errorf("query for createWorkplace failed: %s", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("inserted row missing")
	}

	var workplace models.Workplace
	err = rows.Scan(&workplace.Id, &workplace.UserId, &workplace.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to map row into workplace: %s", err)
	}

	return &workplace, nil
}
