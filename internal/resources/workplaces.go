package resources

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"fmt"
	"time"
)

type WorkplaceResources struct {
	app  *core.App
	user *models.User
}

func NewWorkplaceResources(app *core.App, user *models.User) *WorkplaceResources {
	return &WorkplaceResources{
		app:  app,
		user: user,
	}
}

func (r *WorkplaceResources) List() ([]models.Workplace, error) {
	statement, err := r.app.Database().Prepare("select id, user_id, name from workplaces where user_id = $1 order by id")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findWorkplaces: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(r.user.Id)
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

func (r *WorkplaceResources) Create(name string) (*models.Workplace, error) {
	statement, err := r.app.Database().Prepare("insert into workplaces (user_id, name, created_at, updated_at) values ($1, $2, $3, $4) returning id, user_id, name")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for createWorkplace: %s", err)
	}
	defer statement.Close()

	now := time.Now().UTC()
	rows, err := statement.Query(r.user.Id, name, now, now)
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
