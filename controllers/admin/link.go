package admin

import (

	"beego_blog/models"
	"github.com/astaxie/beego/orm"
	"strings"
)

type LinkController struct {
	baseController
}

func (this *LinkController) Add() {
	if this.Ctx.Request.Method == "POST" {
		sitename := this.GetString("sitename")
		url := this.GetString("url")
		rank, err := this.GetInt("rank")
		if err != nil {
			rank = 0
		}
		var link = &models.Link{Sitename:sitename, Url:url, Rank:rank}
		if err = link.Insert(); err != nil {
			this.showmsg(err.Error())
		}
		this.Redirect("/admin/link/list", 302)
	}
	this.display()
}

func (this *LinkController) List() {
	this.display()
	var list []*models.Link
	orm.NewOrm().QueryTable(new(models.Link)).OrderBy("-rank").All(&list)
	this.Data["list"] = list
}

func (this *LinkController) Edit() {
	id, _ := this.GetInt("id")
	link := &models.Link{Id:id}
	if err := link.Read(); err != nil {
		this.showmsg("友链不存在!")
	}
	if this.Ctx.Request.Method == "POST" {
		sitename := strings.TrimSpace(this.GetString("sitename"))
		url := strings.TrimSpace(this.GetString("url"))
		rank, err := this.GetInt("rank")
		if err != nil {
			rank = 0
		}
		link.Sitename = sitename
		link.Url = url
		link.Rank = rank
		link.Update()
		this.Redirect("/admin/link/list", 302)
	}

	this.display()
	this.Data["link"] = link
}


func (this *LinkController) Delete() {
	id, err := this.GetInt("id")
	if err != nil {
		this.showmsg("删除失败!")
	}
	link := &models.Link{Id:id}
	if err = link.Read(); err == nil {
		link.Delete()
	}
	this.Redirect("/admin/link/list", 302)
}