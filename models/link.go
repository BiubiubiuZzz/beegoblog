package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//友情链接
type Link struct {
	Id int //主键
	Sitename string `orm:"size(80)"`//网站名称
	Url string `orm:"size(200)"`//网址
	Rank int //排序值
}

//获取ling结构体在数据库中对应的表名
func (link *Link) TableName() string {
	//在配置文件中获取表的前缀
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "link"
}

//插入
func (link *Link) Insert() error {
	if _, err := orm.NewOrm().Insert(link); err != nil {
		return err
	}
	return nil
}

//在tb_link表中读取指定字段的信息
func (link *Link) Read(fields ...string) error {
	if err := orm.NewOrm().Read(link, fields...); err != nil {
		return err
	}
	return nil
}

//更新tb_link表中指定的字段
func (link *Link) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(link, fields...); err != nil {
		return err
	}
	return nil
}

//删除
func (link *Link) Delete() error {
	if _, err := orm.NewOrm().Delete(link); err != nil {
		return err
	}
	return nil
}
