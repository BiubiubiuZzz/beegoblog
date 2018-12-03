package admin

import (
	"strings"
	"github.com/astaxie/beego/validation"
	"beego_blog/models"
	"github.com/astaxie/beego/orm"
)

type UserController struct {
	baseController
}

func (this *UserController) Add() {
	input := make(map[string]string)
	errmsg := make(map[string]string)
	if this.Ctx.Request.Method == "POST" {
		username := strings.TrimSpace(this.GetString("username"))
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")
		input["username"] = username
		input["password"] = password
		input["password2"] = password2
		input["email"] = email
		valid := validation.Validation{}
		if result := valid.Required(username, "username"); !result.Ok {
			errmsg["username"] = "请输入用户名!"
		}else if result := valid.MaxSize(username, 15, "username"); !result.Ok {
			errmsg["username"] = "用户名长度不能大于15个字符!"
		}
		if result := valid.Required(password, "password"); !result.Ok {
			errmsg["password"] = "请输入密码!"
		}
		if result := valid.Required(password2, "password2"); !result.Ok {
			errmsg["password2"] = "请再次输入密码!"
		}else if password != password2 {
			errmsg["password2"] = "两次输入的密码不一致!"
		}
		if result := valid.Required(email, "email"); !result.Ok {
			errmsg["email"] = "请输入email地址!"
		}else if result := valid.Email(email, "email"); !result.Ok {
			errmsg["email"] = "Email无效!"
		}

		if active > 0 {
			active = 1
		}else {
			active = 0
		}
		if len(errmsg) == 0 {
			var user = &models.User{}
			user.Username = username
			user.Password = models.Md5([]byte(password))
			user.Email = email
			user.Active = active
			if err := user.Insert(); err != nil {
				this.showmsg(err.Error())
			}
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["input"] = input
	this.Data["errmsg"] = errmsg
	this.display()
}

func (this *UserController) List() {
	var list []*models.User
	query := orm.NewOrm().QueryTable(new(models.User))
	count, _ := query.Count()
	query.OrderBy("-id").Limit(this.pager.Pagesize, (this.pager.Page-1)*this.pager.Pagesize).All(&list)
	this.Data["list"] = list
	this.pager.SetTotalnum(int(count))
	this.pager.SetUrlpath("/admin/user/list?page=%d")
	this.Data["pagebar"] = this.pager.ToString()
	this.display()
}


func (this *UserController) Edit() {
	id, _ := this.GetInt("id")
	user := &models.User{Id:id}
	if err := user.Read(); err != nil {
		this.showmsg("用户不存在!")
	}
	errmsg := make(map[string]string)
	if this.Ctx.Input.Method() == "POST" {
		password := strings.TrimSpace(this.GetString("password"))
		password2 := strings.TrimSpace(this.GetString("password2"))
		email := strings.TrimSpace(this.GetString("email"))
		active, _ := this.GetInt("active")
		valid := validation.Validation{}
		if password != "" {
			if request := valid.Required(password2, "password2"); !request.Ok {
				errmsg["password2"] = "请再次输入密码!"
			}else if password != password2 {
				errmsg["password2"] = "两次输入的密码不一致!"
			}else {
				user.Password = models.Md5([]byte(password))
			}
		}
		if result := valid.Required(email, "email"); !result.Ok {
			errmsg["email"] = "请输入Email地址!"
		}else if result := valid.Email(email, "email"); !result.Ok {
			errmsg["email"] = "Email无效!"
		}else {
			user.Email = email
		}
		if active > 0 {
			user.Active = 1
		}else {
			user.Active = 0
		}
		if len(errmsg) == 0 {
			user.Update()
			this.Redirect("/admin/user/list", 302)
		}
	}
	this.Data["user"] = user
	this.Data["errmsg"] = errmsg
	this.display()
}

func (this *UserController) Delete() {
	id, _ := this.GetInt("id")
	if id == 7 {
		this.showmsg("不能删除ID为7的用户!")
	}
	user := &models.User{Id:id}
	if user.Read() == nil {
		user.Delete()
	}
	this.Redirect("/admin/user/list", 302)
}

