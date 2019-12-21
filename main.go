package main

import (
	_ "goblog/routers"
	"github.com/astaxie/beego"
)



func main() {
	beego.AddFuncMap("prevpage",DoPrevPage)
	beego.AddFuncMap("nextpage",DoNextPage)
	beego.Run()
}

//分页
func DoPrevPage(data int)(int){
	pageIndex:=data-1
	return pageIndex
}

func DoNextPage(data int)(int){
	pageIndex:=data+1
	return pageIndex
}