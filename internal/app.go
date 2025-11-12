package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	Config  *Config
	Secrets *Secrets

	database *sql.DB
}

func (application *Application) Database() *sql.DB {
	return application.database
}

func (application *Application) WithTransaction(c *gin.Context, f func(*sql.Tx) error) error {
	ctx := context.WithoutCancel(c.Request.Context())
	tx, err := application.Database().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %s", err)
	}

	defer func() {
		p := recover()
		if p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	err = f(tx)
	return err
}

func LoadApplication(configPath string) (*Application, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	secrets, err := LoadSecrets()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", config.Database.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %s", err)
	}

	app := Application{
		Config:   config,
		Secrets:  secrets,
		database: db,
	}
	return &app, nil
}
