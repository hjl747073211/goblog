package models
import(
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

//用户
type User struct{
	Id int
	Name string `orm:"unique"`
	Pwd string
}


//分类
type Category struct {
	Id int
	Category string
	Remark string
}

//封面图
type Image struct {
	Id int 
	Image string
	Title string
	Content string
	UpdateTime time.Time `orm:"auto_now"`
}

//博客
type Blog struct {
	Id int
	Img string
	Title string `orm:"size(128);unique" valid:"Required"`
	Remark string
	Content string `orm:"type(text)"`
	Creater string 
	Category string
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
	LookNum int
	TalkNum int
	LikeNum int
}

//评论
type Talk struct {
	Id int 
	Name string
	Talk string
	CreateTime time.Time `orm:"auto_now"`
	BlogId int
}

func init(){
	//设置数据库基本信息,default一般不用,test是库名
	orm.RegisterDataBase("default","mysql","root:root@tcp(127.0.0.1:3306)/test?charset=utf8&loc=Local")
	//映射model数据
	orm.RegisterModel(new(User),new(Article),new(Category),new(Image),new(Blog),new(Talk) )//生成user表
	//生成表，fasle,是否更新(不要选true，会把表清空)，true是否可见
	orm.RunSyncdb("default",false,true)
}