package admin

import (
	"beego_blog/models"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
	"fmt"
	"strconv"
	"math/rand"
)

type ArticleController struct {
	baseController
}

func (this *ArticleController) List() {
	searchtype := this.GetString("searchtype")
	keyword := this.GetString("keyword")
	status, _ := this.GetInt("status")
	var list []*models.Post
	query := orm.NewOrm().QueryTable(new(models.Post)).Filter("status", status)
	if keyword != "" {
		switch searchtype {
		case "title":
			query = query.Filter("title__icontains", keyword)
		case "author":
			query = query.Filter("author__icontains", keyword)
		case "tag":
			query = query.Filter("tags__icontains", keyword)
		}
	}

	count, _ := query.Count()
	if count > 0 {
		query.Limit(this.pager.Pagesize, (this.pager.Page-1)*this.pager.Pagesize).All(&list)
	}
	this.Data["status"] = status
	this.Data["count_1"], _ = orm.NewOrm().QueryTable(&models.Post{}).Filter("status", 1).Count()
	this.Data["count_2"], _ = orm.NewOrm().QueryTable(&models.Post{}).Filter("status", 2).Count()
	this.Data["searchtype"] = searchtype
	this.Data["keyword"] = keyword
	this.Data["list"] = list
	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath(fmt.Sprintf("/admin/article/list?status=%d&searchtype=%s&keyword=%s&page=%s", status, searchtype, keyword, "%d"))
	this.Data["pagebar"] = this.pager.ToString()
	this.display()
}

func (this *ArticleController) Add() {
	this.display()
}


func (this *ArticleController) Save() {
	var post models.Post
	post.Title = strings.TrimSpace(this.GetString("title"))
	if post.Title == "" {
		this.showmsg("标题不能为空!")
	}
	post.Color = strings.TrimSpace(this.GetString("color"))
	if strings.TrimSpace(this.GetString("istop")) == "1" {
		post.Istop = 1
	}
	tags := strings.TrimSpace(this.GetString("tags"))
	timestr := strings.TrimSpace(this.GetString("posttime"))
	post.Status, _ = this.GetInt("status")
	post.Content = this.GetString("content")
	post.Userid = this.userid
	post.Author = this.username
	post.Updated = this.getTime()
	rand.Seed(time.Now().Unix())
	var r = rand.Intn(11)
	post.Cover = "/static/upload/blog" + fmt.Sprintf("%d", r) + ".jpg"
	posttime, err := time.Parse("2006-01-02 15:04:05", timestr)
	if err == nil {
		post.Posttime = posttime
	}
	if err1 := post.Insert(); err1 != nil {
		fmt.Println("err = ", err1)
	}

	addtags := make([]string, 0)
	if tags  != "" {
		tagarr := strings.Split(tags, ",")
		for _, v := range tagarr {
			if tag := strings.TrimSpace(v); tags != "" {
				exists := false
				for _, vv := range addtags {
					if vv == tag {
						exists = true
						break
					}
				}
				if !exists {
					addtags = append(addtags, tag)
				}
			}
		}
	}
	if len(addtags) > 0 {
		for _, v := range addtags {
			tag := &models.Tag{Name :v}
			if err := tag.Read("Name"); err == orm.ErrNoRows {
				tag.Count = 1
				tag.Insert()
			}else {
				tag.Count += 1
				tag.Update("Count")
			}
			tp := &models.TagPost{Tagid:tag.Id, Postid:post.Id, Poststatus:post.Status, Posttime:this.getTime()}
			tp.Insert()
		}
		post.Tags = "," + strings.Join(addtags, ",") + ","//strings.Join(addtags, "+")  [Go,Python,C++,Java]  ==>,Go,Python,C++,Java,
	}
	post.Updated = this.getTime()
	post.Update("tags", "updated")
	this.Redirect("/admin/article/list", 302)
}

func (this *ArticleController) Edit() {
	id, _ := this.GetInt("id")
	post := &models.Post{Id:id}
	if post.Read() != nil {
		this.Abort("404")
	}
	post.Tags = strings.Trim(post.Tags, ",")
	this.Data["post"] = post
	this.Data["posttime"] = post.Posttime.Format("2006-01-02 15:04:05")
	this.display()
}


