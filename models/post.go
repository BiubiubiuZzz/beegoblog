package models

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
	"strings"
)

//文章
type Post struct {
	//主键
	Id int
	//用户id
	Userid int
	//作者
	Author string `orm:"size(15)"`
	//标题
	Title string `orm:"size(100)"`
	//标题颜色
	Color string `orm:"size(7)"`
	//文章内容
	Content string `orm:"type(text)"`
	//标签
	Tags string `orm:"size(100)"`
	//访问量
	Views int
	//文章状态
	Status int
	//发表时间
	Posttime time.Time `orm:"type(datetime)"`
	//更新时间
	Updated time.Time `orm:"type(datetime)"`
	//是否置顶
	Istop int
	//封面
	Cover string `orm:"size(70)"`
}

//获取post结构体在数据库中对应的表名
func (post *Post) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	fmt.Println("文章表对应的名称 = ", dbprefix + "post")
	return dbprefix + "post"//tb_post
}

//插入
func (post *Post) Insert() error {
	if _, err := orm.NewOrm().Insert(post); err != nil {
		return err
	}
	return nil
}

//在tb_post表中 读取指定字段的信息
func (post *Post) Read(fields ...string) error {
	if err := orm.NewOrm().Read(post, fields...); err != nil {
		return err
	}
	return nil
}

//在tb_post表中更新指定字段信息
func (post *Post) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(post, fields...); err != nil {
		return err
	}
	return nil
}

//删除文章，需要删除标签文章表中对应的记录，且需要更新标签表中的count字段
func (post *Post) Delete() error {
	//判断文章的标签是否为空
	if post.Tags != ""{
		//创建orm对象
		orm := orm.NewOrm()
		//获得tb_tagpost表的句柄，并设置过滤条件:postid=post.Id
		query := orm.QueryTable(&TagPost{}).Filter("postid", post.Id)
		//创建文章标签指针结构体，用于存储查询结果
		var tagpost []*TagPost
		//查询结果，并判断查询到的记录数是否大于0且是否有错误
		if n, err := query.All(&tagpost); n > 0 && err == nil {
			//遍历查询结果
			for i := 0; i < len(tagpost); i++ {
				//创建标签，并初始化id
				var tag = &Tag{Id:tagpost[i].Tagid}
				//根据标签id查询标签
				if err = tag.Read(); err == nil && tag.Count > 0 {
					//count减一
					tag.Count--
					//更新count字段
					tag.Update("count")
				}
			}
		}
		//删除文章标签表中对应的记录
		orm.QueryTable(&TagPost{}).Filter("postid", post.Id).Delete()
	}
	//删除文章
	if _, err := orm.NewOrm().Delete(post); err != nil {
		return err
	}
	return nil
}



func (this *Post) Excerpt() string {
	return this.Content
}

func (this *Post) Link() string {
	//fmt.Println("Link被调用了")
	return "/article/" + strconv.Itoa(this.Id)//article/
}

func (this *Post) ColorTitle() string {
	if this.Color != "" {
		return fmt.Sprintf("<span style=\"color:%s\">%s</span>", this.Color, this.Title)
	}
	return this.Title
}

func (this *Post) TagsLink() string {
	if this.Tags == "" {
		return ""
	}
	tagslink := strings.Trim(this.Tags, ",")
	return tagslink
}


//获取上一篇文章(文章id小于当前文章的id)和下一篇文章(文章id大于当前文章的id)
func (this *Post) GetPreAndNext() (pre, next *Post) {
	//创建文章结构体
	pre = &Post{}
	next = &Post{}
	//获得tb_post表的结构体，并根据指定的过滤条件(文章id小于当前文章的id,状态正常)查询一条记录
	//SELECT * FROM tb_post WHERE id = num - 1; num为当前文章的id
	err := orm.NewOrm().QueryTable(new(Post)).Filter("id__lt", this.Id).Filter("status", 0).Limit(1).One(pre)
	//处理错误
	if err != nil {
		pre = nil
	}
	//获得tb_post表的结构体，并根据指定的过滤条件(文章id大于当前文章的id,状态正常)查询一条记录
	//SELECT * FROM tb_post WHERE id = num + 1; num为当前文章的id
	err = orm.NewOrm().QueryTable(new(Post)).Filter("id__gt", this.Id).Filter("status", 0).Limit(1).One(next)
	//处理错误
	if err != nil {
		pre = nil
	}
	return
}
