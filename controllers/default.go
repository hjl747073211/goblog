package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"goblog/models"
	"time"
	"strconv"
	"path"
	"strings"
	"math"
	"bytes"
	"encoding/gob"
	// "fmt"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "login.html"
}

func (c *MainController) Logout() {
	c.DelSession("goblog_user") 
	c.TplName = "login.html"
}

func (c *MainController) ShowLogin() {
	//获取cookie
	// name:=c.Ctx.GetCookie("goblog_user")
	//获取session
	// name:=c.GetSession("goblog_user")
	// beego.Info(name)
	// if name!="" && name!=nil{
	// 	c.Data["goblog_user"]=name
	// 	c.TplName = "admin_index.html"
	// }else{
		c.TplName = "login.html"
	// }

}

func (c *MainController) DoLogin() {
	userName:=c.GetString("username")
	pwd:=c.GetString("password")
	if userName=="" || pwd == ""{
		beego.Info("数据不能为空")
		c.Redirect("/login",302)
		return
	}
	o:= orm.NewOrm()
	user := models.User{}
	user.Name=userName
	err:=o.Read(&user,"Name")
	if err!=nil{
		beego.Info("用户不存在")
		// c.TplName="login.html"
		c.Redirect("/login",302)
		return
	}
	if user.Pwd != pwd{
		beego.Info("pwd err")
		c.Redirect("/login",302)
		return
	}

	//设置cookie
	// c.Ctx.SetCookie("goblog_user",userName,time.Second*3600)
	//设置session
	c.SetSession("goblog_user",userName)

	c.Data["goblog_user"]=userName
	c.TplName = "admin_index.html"

}

func (c *MainController) ShowAdminIndex() {
	c.TplName = "admin_index.html"
}

func (c *MainController) ShowCategory() {
	o:=orm.NewOrm()
	var categories []models.Category
	count,err:=o.QueryTable("Category").All(&categories)
	if err!=nil{
		beego.Info("select err")
		return 
	}
	c.Data["count"]=count
	c.Data["categories"]=categories
	c.TplName = "admin_category.html"
}

func (c *MainController) ShowCategoryAdd() {
	c.TplName = "admin_category_add.html"
}

func (c *MainController) DoCategoryAdd() {
	category:=c.GetString("category")
	remark:=c.GetString("remark")

	if category == "" || remark == ""{
		beego.Info("数据不能为空")
		c.Redirect("/categoryadd",302)
		return 
	}
	o := orm.NewOrm()
	cate:=models.Category{}
	cate.Category=category
	cate.Remark=remark
	_,err:=o.Insert(&cate)
	if err!=nil{
		beego.Info("插入失败")
		c.Redirect("/categoryadd",302)
		return 
	}
	c.Ctx.WriteString("添加category成功")	
	// c.TplName="admin_category_add.html"
	// c.Redirect("/categoryadd",200)  //跳转状态码302,写错不跳

}

func (c *MainController) ShowPic() {
	var images []models.Image

	//redis取出数据
	conn,_:=redis.Dial("tcp",":6379")
	defer conn.Close()
	
	buffer,err:=redis.Bytes(conn.Do("get","images"))

	//redis获取不到，就从数据库查，否则就直接取redis
	if err!=nil{
		o:=orm.NewOrm()
		//查询image表的所有数据，放在images里边
		_,err:=o.QueryTable("Image").All(&images)
		if err!=nil{
			beego.Info("select err")
			return 
		}
	}else{
		dec := gob.NewDecoder(bytes.NewReader(buffer))
		dec.Decode(&images)
		beego.Info(images)
	
	}


	c.Data["count"]=len(images)
	c.Data["images"]=images
	c.TplName = "admin_big_pic.html"
}

func (c *MainController) ShowPicAdd() {
	c.TplName = "admin_big_pic_add.html"
}

