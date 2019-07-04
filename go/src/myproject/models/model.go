package models

import (
	//"database/sql"
	//	"fmt"
	//"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct {
	Id       int
	Name     string
	PassWord string

	Articles []*Article `orm:"reverse(many)"`
}

type Article struct {
	Id       int       `orm:"pk;auto"`
	ArtiName string    `orm:"size(20)"`
	Atime    time.Time `orm:"auto_now"`
	Acount   int       `orm:"default(0);null"`
	Acontent string    `orm:"size(500)"`
	Aimg     string    `orm:"size(100)"`

	ArticleType *ArticleType `orm:"rel(fk)"`
	Users       []*User      `orm:"rel(m2m)"`
}
type ArticleType struct {
	Id      int
	TypeNme string `orm:"size(20)"`

	Articles []*Article `orm:"reverse(many)"`
}

func init() {
	//// 操作数据库代码
	//conn,err := sql.Open("mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	//defer conn.Close()//随手关闭数据库是个好习惯
	//// 第一个参数时数据库驱动
	//if err !=nil{
	//	beego.Info("连接错误")
	//
	//	beego.Error("连接错误",err)
	//}
	//// 插入数据 执行语句
	//_,err1:=conn.Exec("insert into itcast value (?,?)","123456","haha")
	//
	//fmt.Println(err1)
	//
	//// 查询数据
	//res,err2:=conn.Query("select name from itcast")
	//var name string
	//fmt.Println(err2)
	//for res.Next(){
	//
	//	res.Scan(&name)
	//	beego.Info(name)
	//}

	// orm操作数据库

	orm.RegisterDataBase("default", "mysql", "root:7894561230@(127.0.0.1:3306)/test?charset=utf8")

	// 创建表
	orm.RegisterModel(new(User), new(Article), new(ArticleType))

	// 生成表
	orm.RunSyncdb("default", false, true) // 第二个参数默认 false的话每次重启项目 不会从新生成表

}
