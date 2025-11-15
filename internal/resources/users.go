package resources

import (
	"azarole/internal/core"
	"azarole/internal/models"
	"fmt"
)

type UserResources struct {
	app *core.App
}

func NewUserResources(app *core.App) *UserResources {
	return &UserResources{
		app: app,
	}
}

func (r *UserResources) Find(userId models.UserId) (*models.User, error) {
	statement, err := r.app.Database().Prepare("select id from users where id = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for loadCurrentUser: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("query for loadCurrentUser failed: %s", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var user models.User
	rows.Scan(&user.Id)

	return &user, nil
}