func (c *MainController) DoPicAdd() {
	title:=c.GetString("title")
	content:=c.GetString("content")
	if title=="" || content=="" {
		c.Ctx.WriteString("请输入标题和内容")
		return 
	}

	//上传图片文件,f用不上，直接关闭
	f,h,err := c.GetFile("image")
	defer f.Close()

	//限定格式
	fileext:=path.Ext(h.Filename)
	if fileext!=".jpg" && fileext!=".png" {
		beego.Info("upload fail")
		return
	}
	//限定大小
	if h.Size < 100000 {
		beego.Info("upload file size err")
	}
	//当前时间戳（纳秒）
	now := time.Now().UnixNano()


	var imagePath string
	if err!=nil{
		beego.Info("upload fail")
		return 
	}else{
		//if里面定义的变量作用域只在if，外边无法使用，这里我将imagePath定在外面
		imagePath="./static/upload/"+strconv.FormatInt(now,10)+h.Filename 
		c.SaveToFile("image",imagePath)
	}

	o := orm.NewOrm()
	img:=models.Image{}
	img.Image=imagePath
	img.Title=title
	img.Content=content
	_ ,errInsert := o.Insert(&img)
	if errInsert!=nil{
		beego.Info("插入失败")
		c.Redirect("/picadd",302)
		return 
	}

	
	var images []models.Image
	//查询image表的所有数据，放在images里边
	_,errSel:=o.QueryTable("Image").All(&images)
	if errSel!=nil{
		beego.Info("select err")
		return 
	}
	// 添加新数据时候，更新redis数据
	conn,_:=redis.Dial("tcp",":6379")
	defer conn.Close()
	//序列化后存入
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	enc.Encode(images)

	conn.Do("set","images",buffer.Bytes())
	c.Redirect("/pic",302)

}

func (c *MainController) ShowIndex() {
	o:=orm.NewOrm()
	var blogs []models.Blog
	var hotblogs []models.Blog
	var big_imgs []models.Image
	var talks []models.Talk
	_,err_img:=o.QueryTable("Image").All(&big_imgs)
	qs:=o.QueryTable("Blog")
	totalcount,err_count:=qs.Count()

	//分页
	pageSize:=5
	//向上取整，Math.Floor向下取整
	pageCount:=int(math.Ceil(float64(totalcount)/float64(pageSize)) )
	
	
	pageIndex,pageindex_err:= c.GetInt("pageIndex")
	if pageindex_err!=nil{
		pageIndex=1  //默认首页
	}
	
	pageStart:=(pageIndex-1)*pageSize
	_,err:=o.QueryTable("Blog").Limit(pageSize,pageStart).All(&blogs)
	
	FirstPage:=false
	if pageIndex==1 {
		FirstPage=true
	}
	c.Data["FirstPage"]=FirstPage
	LastPage:=false
	if pageIndex==pageCount {
		LastPage=true
	}
	c.Data["LastPage"]=LastPage



	//使用QueryBuilder查询
	qb,_:=orm.NewQueryBuilder("mysql")
	qb.Select("*").From("talk").OrderBy("id").Desc().Limit(5)
	sql:=qb.String()
	beego.Info(sql)
	o.Raw(sql).QueryRows(&talks)
	c.Data["talks"]=talks
	//使用QueryBuilder查询，热门文章
	qb1,_:=orm.NewQueryBuilder("mysql")
	qb1.Select("*").From("blog").OrderBy("look_num").Desc().Limit(3)
	sql2:=qb1.String()
	o.Raw(sql2).QueryRows(&hotblogs)
	c.Data["hotblogs"]=hotblogs

	if err!=nil && err_img!=nil&& err_count!=nil {
		beego.Info("select err")
		return 
	}
	c.Data["big_imgs"]=big_imgs
	c.Data["count"]=totalcount
	c.Data["blogs"]=blogs
	c.Data["pagecount"]=pageCount
	c.Data["pageIndex"]=pageIndex
	
	c.TplName = "index.html"
}


func (c *MainController) ShowAddBlog() {
	o:=orm.NewOrm()
	var categories []models.Category
	_,err:=o.QueryTable("Category").All(&categories)
	if err!=nil{
		beego.Info("select err")
		return 
	}
	c.Data["categories"]=categories
	c.TplName = "admin_blog_add.html"
}

