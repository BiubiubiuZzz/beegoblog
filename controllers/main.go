package controllers

import (
	"github.com/astaxie/beego"
	"beego_blog/models"
	"github.com/astaxie/beego/orm"
	"fmt"
	"strconv"
)

type MainController struct {
	beego.Controller
	Pager *models.Pager
}

func (this *MainController) Prepare() {
	var page int
	var pagesize int
	var err error
	if page, err = strconv.Atoi(this.Ctx.Input.Param(":page")); err != nil {
		page = 1
	}
	pagesize = 2
	this.Pager = models.NewPager(page, pagesize, 0, "")

}


func (this *MainController) display(tplname string) {
	theme := "double"
	this.Layout =  theme + "/layout.html"
	this.TplName = theme + "/" + tplname + ".html"
	this.LayoutSections = make(map[string]string)
	this.LayoutSections["head"] = theme + "/head.html"
	this.LayoutSections["foot"] = theme + "/foot.html"
	if tplname == "index" {
		this.LayoutSections["banner"] = theme + "/banner.html"
		this.LayoutSections["middle"] = theme + "/middle.html"
		this.LayoutSections["right"] = theme + "/right.html"
	}else if tplname == "life" {//判断是不是成长录
		this.LayoutSections["right"] = theme + "/right.html"
	}

}


func (this *MainController) setHeadMeater() {
	this.Data["title"] = beego.AppConfig.String("title")
	this.Data["keywords"] = beego.AppConfig.String("keywords")
	this.Data["description"] = beego.AppConfig.String("description")
}

func (this *MainController) setRight() {
	this.Data["latestblog"] = models.GetLatestBlog()
	this.Data["hotblog"] = models.GetHotBlog()
	this.Data["links"] = models.GetLinks()
}

func (this *MainController) Index() {
	this.display("index")
	this.setHeadMeater()
	this.setRight()

	var list []*models.Post
	post := models.Post{}

	query := orm.NewOrm().QueryTable(&post).Filter("status", 0)
	count, _ := query.Count()
	this.Pager.SetTotalnum(int(count))
	this.Pager.SetUrlpath("/index%d.html")
	if count > 0 {
		_, err := query.OrderBy("-istop", "-views").Limit(this.Pager.Pagesize, (this.Pager.Page - 1)*this.Pager.Pagesize).All(&list)
		if err != nil {
			fmt.Println("err = ", err)
		}
	}
	this.Data["list"] = list
	this.Data["pagebar"] = this.Pager.ToString()
}

func (this *MainController) Show() {
	this.display("article")
	this.setHeadMeater()
	this.setRight()
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		this.Redirect("/404", 302)
	}
	post := new(models.Post)
	post.Id = id
	err = post.Read()
	if err != nil {
		this.Redirect("/404", 302)
	}
	post.Views++
	post.Update("Views")
	this.Data["post"] = post
	pre, next := post.GetPreAndNext()
	this.Data["pre"] = pre
	this.Data["next"] = next
	this.Data["smalltitle"] = "文章详情"
}


func (this *MainController) About() {
	this.display("about")
	this.setHeadMeater()
}


func (this *MainController) BlogList() {
	this.display("life")
	this.setHeadMeater()
	this.setRight()
	var list []*models.Post
	query := orm.NewOrm().QueryTable(new(models.Post)).Filter("status", 0)
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-istop", "-posttime").Limit(this.Pager.Pagesize, (this.Pager.Page-1)*this.Pager.Pagesize).All(&list)
	}
	this.Data["list"] = list
	this.Pager.SetUrlpath("/life%d.html")
	this.Pager.SetTotalnum(int(count))
	this.Data["pagebar"] = this.Pager.ToString()
}

func (this *MainController) Mood() {
	this.display("mood")
	this.setHeadMeater()
	var list []*models.Mood
	query := orm.NewOrm().QueryTable(new(models.Mood))
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-posttime").Limit(this.Pager.Pagesize, (this.Pager.Page - 1)*this.Pager.Pagesize).All(&list)
	}
	this.Data["list"] = list
	this.Pager.SetTotalnum(int(count))
	this.Pager.SetUrlpath("/mood%d.html")
	this.Data["pagebar"] = this.Pager.ToString()
}


func (this *MainController) Go404() {
	this.display("404")
	this.setHeadMeater()
}
















