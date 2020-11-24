package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goblog/app/models/article"
	"goblog/app/policies"
	"goblog/app/requests"
	"goblog/pkg/auth"
	"goblog/pkg/route"
	"goblog/pkg/view"
	"net/http"
)

// ArticlesController 处理静态页面
type ArticlesController struct {
	BaseController
}

// Show 文章详情页面
func (ac *ArticlesController) Show(ctx *gin.Context) {

	// 1. 获取 URL 参数
	//id := route.GetRouteVariable("id", ctx)
	id := ctx.Param("id")
	fmt.Println(id)

	// 2. 读取对应的文章数据
	article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResposeForSQLError(ctx.Writer, err)
	} else {
		// ---  4. 读取成功，显示文章 ---
		view.Render(ctx.Writer, view.D{
			"Article":          article,
			"CanModifyArticle": policies.CanModifyArticle(article),
		}, "articles.show", "articles._article_meta")
	}
}

// Index 文章列表页
func (ac *ArticlesController) Index(ctx *gin.Context) {

	// 1. 获取结果集
	articles, pagerData, err := article.GetAll(ctx, 2)

	if err != nil {
		ac.ResposeForSQLError(ctx.Writer, err)
	} else {

		// ---  2. 加载模板 ---
		view.Render(ctx.Writer, view.D{
			"Articles":  articles,
			"PagerData": pagerData,
		}, "articles.index", "articles._article_meta")
	}
}

// Create 文章创建页面
func (*ArticlesController) Create(w http.ResponseWriter, r *http.Request) {
	view.Render(w, view.D{}, "articles.create", "articles._form_field")
}

// Store 文章创建页面
func (*ArticlesController) Store(w http.ResponseWriter, r *http.Request) {

	// 1. 初始化数据
	currentUser := auth.User()
	_article := article.Article{
		Title:  r.PostFormValue("title"),
		Body:   r.PostFormValue("body"),
		UserID: currentUser.ID,
	}

	// 2. 表单验证
	errors := requests.ValidateArticleForm(_article)

	// 3. 检测错误
	if len(errors) == 0 {
		// 创建文章
		_article.Create()
		if _article.ID > 0 {
			indexURL := route.Name2URL("articles.show", "id", _article.GetStringID())
			http.Redirect(w, r, indexURL, http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建文章失败，请联系管理员")
		}
	} else {
		view.Render(w, view.D{
			"Article": _article,
			"Errors":  errors,
		}, "articles.create", "articles._form_field")
	}
}

// Edit 文章更新页面
func (ac *ArticlesController) Edit(ctx *gin.Context) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", ctx)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResposeForSQLError(ctx.Writer, err)
	} else {

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResposeForUnauthorized(ctx)
		} else {
			// 4. 读取成功，显示编辑文章表单
			view.Render(ctx.Writer, view.D{
				"Article": _article,
				"Errors":  view.D{},
			}, "articles.edit", "articles._form_field")
		}
	}
}

// Update 更新文章
func (ac *ArticlesController) Update(ctx *gin.Context) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", ctx)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResposeForSQLError(ctx.Writer, err)
	} else {
		// 4. 未出现错误

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResposeForUnauthorized(ctx)
		} else {

			// 4.1 表单验证
			_article.Title = ctx.PostForm("title")
			_article.Body = ctx.PostForm("body")

			errors := requests.ValidateArticleForm(_article)

			if len(errors) == 0 {

				// 4.2 表单验证通过，更新数据
				rowsAffected, err := _article.Update()

				if err != nil {
					// 数据库错误
					ctx.Writer.WriteHeader(http.StatusInternalServerError)
					fmt.Fprint(ctx.Writer, "500 服务器内部错误")
					return
				}

				// √ 更新成功，跳转到文章详情页
				if rowsAffected > 0 {
					showURL := route.Name2URL("articles.show", "id", id)
					ctx.Redirect(http.StatusFound, showURL)
				} else {
					fmt.Fprint(ctx.Writer, "您没有做任何更改！")
				}
			} else {

				// 4.3 表单验证不通过，显示理由
				view.Render(ctx.Writer, view.D{
					"Article": _article,
					"Errors":  errors,
				}, "articles.edit", "articles._form_field")
			}
		}
	}
}

// Delete 删除文章
func (ac *ArticlesController) Delete(ctx *gin.Context) {

	// 1. 获取 URL 参数
	id := route.GetRouteVariable("id", ctx)

	// 2. 读取对应的文章数据
	_article, err := article.Get(id)

	// 3. 如果出现错误
	if err != nil {
		ac.ResposeForSQLError(ctx.Writer, err)
	} else {

		// 检查权限
		if !policies.CanModifyArticle(_article) {
			ac.ResposeForUnauthorized(ctx)
		} else {
			// 4. 未出现错误，执行删除操作
			rowsAffected, err := _article.Delete()

			// 4.1 发生错误
			if err != nil {
				// 应该是 SQL 报错了
				ctx.Writer.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(ctx.Writer, "500 服务器内部错误")
			} else {
				// 4.2 未发生错误
				if rowsAffected > 0 {
					// 重定向到文章列表页
					indexURL := route.Name2URL("articles.index")
					ctx.Redirect(http.StatusFound, indexURL)
				} else {
					// Edge case
					ctx.Writer.WriteHeader(http.StatusNotFound)
					fmt.Fprint(ctx.Writer, "404 文章未找到")
				}
			}
		}
	}
}
