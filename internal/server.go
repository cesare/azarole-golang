package app

import "github.com/gin-gonic/gin"

func Engine(config *Config) (*gin.Engine, error) {
	engine := gin.Default()
	return engine, nil
}
