package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

type Tag struct {
	//主键id
	Id int
	//标签名称
	Name string `orm:"size(20)"`
	//和某个标签相关的文章数量
	Count int
}

//在数据库中获取Tag结构体所对应的表名
func (tag *Tag) TableName() string {
	dbprefix := beego.AppConfig.String("dbprefix")
	return dbprefix + "tag"
}

//插入
func (tag *Tag) Insert() error {
	if _, err := orm.NewOrm().Insert(tag); err != nil {
		return err
	}
	return nil
}

//在数据库中根据指定字段读取相关信息
func (tag *Tag) Read(fields ...string) error {
	if err := orm.NewOrm().Read(tag, fields...); err != nil {
		return err
	}
	return nil
}

//根据指定字段更新tb_tag表
func (tag *Tag) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(tag, fields...); err != nil {
		return err
	}
	return nil
}

//删除
/*
思路：根据需要被删除的标签在标签文章表中查找相关记录，
从这些记录中去除对应的文章id，然后根据文章id在文章表中找到
对应的记录，将标签名替换为逗号，最后删除标签文章表中相关记录。
*/
func (tag *Tag) Delete() error {
	//创建TagPost指针切片
	var list []*TagPost
	//获取Post结构体在数据库中所对应表的表名
	table := new(Post).TableName()
	//在tb_tagpost表中查询出和tagid相关的记录
	orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).All(&list)
	if len(list) > 0 {
		//创建切片，用于存储文章id
		ids := make([]string, 0, len(list))
		for _, tagpost := range list {
			ids = append(ids, strconv.Itoa(tagpost.Postid))// 1 2 3==1,2,3
		}
		//update tb_post set tags = REPLACE(tags, ?, ',') where id in(1, 2, 3) tag.Name ,go,java, ==>,java,
		orm.NewOrm().Raw("update " + table + " set tags = REPLACE(tags, ?, ',') where id in (" + strings.Join(ids, ",")+")", ","+tag.Name+",").Exec()
		orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).Delete()
	}
	//删除标签
	if _, err := orm.NewOrm().Delete(tag); err != nil {
		return err
	}
	return nil
}

//更新tb_tag表中的count字段
func (tag *Tag) UpCount() {
	//获得标签文章表的句柄并设置过滤条件：tagid=tag.id,获取符合条件的记录的数量
	count, err := orm.NewOrm().QueryTable(&TagPost{}).Filter("tagid", tag.Id).Count()
	//将数量转换为整形
	newcount := int(count)
	//如果没有错误并且和该标签相关的文章数量和数据库中的数量不等
	if err == nil && newcount != tag.Count{
		//设置count字段
		tag.Count = newcount
		//更新count字段
		tag.Update("count")
	}

}

//合并标签，在标签文章表中找到和需要被删除的标签相关的记录，
//将这些记录中的标签id替换为目标标签的id，在标签文章表中根据标签id找到
//对应的文章id，然后通过文章id找到对应的记录，将这些记录中的标签名称替换成
//目标标签的名称
func (tag *Tag) MergeTo(to *Tag) {
	//创建标签文章指针结构体，用于存储查询结果
	var list []*TagPost
	//创建标签文章结构体
	var tp TagPost
	//获取标签文章表的句柄
	query := orm.NewOrm().QueryTable(&tp)
	//设置过滤条件tagid=tag.Id,并查询所有和结果的记录
	query.Filter("tagid", tag.Id).All(&list)
	//判断查询到的记录数量是否大于0
	if len(list) > 0 {
		//创建切片，用于存储和该标签相关的文章id
		ids := make([]string, 0, len(list))
		//遍历查询结果
		for _, v := range list {
			//将标签文章中的文章id追加到ids中
			ids = append(ids, strconv.Itoa(v.Postid))
		}
		//将标签文章表中相关记录中的标签id更新为目标标签的名称
		query.Filter("tagid", tag.Id).Update(orm.Params{"tagid":to.Id})
		//将文章表中相关记录的标签的名称替换为目标标签的名称
		//update tb_post set tags = Replace(tags, tag.Name, to.Name) where id (ids)
		orm.NewOrm().Raw("Update " + new(Post).TableName() + " set tags = Replace(tags, ?, ?) where id in (" + strings.Join(ids, ",") +")", ","+tag.Name+",", ","+to.Name+",").Exec()
	}
}


