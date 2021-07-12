package main

import (
	"fmt"

	"github.com/GeekChomolungma/Everest-Base-Camp/chomows/huobi"
	"github.com/GeekChomolungma/Everest-Base-Camp/config"
	"github.com/GeekChomolungma/Everest-Base-Camp/handlerfactory"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("WELCOME TO BASE CAMP, READY TO SUBDUE CHOMOLUNGMA!")

	// config server
	config.Setup()

	// init factory
	handlerfactory.Init()

	// register gin server
	r := gin.Default()

	// for gateway restful action
	r.POST("/api/v1/Chomolungma/entrypoint", handlerfactory.FactoryImport)
	r.GET("/ws", huobi.WebsocketHandler)
	r.GET("/ws/v2", huobi.WebsocketHandlerV2)

	// server run!
	r.Run(config.ServerSetting.Host)
}
