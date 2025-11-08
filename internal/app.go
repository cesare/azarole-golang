package app

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Application struct {
	Config  Config
	Secrets Secrets

	database *sql.DB
}

func (application *Application) Database() *sql.DB {
	return application.database
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
		Config:   *config,
		Secrets:  *secrets,
		database: db,
	}
	return &app, nil
}
