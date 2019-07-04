package main

import (
	"github.com/astaxie/beego"
	_ "myproject/models"
	_ "myproject/routers"
)

func main() {
	beego.AddFuncMap("prepage", ShowPrePage)
	beego.AddFuncMap("nextpage", ShowNextPage)
	beego.Run()
}

// 后台定义一个函数
func ShowPrePage(pageIndex int) int {
	if pageIndex == 1 {
		return pageIndex
	}

	return pageIndex - 1

}

func ShowNextPage(pageIndex int, pageCount int) int {

	if pageIndex == pageCount {
		return pageIndex
	}
	return pageIndex + 1

}
