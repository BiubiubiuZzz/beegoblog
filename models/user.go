package models

import (

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//用户
type User struct {
	//主键
	Id int
	//用户名
	Username string `orm:"size(15)"`
	//密码
	Password string `orm:"size(32)"`
	//邮箱
	Email string `orm:"size(50)"`
	//登陆次数
	Logincount int
	//用户key
	Authkey string `orm:"size(10)"`
	//是否激活
	Active int
}

//获取user结构体在数据库中所对应的表的表名
func (user *User) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "user"
}

//插入
func (user *User) Insert() error {
	if _, err := orm.NewOrm().Insert(user); err != nil {
		return err
	}
	return nil
}

//读取
func (user *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(user, fields...); err != nil {
		return err
	}
	return nil
}

//更新
func (user *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(user, fields...); err != nil {
		return err
	}
	return nil
}

//删除
func (user *User) Delete() error {
	if _, err := orm.NewOrm().Delete(user); err != nil {
		return err
	}
	return nil
}

