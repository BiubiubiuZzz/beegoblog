package admin

import (
	"github.com/astaxie/beego"
	"strings"
	"beego_blog/models"
	"strconv"

	"fmt"
	"time"
)

type baseController struct {
	beego.Controller
	userid int//用户id
	username string//用户姓名
	controllerName string //控制器的名称
	actionName string //处理函数的名称
	pager *models.Pager
}

func (this *baseController) Prepare() {
	controllerName, actionName := this.GetControllerAndAction()
	this.controllerName = strings.ToLower(controllerName[:len(controllerName)-10])
	this.actionName = strings.ToLower(actionName)

	this.auth()
	this.checkPermission()

	page, err := this.GetInt("page")
	if err != nil {
		page = 1
	}
	pagesize := 2

	this.pager = models.NewPager(page, pagesize, 0, "")
}

func (this *baseController) checkPermission() {
	if this.userid != 7 && this.controllerName == "user" {
		this.showmsg("抱歉,只有超级管理员才能进行该操作!")
	}
}


func (this *baseController) display(tplname ...string) {
	moduleName := "admin/"
	this.Layout = moduleName + "layout.html"
	this.Data["version"] = beego.AppConfig.String("version")
	this.Data["adminname"] = this.username
	if len(tplname) == 1 {
		this.TplName = moduleName + tplname[0] + ".html"
	}else {
		this.TplName = moduleName + this.controllerName + "_" + this.actionName + ".html"
	}
}


func (this *baseController) showmsg(msg ...string) {
	this.display("showmsg")
	if len(msg) == 1 {
		msg = append(msg, this.Ctx.Request.Referer())
	}

	this.Data["msg"] = msg[0]
	this.Data["redirect"] = msg[1]
	this.Render()
	this.StopRun()
}

func (this *baseController) auth() {
	if this.controllerName == "account" && this.actionName == "login" {

	}else {
		arr := strings.Split(this.Ctx.GetCookie("auth"), "|")
		if len(arr) == 2 {
			idstr, password := arr[0], arr[1]
			id, _ := strconv.Atoi(idstr)
			if id > 0 {
				user := new(models.User)
				user.Id = id
				if user.Read() == nil && password == user.Password {
					this.userid = user.Id
					this.username = user.Username
				}
			}
		}
		if this.userid == 0 {
			fmt.Println("即将重定向")
			this.Redirect("/admin/login", 302)
		}
	}
}


func (this *baseController) getTime() time.Time {
	return time.Now().UTC()
}




