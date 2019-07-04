package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"myproject/models"
	"path"
	"time"
)

type ArticleController struct {
	beego.Controller
}

// 文章列表页展示 以及翻页
func (this *ArticleController) ShowArticleList() {

	// 高级查询
	// 指定表
	userName := this.GetSession("username")
	if userName == nil {
		this.Redirect("/login", 301)
		return
	}
	o := orm.NewOrm()

	qs := o.QueryTable("Article")
	var articles []models.Article

	//_,err:=qs.All(&articles)
	//if err!=nil{
	//	fmt.Println("查询数据错误")
	//
	//}
	typeName := this.GetString("select")

	var count int64

	// 查询总记录数
	//count,_=qs.Count()
	pageSize := 2

	//pageCount:=count/int64(pageSize)

	// 天花板函数 和地板函数

	pageIndex, err := this.GetInt("pageIndex")
	if err != nil {
		pageIndex = 1
	}
	// 获取分页数据 第一个参数 获取几条 第二个参数 从那条数据开始获取
	start := (pageIndex - 1) * pageSize
	//_,err=qs.Limit(pageSize,start).RelatedSel("ArticleType").All(&articles)

	if typeName == "" {
		count, _ = qs.Count()

	} else {
		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).Count()

	}
	pageCount := math.Ceil(float64(count) / float64(pageSize))

	// 获取文章类型
	var types models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	if typeName == "" {

		qs.Limit(pageSize, start).RelatedSel("ArticleType").All(&articles)

	} else {

		qs.Limit(pageSize, start).RelatedSel("ArticleType").Filter("ArticleType__TypeName", typeName).All(&articles)

	}

	this.Data["pageIndex"] = pageIndex
	this.Data["count"] = count
	this.Data["pageCount"] = int(pageCount)
	this.Data["articles"] = articles
	this.Data["types"] = types
	this.Data["typeName"] = typeName

	this.TplName = "index.html"

}

// 展示文章添加
func (this *ArticleController) ShowAddArticle() {

	o := orm.NewOrm()

	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	this.Data["types"] = types
	this.TplName = "add.html"
}

// 添加文章
func (this *ArticleController) HandleAddArticle() {

	// 1获取数据
	articlename := this.GetString("articleName")
	content := this.GetString("content")

	// 2校验数据
	if articlename == "" || content == "" {
		this.Data["errmsg"] = "添加不完整"
		this.TplName = "add.html"
		return
	}

	// 处理文件上传
	file, head, err := this.GetFile("uploadname") // head 为文件头

	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = "add.html"
		return
	}
	defer file.Close()

	this.SaveToFile("uploadname", "./static/img"+head.Filename)

	// 上传文件注意点 文件大小 文件格式 防止重名
	if head.Size > 1024*5 {
		this.Data["errmsg"] = "文件过大 请重新上传"
		this.TplName = "add.html"
		return
	}

	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误 请重新上传"
		this.TplName = "add.html"
		return
	}

	// 防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	// 存储
	this.SaveToFile("uploadname", "./static/img/"+fileName)

	// 3处理数据
	o := orm.NewOrm()

	var article models.Article

	article.ArtiName = articlename
	article.Acontent = content

	article.Aimg = "/static/img" + fileName

	TypeName := this.GetString("select")
	var articleType models.ArticleType

	articleType.TypeNme = TypeName

	o.Read(&articleType)

	article.ArticleType = &articleType

	o.Insert(&article)

	this.Redirect("/showArticleList", 301)

}

