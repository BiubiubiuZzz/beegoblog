package admin

import (
	"strings"
	"beego_blog/models"
	"time"
	"github.com/astaxie/beego/orm"
	"fmt"
	"math/rand"
)

type MoodController struct {
	baseController
}

func (this *MoodController) Add() {
	if this.Ctx.Request.Method == "POST" {
		content := strings.TrimSpace(this.GetString("content"))
		var mood models.Mood
		mood.Content = content
		rand.Seed(time.Now().Unix())
		var r = rand.Intn(11)
		mood.Cover = "/static/upload/blog" + fmt.Sprintf("%d", r) + ".jpg"
		mood.Posttime = time.Now()
		if err := mood.Insert(); err != nil {
			this.showmsg(err.Error())
		}
		this.Redirect("/admin/mood/list", 302)
	}
	this.display()
}

//说说列表
func (this  *MoodController) List() {
	this.display()
	var list []*models.Mood
	query := orm.NewOrm().QueryTable(new(models.Mood))
	count, _ := query.Count()
	if count > 0 {
		query.OrderBy("-id").Limit(this.pager.Pagesize, (this.pager.Page-1)*this.pager.Pagesize).All(&list)
	}
	this.Data["list"] = list
	this.pager.SetUrlpath("/admin/mood/list?page=%d")
	this.pager.SetTotalnum(int(count))
	this.Data["pagebar"] = this.pager.ToString()
}

func (this *MoodController) Delete() {
	id, err := this.GetInt("id")
	if err != nil {
		this.showmsg("删除失败!")
	}
	mood := models.Mood{Id:id}
	if err = mood.Read(); err == nil {
		mood.Delete()
	}
	this.Redirect("/admin/mood/list", 302)
}
















