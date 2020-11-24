package bootstrap

import (
	"github.com/gin-gonic/gin"
	"goblog/routes"
)

// SetupRoute 路由初始化
func SetupRoute(server *gin.Engine) {

	routes.RegisterWebRoutes(server)

}
