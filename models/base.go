package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"crypto/md5"
	"fmt"
)

func init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	dburl := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	// set default database
	orm.RegisterDataBase("default", "mysql", dburl, 30)
	// register model
	orm.RegisterModel(new(Link), new(Mood),  new(Post), new(Tag), new(TagPost), new(User))
}


//从缓存中获取最新的4篇文章
func GetLatestBlog() []*Post {

	//创建文章指针切片
	var result []*Post
	post := Post{}
	//从tb_post表中过滤状为0，能通过id进行访问的文章
	query := orm.NewOrm().QueryTable(&post).Filter("status", 0)
	//查询数量
	count, _ := query.Count()
	//判断数量是否大于0
	if count > 0 {
		//如果数量大于0，则通过发表时间倒序查找4片最新文章
		query.OrderBy("-posttime").Limit(4).All(&result)
	}

	return result
}

//从tb_post表中获取最热门的4片文章
func GetHotBlog() []*Post {

	//创建文章指针切片
	var result []*Post
	post := Post{}
	//从tb_post表中过滤状为0，能通过id进行访问的文章
	query := orm.NewOrm().QueryTable(&post).Filter("status", 0)
	//查询数量
	count, _ := query.Count()
	//判断数量是否大于0
	if count > 0 {
		//按点击量倒序排序，查询出4篇文章
		query.OrderBy("-views").Limit(4).All(&result)
	}
	return result
}

//从tb_link表中获取友情链接
func GetLinks() []*Link {
	//判断缓存中是否存在links，如果不存在从数据库中获取，否则从缓存中获取
		//创建友情链接切片
		var result []*Link
		link := Link{}
		//得到tb_link表的句柄
		query := orm.NewOrm().QueryTable(&link)
		//得到友情链接的数量
		count, _ := query.Count()
		//判断数量是否大于0
		if count > 0 {
			//根据排序等级倒序排列
			query.OrderBy("-rank").All(&result)
		}
		//存入缓存

	return result
}


func Md5(buf []byte) string {
	mymd5 := md5.New()
	mymd5.Write(buf)
	result := mymd5.Sum(nil)
	return fmt.Sprintf("%x", result)
}