func (c *MainController) DoAddBlog() {
	title:=c.GetString("title")
	creater:=c.GetString("creater")
	remark:=c.GetString("remark")
	category:=c.GetStrings("categorys")
	// beego.Info(category,string(subsCodes))
	content:=c.GetString("editorValue")

	if title=="" || creater=="" || remark==""{
		c.Ctx.WriteString("请全部填写完毕再提交哦！")
		return 
	}

	//上传图片文件,f用不上，直接关闭
	f,h,err := c.GetFile("img")
	defer f.Close()

	//限定格式
	fileext:=path.Ext(h.Filename)
	if fileext!=".jpg" && fileext!=".png" {
		beego.Info("upload fail")
		return
	}
	//限定大小
	if h.Size < 100000 {
		beego.Info("upload file size err")
	}
	//当前时间戳（纳秒）
	now := time.Now().UnixNano()


	var imagePath string
	if err!=nil{
		beego.Info("upload fail")
		return 
	}else{
		//if里面定义的变量作用域只在if，外边无法使用，这里我将imagePath定在外面
		imagePath="./static/upload/"+strconv.FormatInt(now,10)+h.Filename 
		c.SaveToFile("img",imagePath)
	}

	o := orm.NewOrm()
	blog:=models.Blog{}
	blog.Img=imagePath
	blog.Title=title
	blog.Remark=remark
	blog.Creater=creater
	blog.Category=strings.Join(category,",")  //将切片分隔成字符串
	blog.Content=content
	_ ,errInsert := o.Insert(&blog)
	if errInsert!=nil{
		beego.Info("提交失败")
		c.Redirect("/blog",302)
		return 
	}
	c.Redirect("/bloglist",302)
}

func (c *MainController) ShowBlog() {
	o:=orm.NewOrm()
	var blogs []models.Blog
	
	count,err:=o.QueryTable("Blog").All(&blogs)
	if err!=nil{
		beego.Info("select err")
		return 
	}

	c.Data["count"]=count
	c.Data["blogs"]=blogs
	c.TplName = "admin_blog.html"
}


func (c *MainController) ShowDetail() {
	//文章内容
	id,err:=c.GetInt("id")
	if err!=nil{
		beego.Info("get id err")
		return
	}
	o:=orm.NewOrm()
	detail:=models.Blog{Id:id}
	detail_err:=o.Read(&detail)
	if detail_err!=nil{
		beego.Info(err)
	}
	looknum:=detail.LookNum

	detail.LookNum = looknum+1
	_,err2 := o.Update(&detail);
	if err2!=nil {
		beego.Info(err)
		return
	}

	c.Data["detail"]=detail

	//下方推荐内容
	var blogs []models.Blog
	//使用QueryBuilder查询
	qb1,err:=orm.NewQueryBuilder("mysql")
	qb1.Select("id","img","title","category","create_time").From("blog").OrderBy("id").Desc().Limit(3)
	sql1:=qb1.String()
	o.Raw(sql1).QueryRows(&blogs)
	c.Data["blogs"]=blogs


	//评论内容
	var talks []models.Talk
	//使用QueryBuilder查询
	qb2,err:=orm.NewQueryBuilder("mysql")
	qb2.Select("name","talk","create_time").From("talk").Where("blog_id=?")
	sql2:=qb2.String()
	o.Raw(sql2,id).QueryRows(&talks)
	c.Data["talks"]=talks
	c.Data["talks_count"]=len(talks)

	c.TplName = "detail.html"
}


func (c *MainController) DoAddTalk() {
	id,err:=c.GetInt("id")
	if err!=nil{
		beego.Info("get id err")
		return 
	}
	name:=c.GetString("name")
	content:=c.GetString("talk")

	o := orm.NewOrm()
	talk:=models.Talk{}
	talk.Name=name
	talk.Talk=content
	talk.BlogId=id
	_ ,err1 := o.Insert(&talk)
	if err1!=nil{
		beego.Info("评论失败")
		return 
	}

	detail:=models.Blog{Id:id}
	detail_err:=o.Read(&detail)
	if detail_err!=nil{
		beego.Info(err)
	}
	talknum:=detail.TalkNum

	detail.TalkNum = talknum+1
	_,err2 := o.Update(&detail);
	if err2!=nil {
		beego.Info(err)
		return
	}

	c.Redirect("/detail?id="+strconv.Itoa(id),302)
}