func (this *ArticleController) Update() {
	var post models.Post
	id, err := this.GetInt("id")
	if err == nil {
		post.Id = id
		if post.Read() != nil {
			this.Redirect("/admin/article/list", 302)
		}
	}
	post.Title = strings.TrimSpace(this.GetString("title"))
	if post.Title == "" {
		this.showmsg("标题不能为空!")
	}
	post.Color = strings.TrimSpace(this.GetString("color"))
	if strings.TrimSpace(this.GetString("istop")) == "1" {
		post.Istop = 1
	}
	tags := strings.TrimSpace(this.GetString("tags"))
	timestr := strings.TrimSpace(this.GetString("posttime"))
	post.Status, _ = this.GetInt("status")
	post.Content = this.GetString("content")
	post.Updated = this.getTime()
	if posttime, err := time.Parse("2006-01-02 15:04:05", timestr); err == nil {
		post.Posttime = posttime
	}
	if strings.Trim(post.Tags, ",") == tags {
		post.Update("title", "color", "istop", "status", "content", "posttime")
		this.Redirect("/admin/article/list", 302)
	}

	if post.Tags != "" {
		var tagpost models.TagPost
		query := orm.NewOrm().QueryTable(&tagpost).Filter("postid", post.Id)
		var tagpostarr []*models.TagPost
		if n, err := query.All(&tagpostarr); n > 0 && err == nil {
			for i := 0; i < len(tagpostarr); i++ {
				var tag = &models.Tag{Id:tagpostarr[i].Tagid}
				if err = tag.Read(); err == nil && tag.Count > 0 {
					tag.Count--
					tag.Update("count")
				}
			}
		}
		query.Delete()
	}
	addtags := make([]string, 0)
	if tags != "" {
		tagarr := strings.Split(tags, ",")
		for _, v := range tagarr {
			if tag := strings.TrimSpace(v); tag != "" {
				exists := false
				for _, vv := range addtags {
					if vv == tag {
						exists = true
						break
					}
				}
				if !exists{
					addtags = append(addtags, tag)
				}
			}
		}
	}
	if len(addtags) > 0 {
		for _, v := range addtags {
			tag := &models.Tag{Name :v}
			if err := tag.Read("Name"); err == orm.ErrNoRows {
				tag.Count = 1
				tag.Insert()
			}else {
				tag.Count += 1
				tag.Update("Count")
			}
			tp := &models.TagPost{Tagid:tag.Id, Postid:post.Id, Poststatus:post.Status, Posttime:this.getTime()}
			tp.Insert()
		}
		post.Tags = "," + strings.Join(addtags, ",") + ","
	}
	post.Update("title", "color", "istop", "tags", "status", "content", "updated", "posttime")
	this.Redirect("/admin/article/list", 302)
}


func (this *ArticleController) Delete() {
	id, _ := this.GetInt("id")

	post := &models.Post{Id:id}
	if post.Read() == nil {
		post.Delete()
	}
	this.Redirect("/admin/article/list", 302)
}


func (this *ArticleController) Batch() {
	ids := this.GetStrings("ids[]")
	op := this.GetString("op")
	idarr := make([]int, 0)
	for _, v := range ids {
		if id, _ := strconv.Atoi(v); id > 0 {
			idarr = append(idarr, id)
		}
	}
	var post models.Post
	query := orm.NewOrm().QueryTable(&post)
	switch op {
	case "topub":
		query.Filter("id__in", idarr).Update(orm.Params{"status":0})
	case "todrafts":
		query.Filter("id__in", idarr).Update(orm.Params{"status":1})
	case "totrash":
		query.Filter("id__in", idarr).Update(orm.Params{"status":2})
	case "delete":
		for _, id := range idarr {
			obj := models.Post{Id:id}
			if obj.Read() == nil {
				obj.Delete()
			}
		}
	}
	this.Redirect(this.Ctx.Request.Referer(), 302)
}










