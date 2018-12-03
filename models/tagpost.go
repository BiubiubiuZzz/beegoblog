package models

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type TagPost struct {
	//主键
	Id int
	//标签id
	Tagid int
	//文章id
	Postid int
	//文章状态
	Poststatus int
	//发表时间
	Posttime time.Time `orm:"type(datetime)"`
}

//获取TagPost结构体在数据库中对应表的表名
func (tagpost *TagPost) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "tag_post"
}

//插入
func (tagpost *TagPost) Insert() error {
	if _, err := orm.NewOrm().Insert(tagpost); err != nil {
		return err
	}
	return nil
}

//在tb_tag_post表中根据指定字段读取相关信息
func (tagpost *TagPost) Read(fields ...string) error {
	if err := orm.NewOrm().Read(tagpost, fields...); err != nil {
		return err
	}
	return nil
}

//根据自定字段更新tb_tag_post表中的信息
func (tagpost *TagPost) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(tagpost, fields...); err != nil {
		return err
	}
	return nil
}

//删除
func (tagpost *TagPost) Delete() error {
	if _, err := orm.NewOrm().Delete(tagpost); err != nil {
		return err
	}
	return nil
}





