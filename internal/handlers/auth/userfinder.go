package auth

import (
	app "azarole/internal"
	"azarole/internal/models"
	"database/sql"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type FinderResult struct {
	UserId models.UserId
}

type UserFinder struct {
	application *app.Application
	identifier  string
}

func NewUserFinder(application *app.Application, identifier string) *UserFinder {
	return &UserFinder{
		application: application,
		identifier:  identifier,
	}
}

func (finder *UserFinder) Execute(c *gin.Context) (*FinderResult, error) {
	var result FinderResult

	err := finder.application.WithTransaction(c, func(tx *sql.Tx) error {
		r, err := finder.findUserId(tx)
		if err != nil {
			return err
		}

		if r != nil {
			result.UserId = r.UserId
			return nil
		}

		user, err := finder.createUser(tx)
		if err != nil {
			return err
		}

		err = finder.createGoogleAuthUser(tx, user.Id)
		if err != nil {
			return err
		}

		result.UserId = user.Id
		return nil
	})

	return &result, err
}

func (finder *UserFinder) findUserId(tx *sql.Tx) (*FinderResult, error) {
	statement, err := tx.Prepare("select user_id as id from google_authenticated_users where uid = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for findUserId: %s", err)
	}
	defer statement.Close()

	rows, err := statement.Query(finder.identifier)
	if err != nil {
		return nil, fmt.Errorf("query for findUserId failed: %s", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	var result FinderResult
	err = rows.Scan(&result.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row for findUserId: %s", err)
	}

	return &result, nil
}

func (finder *UserFinder) createUser(tx *sql.Tx) (*models.User, error) {
	statement, err := tx.Prepare("insert into users (created_at) values ($1) returning id")
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement for createUser: %s", err)
	}
	defer statement.Close()

	var user models.User
	now := time.Now().UTC()
	err = statement.QueryRow(now).Scan(&user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to insert new user: %s", err)
	}

	return &user, nil
}

func (finder *UserFinder) createGoogleAuthUser(tx *sql.Tx, userId models.UserId) error {
	statement, err := tx.Prepare("insert into google_authenticated_users (user_id, uid, created_at) values ($1, $2, $3)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement for createGoogleAuthUser: %s", err)
	}
	defer statement.Close()

	now := time.Now().UTC()
	_, err = statement.Exec(userId, finder.identifier, now)
	if err != nil {
		return fmt.Errorf("failed to insert google_authenticated_user: %s", err)
	}

	return nil
}
