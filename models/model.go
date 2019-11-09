package models
import(
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type User struct{
	Id int
	Name string `orm:"unique"`
	Pwd string
}

type Article struct {
	Id int
	Aname string `orm:"size(20)"`  //设置长度
	Atime time.Time
	Account int `orm:"default(0)"` //设置默认值
	Acontent string
	Aimg string
}

type Category struct {
	Id int
	Category string
	Remark string
}

type Image struct {
	Id int 
	Image string
	Title string
	Content string
	UpdateTime time.Time `orm:"auto_now"`
}

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