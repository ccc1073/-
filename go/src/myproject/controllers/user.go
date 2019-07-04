package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"myproject/models"
)

func (c *MainController) showGet() {
	// 获取orm对象
	o := orm.NewOrm()
	// 执行某个操作函数 增删改查
	// 插入操作
	var user models.User
	user.Name = "lin"
	user.PassWord = "123456"

	// 查询操作
	var result models.User
	result.Id = 1

	err := o.Read(&user, "Id")

	fmt.Println(err)

	// 插入操作
	count, err := o.Insert(&user)
	if err != nil {
		beego.Error("插入失败")
	}
	beego.Info(count)

	// 更新操作
	var person models.User
	person.Id = 1
	// 更新之前需要先进行查询操作
	err2 := o.Read(&person)
	if err2 != nil {
		fmt.Println("查询有误更新失败")
	}
	person.Name = "lala"
	n, err := o.Update(&person)
	if err != nil {
		fmt.Println("更新失败")
	}
	beego.Info(n)

	// 删除操作
	var person2 models.User

	person2.Id = 1

	n, err3 := o.Delete(&person2)
	if err3 != nil {
		fmt.Println(err3)
	}
	beego.Info(n)

}

type UserContorller struct {
	beego.Controller
}

// 显示注册页面
func (this *UserContorller) Register() {

	this.TplName = "register.html"
}

func (this *UserContorller) HandlePost() {

	// 1.获取数据
	username := this.GetString("username")
	pwd := this.GetString("password")

	beego.Info(username, pwd)

	// 2.校验数据
	if username == "" || pwd == "" {
		this.Data["errmsg"] = "注册数据不完整"
		beego.Info("注册数据不完整")
		this.TplName = "register.html"
		return
	}

	// 3.操作数据
	o := orm.NewOrm()
	var user models.User

	user.Name = username
	user.PassWord = pwd
	o.Insert(&user)
	// 4.返回页面
	//this.Ctx.WriteString("注册成功")
	// 跳转
	this.Redirect("/login", 301)
}

// 展示登录页面
func (this *UserContorller) ShowLogin() {

	userName := this.Ctx.GetCookie("userName")
	if userName == "" {
		this.Data["userName"] = ""
		this.Data["checked"] = ""
	} else {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	}
	this.TplName = "login.html"
}

func (this *UserContorller) HandleLogin() {
	// 获取数据
	username := this.GetString("username")
	pwd := this.GetString("password")
	if username == "" || pwd == "" {
		this.Data["errmsg"] = "登录数据不完整"
		beego.Info("登录数据不完整")
		this.TplName = "login.html"
		return
	}

	// 查询
	o := orm.NewOrm()

	var user models.User

	user.Name = username

	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["errmsg"] = "用户不存在"
		this.TplName = "login.html"
		return
	}
	if user.PassWord != pwd {
		this.Data["errmsg"] = "密码错误"
		this.TplName = "login.html"
		return
	}

	// 返回
	//this.Ctx.WriteString("登录成功")
	data := this.GetString("remember")

	if data == "on" {
		this.Ctx.SetCookie("username", username, 100)
	} else {
		this.Ctx.SetCookie("username", "")
	}

	this.SetSession("username", username)

	this.Redirect("/showArticleList", 301)
}

// 退出登录
func (this *UserContorller) Logout() {

	// 删除session
	this.DelSession("username")
	this.Redirect("/login", 301)
}
