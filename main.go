package main

import (
	"fmt"

	"github.com/GeekChomolungma/Everest-Base-Camp/config"
	"github.com/GeekChomolungma/Everest-Base-Camp/handlerfactory"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("WELCOME TO QUANTS GATEWAY")

	// config server
	config.Setup()

	// init factory
	handlerfactory.Init()

	// register gin server
	r := gin.Default()

	// for gateway restful action
	r.POST("/api/v1/Chomolungma/entrypoint", handlerfactory.FactoryImport)

	// server run!
	r.Run(config.ServerSetting.Host)
}
