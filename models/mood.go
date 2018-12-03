package models

import (
	"time"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//说说
type Mood struct {
	Id int//主键
	Content string `orm:"type(text)"`//说说内容
	Cover string `orm:"size(70)"`//封面
	Posttime time.Time `orm:type(datetime)`//发表时间
}

//获取mood结构体在数据库中对应的表明
func (mood *Mood) TableName() string {
	//在配置文件中获取表的前缀
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "mood"
}

//插入
func (mood *Mood) Insert() error {
	if _, err := orm.NewOrm().Insert(mood); err != nil {
		return err
	}
	return nil
}

//在tb_mood表中读取指定字段的信息
func (mood *Mood) Read(fields ...string) error {
	if err := orm.NewOrm().Read(mood, fields...); err != nil {
		return err
	}
	return nil
}

//更新tb_mood表中指定字段
func (mood *Mood) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(mood, fields...); err != nil {
		return err
	}
	return nil
}

//删除
func (mood *Mood) Delete() error {
	if _, err := orm.NewOrm().Delete(mood); err != nil {
		return err
	}
	return nil
}