func (c *MainController) DoAddLike() {
	id,err:=c.GetInt("id")
	if err!=nil{
		beego.Info("get id err")
		return 
	}
	o := orm.NewOrm()
	detail:=models.Blog{Id:id}
	detail_err:=o.Read(&detail)
	if detail_err!=nil{
		beego.Info(err)
	}
	likenum:=detail.LikeNum

	detail.LikeNum = likenum+1
	_,err2 := o.Update(&detail);
	if err2!=nil {
		beego.Info(err)
		return
	}

	c.Redirect("/detail?id="+strconv.Itoa(id),302)
}


func (c *MainController) ShowBlogEdit() {
	id,err:=c.GetInt("id")
	if err!=nil{
		beego.Info("get id err")
		return
	}
	o:=orm.NewOrm()
	detail:=models.Blog{Id:id}
	detail_err:=o.Read(&detail)
	if detail_err!=nil{
		beego.Info(err)
	}


	var categories []models.Category
	_,err1:=o.QueryTable("Category").All(&categories)
	if err1!=nil{
		beego.Info("select err")
		return 
	}
	c.Data["categories"]=categories
	beego.Info(categories)
	c.Data["checkcategory"]=strings.Split(detail.Category, ",")


	c.Data["detail"]=detail
	c.TplName = "admin_blog_edit.html"
}

func (c *MainController) DoEdit() {

	title:=c.GetString("title")
	creater:=c.GetString("creater")
	remark:=c.GetString("remark")
	category:=c.GetStrings("categorys")
	content:=c.GetString("editorValue")

	if title=="" || creater=="" || remark==""{
		c.Ctx.WriteString("请全部填写完毕再提交哦！")
		return 
	}


	id,err:=c.GetInt("id")
	if err!=nil{
		beego.Info("get id err")
		return 
	}
	o := orm.NewOrm()
	detail:=models.Blog{Id:id}
	detail_err:=o.Read(&detail)
	if detail_err!=nil{
		beego.Info(err)
	}



	//上传图片文件,f用不上，直接关闭
	f,h,err := c.GetFile("img")

	if err==nil{
		defer f.Close()
	}
	
	imagePath:=detail.Img

	if err==nil{
		//限定格式
		fileext:=path.Ext(h.Filename)
		if fileext!=".jpg" && fileext!=".png" {
			beego.Info("upload fail")
			return
		}
		//限定大小
		if h.Size < 100000 {
			beego.Info("upload file size err")
		}
		//当前时间戳（纳秒）
		now := time.Now().UnixNano()

		imagePath="./static/upload/"+strconv.FormatInt(now,10)+h.Filename 
		c.SaveToFile("img",imagePath)
	}


	detail.Img=imagePath
	detail.Title=title
	detail.Remark=remark
	detail.Creater=creater
	detail.Category=strings.Join(category,",")  //将切片分隔成字符串
	detail.Content=content

	_ ,errUpdate := o.Update(&detail)
	if errUpdate!=nil{
		beego.Info("更新失败")
		c.Redirect("/edit?id="+strconv.Itoa(id),302)
		return 
	}
	c.Redirect("/bloglist",302)


}



func (c *MainController) DoDel() {
	id,id_err:=c.GetInt("id")
	if id_err!=nil{
		beego.Info("get id err")
		return 
	}
	o := orm.NewOrm()
	detail:=models.Blog{Id:id}
	detail_err:=o.Read(&detail)
	if detail_err!=nil{
		beego.Info(detail_err)
		return
	}
	_,err:=o.Delete(&detail)
	if err!=nil{
		beego.Info(err)
		return
	}

	c.Redirect("/bloglist",302)
}