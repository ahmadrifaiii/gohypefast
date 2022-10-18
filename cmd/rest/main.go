package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/labstack/echo/v4"
	"hypefast.io/services/config"
	"hypefast.io/services/config/env"
	"hypefast.io/services/config/redisconf"
	rest "hypefast.io/services/pkg/server/rest"
	"hypefast.io/services/pkg/utils/response"
	"hypefast.io/services/router"
)

func main() {
	env.LoadEnv()
	ctx := context.Background()

	// redis init
	redisServer := redisconf.GetMasterRedis()

	// configuration initialize
	conf := config.Configuration{
		RedisConnect: redisServer,
	}

	// server initialize
	server := rest.HTTPServer()

	// router handle
	router.SetupRouterShorter(ctx, server, conf)

	server.GET("/health", func(c echo.Context) error {
		return response.Success(c, response.Response{
			Message: "OK",
		})
	})

	// start echo server
	rest.StartServer(server)

	// Shutdown with gracefull handler
	rest.ShutdownServer(server)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