// 展示文章详情页面
func (this *ArticleController) ShowArticleDetail() {

	// 获取数据
	id, err := this.GetInt("articleId")

	if err != nil {
		beego.Info("传递连接错误")
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.QueryTable("Article").RelatedSel("ArticleType").Filter("Id", id).One(&article)
	//o.Read(&article)

	article.Acount += 1
	o.Update(&article)

	// 多对多插入浏览记录
	m2m := o.QueryM2M(&article, "Users")
	userName := this.GetSession("username")
	if userName == nil {
		this.Redirect("/login", 301)
		return
	}
	var user models.User
	// interface断言
	user.Name = userName.(string)
	o.Read(&user, "Name")
	m2m.Add(user)

	var users []models.User
	// 查询多对多
	//o.LoadRelated(&article,"Users")
	o.QueryTable("User").Filter("Articles__Article__Id", id).Distinct().All(&users)
	this.Data["article"] = article
	this.Data["users"] = users

	this.TplName = "content.html"
}

// 编辑文章页面
func (this *ArticleController) ShowUpdateArticle() {

	id, err := this.GetInt("articleId")

	if err != nil {
		beego.Info("请求文章错误")
		return
	}
	// 数据处理
	o := orm.NewOrm()
	var article models.Article

	article.Id = id
	o.Read(&article)

	// 返回视图
	this.Data["article"] = article
	this.TplName = "update.html"

}

// 封装文件处理函数
func UploadFile(this *beego.Controller, filePath string, view string) string {

	// 处理文件上传
	file, head, err := this.GetFile(filePath) // head 为文件头

	if head.Filename == "" {
		return "NoImg"
	}

	if err != nil {
		this.Data["errmsg"] = "文件上传失败"
		this.TplName = view
		return ""
	}
	defer file.Close()

	this.SaveToFile(filePath, "./static/img"+head.Filename)

	// 上传文件注意点 文件大小 文件格式 防止重名
	if head.Size > 1024*5 {
		this.Data["errmsg"] = "文件过大 请重新上传"
		this.TplName = view
		return ""
	}

	ext := path.Ext(head.Filename)
	if ext != ".jpg" && ext != ".png" && ext != ".jpeg" {
		this.Data["errmsg"] = "文件格式错误 请重新上传"
		this.TplName = view
		return ""
	}

	// 防止重名
	fileName := time.Now().Format("2006-01-02-15:04:05") + ext
	// 存储
	this.SaveToFile(filePath, "./static/img/"+fileName)

	return "/static/img/" + fileName
}

// 处理编辑界面数据
func (this *ArticleController) HandleUpdateArticle() {
	// 获取数据
	id, err := this.GetInt("articleId")
	content := this.GetString("content")
	articleName := this.GetString("articleName")

	filepath := UploadFile(&this.Controller, "uploadname", "update.html")

	if err != nil || articleName == "" || content == "" || filepath == "" {
		beego.Info("请求错误")
		return
	}

	// 数据处理
	o := orm.NewOrm()
	var article models.Article
	article.Id = id

	err = o.Read(&article)

	if err != nil {
		beego.Info("更新的文章不存在")
	}

	article.ArtiName = articleName
	article.Acontent = content
	if filepath != "NoImg" {
		article.Aimg = filepath
	}

	o.Update(&article)

	this.Redirect("/showArticleList", 302)

}

// 删除文章
func (this *ArticleController) DeleteArticle() {

	id, err := this.GetInt("articleId")

	if err != nil {
		beego.Info("删除文章请求路径错误")
		return
	}

	o := orm.NewOrm()
	var article models.Article
	article.Id = id
	o.Delete(&article)

	this.Redirect("/showArticle", 302)

}

// 展示添加分类
func (this *ArticleController) ShowArticleType() {

	o := orm.NewOrm()
	var types []models.ArticleType
	o.QueryTable("ArticleType").All(&types)

	this.Data["type"] = types

	this.TplName = "addType.html"
}

// 操作添加分类
func (this *ArticleController) HandleAddType() {

	typeName := this.GetString("typeName")

	if typeName == "" {
		beego.Info("信息不完整,请重新输入")

		return
	}
	// 处理数据
	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.TypeNme = typeName
	o.Insert(&articleType)
	// 插入操作

	this.Redirect("/addType", 301)
}

// 删除分类
func (this *ArticleController) DeleteType() {
	id, err := this.GetInt("typeid")

	if err != nil {
		beego.Error("删除类型错误", err)
		return
	}

	o := orm.NewOrm()
	var articleType models.ArticleType
	articleType.Id = id
	o.Delete(&articleType)

	this.Redirect("/article/addType", 302)

}
