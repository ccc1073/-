package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"myproject/controllers"
)

func init() {

	beego.InsertFilter("/article/*", beego.BeforeExec, Filfterfunc)

	beego.Router("/", &controllers.MainController{})
	// 请求指定自定义方法 一个请求指定一个方法
	//beego.Router("login",&controllers.LogineController{},"get:showlogin;post:postfunc")
	// 多个请求指定一个方法
	//beego.Router("/index",&controllers.IndexController{},"get,post:HandleFunc")
	// 给所有请求指定一个方法
	//beego.Router("/index",&controllers.IndexController{},"*:HandleFunc")

	beego.Router("/register", &controllers.UserContorller{}, "get:Register;post:HandlePost") // 注册

	beego.Router("/login", &controllers.UserContorller{}, "get:ShowLogin;post:HandleLogin") // 登录

	beego.Router("/article/showArticleList", &controllers.ArticleController{}, "get:ShowArticleList") //文章首页
	// 展示添加文章页  操作添加文章
	beego.Router("/article/addArticle", &controllers.ArticleController{}, "get:ShowAddArticle;post:HandleAddArticle")

	// 文章详情
	beego.Router("/article/showArticleDetail", &controllers.ArticleController{}, "get:ShowArticleDetail")

	// 编辑文章
	beego.Router("/article/updateArticle", &controllers.ArticleController{}, "get:ShowUpdateArticle;post:HandleUpdateArticle")

	// 删除文章
	beego.Router("/article/deleteArticle", &controllers.ArticleController{}, "get:DeleteArticle")
	// 添加分类
	beego.Router("/article/addType", &controllers.ArticleController{}, "get:ShowArticleType;post:HandleAddType")

	beego.Router("/article/logout", &controllers.UserContorller{}, "get:Logout")

	beego.Router("/article/deletetype", &controllers.ArticleController{}, "get:DeleteType")

}

var Filfterfunc = func(ctx *context.Context) {
	username := ctx.Input.Session("username")
	if username == nil {
		ctx.Redirect(302, "/login")
		return
	}

}
