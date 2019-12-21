package routers

import (
	"goblog/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	
)
//注释
func init() {
	//这里可以匹配路由，在前边判断是否有session
    beego.InsertFilter("/admin", beego.BeforeRouter,FilterFunc)
    beego.InsertFilter("/pic", beego.BeforeRouter,FilterFunc)
    beego.InsertFilter("/category", beego.BeforeRouter,FilterFunc)
    beego.InsertFilter("/blog", beego.BeforeRouter,FilterFunc)

    
    beego.Router("/", &controllers.MainController{})
    beego.Router("/login", &controllers.MainController{},"get:ShowLogin;post:DoLogin")
    beego.Router("/admin", &controllers.MainController{},"get:ShowAdminIndex")
    beego.Router("/category", &controllers.MainController{},"get:ShowCategory")
    beego.Router("/categoryadd", &controllers.MainController{},"get:ShowCategoryAdd;post:DoCategoryAdd")
    beego.Router("/pic", &controllers.MainController{},"get:ShowPic")
    beego.Router("/picadd", &controllers.MainController{},"get:ShowPicAdd;post:DoPicAdd")
    beego.Router("/index", &controllers.MainController{},"get:ShowIndex")
    beego.Router("/blog", &controllers.MainController{},"get:ShowAddBlog;post:DoAddBlog")
    beego.Router("/bloglist", &controllers.MainController{},"get:ShowBlog")
    beego.Router("/detail", &controllers.MainController{},"get:ShowDetail")
    beego.Router("/talk", &controllers.MainController{},"post:DoAddTalk")
    beego.Router("/like", &controllers.MainController{},"get:DoAddLike")
    beego.Router("/edit", &controllers.MainController{},"get:ShowBlogEdit;post:DoEdit")
    beego.Router("/del", &controllers.MainController{},"get:DoDel")
    beego.Router("/logout", &controllers.MainController{},"get:Logout")
  
}

var FilterFunc = func (ctx *context.Context){
	name:=ctx.Input.Session("goblog_user")
	if name==nil{
		ctx.Redirect(302,"/")
	
	}
}