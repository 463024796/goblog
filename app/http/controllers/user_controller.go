package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goblog/app/models/article"
	"goblog/app/models/user"
	"goblog/pkg/logger"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

// UserController 用户控制器
type UserController struct {
	BaseController
}

// Show 用户个人页面
func (uc *UserController) Show(ctx *gin.Context) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", ctx)

	// 2. 读取对应的文章数据
	_user, err := user.Get(id)

	// 3. 如果出现错误
	if err != nil {
		uc.ResposeForSQLError(ctx.Writer, err)
	} else {
		// ---  4. 读取成功，显示用户文章列表 ---
		articles, err := article.GetByUserID(_user.GetStringID())
		if err != nil {
			logger.LogError(err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(ctx.Writer, "500 服务器内部错误")
		} else {
			view.Render(ctx.Writer, view.D{
				"Articles": articles,
			}, "articles.index", "articles._article_meta")
		}
	}
}
