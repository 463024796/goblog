package main

import (
	gin "github.com/gin-gonic/gin"
	"goblog/bootstrap"
	"goblog/config"
)

func init() {
	// 初始化配置信息
	config.Initialize()
}

func main() {
	// 初始化 SQL
	bootstrap.SetupDB()

	server := gin.New()

	// 初始化路由绑定
	bootstrap.SetupRoute(server)

	server.Run(":8885")

	//http.ListenAndServe(":"+c.GetString("app.port"), middlewares.RemoveTrailingSlash(router))
}
