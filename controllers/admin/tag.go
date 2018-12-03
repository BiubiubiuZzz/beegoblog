package admin

import (
	"strconv"
	"strings"
	"beego_blog/models"
	"github.com/astaxie/beego/orm"
)

type TagController struct {
	baseController
}


func (this *TagController) Index() {
	act := this.GetString("act")
	switch act {
	case "batch":
		this.batch()
	default:
		this.List()
	}
}


func (this *TagController) List() {
	var list []*models.Tag
	var tag models.Tag
	query := orm.NewOrm().QueryTable(&tag)
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-count").Limit(this.pager.Pagesize, (this.pager.Page-1)*this.pager.Pagesize).All(&list)
	}
	this.Data["list"] = list
	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath("/admin/tag?page=%d")
	this.Data["pagebar"] = this.pager.ToString()
	this.display("tag_list")
}

func (this *TagController) batch() {
	ids := this.GetStrings("ids[]")
	op := this.GetString("op")
	idarr := make([]int, 0)
	for _, v := range ids {
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}
	switch op {
	case "merge":
		toname := strings.TrimSpace(this.GetString("toname"))
		if toname != "" && len(idarr) > 0 {
			tag := new(models.Tag)
			tag.Name = toname
			if tag.Read("name") != nil {
				tag.Count = 0
				tag.Insert()
			}
			for _, id := range idarr {
				obj := models.Tag{Id:id}
				if obj.Read() == nil {
					obj.MergeTo(tag)
					obj.Delete()
				}
			}
			tag.UpCount()
		}
	case "delete":
		for _, id := range idarr {
			obj := models.Tag{Id:id}
			if obj.Read() == nil {
				obj.Delete()
			}
		}
	}
	this.Redirect("/admin/tag", 302)
}