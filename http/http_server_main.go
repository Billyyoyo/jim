package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jim/http/core"
	"jim/http/router"
)

func main() {
	fmt.Println("------startup http server-----")
	gin.SetMode(core.AppConfig.Server.Mode)
	r := gin.Default()
	r.Use(core.Logger())
	r.Use(core.CROS())
	router.Route(r)
	r.Run(fmt.Sprintf("%s:%d", core.AppConfig.Server.Host, core.AppConfig.Server.Port))
}
