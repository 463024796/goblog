package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// PagesController 处理静态页面
type PagesController struct {
}

// Home 首页
func (*PagesController) Home(ctx *gin.Context) {
	fmt.Fprint(ctx.Writer, "<h1>Hello, 欢迎来到 goblog！</h1>")
}

// About 关于我们页面
func (*PagesController) About(ctx *gin.Context) {
	fmt.Fprint(ctx.Writer, "此博客是用以记录编程笔记，如您有反馈或建议，请联系 "+
		"<a href=\"mailto:summer@example.com\">summer@example.com</a>")
}

// NotFound 404 页面
func (*PagesController) NotFound(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
	fmt.Fprint(ctx.Writer,"<h1>请求页面未找到 :(</h1><p>如有疑惑，请联系我们。</p>")
}
