package routers

import (
	"github.com/astaxie/beego"
	"beego_blog/controllers"
	"beego_blog/controllers/admin"
)

func init() {
	//----------------------------前台页面------------------------------------------------------
	//首页
	beego.Router("/", &controllers.MainController{}, "*:Index")
	//首页分页路由(href="/index2.html")
	beego.Router("/index:page:int.html", &controllers.MainController{}, "*:Index")
	//文章详情(aticle/2)
	beego.Router("/article/:id:int", &controllers.MainController{}, "*:Show")
	//关于我
	beego.Router("/about.html", &controllers.MainController{}, "*:About")
	//成长录
	beego.Router("/life.html", &controllers.MainController{}, "*:BlogList")
	//成长录分页路由(/lift2.html)
	beego.Router("/life:page:int.html", &controllers.MainController{}, "*:BlogList")
	//碎言碎语
	beego.Router("/mood.html", &controllers.MainController{}, "*:Mood")
	//碎言碎语分页路由(/mood2.html)
	beego.Router("/mood:page:int.html", &controllers.MainController{}, "*:Mood")

	////错误页面，匹配任意url
	beego.Router("/:urlname(.+)", &controllers.MainController{}, "*:Go404")

	//----------------------------后台页面-----------------------------------------------------------
	//首页
	beego.Router("/admin", &admin.IndexController{}, "*:Index")

	//----------------------------说说管理-----------------------------------
	//添加说说
	beego.Router("/admin/mood/add", &admin.MoodController{}, "*:Add")
	//说说列表
	beego.Router("/admin/mood/list", &admin.MoodController{}, "*:List")
	//删除说说
	beego.Router("/admin/mood/delete", &admin.MoodController{}, "*:Delete")

	//----------------------------账户管理-----------------------------------
	//登录
	beego.Router("/admin/login", &admin.AccountController{}, "*:Login")
	//退出登陆
	beego.Router("/admin/logout", &admin.AccountController{}, "*:Logout")
	//修改个人资料
	beego.Router("/admin/account/profile", &admin.AccountController{}, "*:Profile")

	//-----------------------------友情链接管理------------------------------
	//添加友情链接
	beego.Router("/admin/link/add", &admin.LinkController{}, "*:Add")
	//友情链接列表
	beego.Router("/admin/link/list", &admin.LinkController{}, "*:List")
	//编辑友情链接
	beego.Router("/admin/link/edit", &admin.LinkController{}, "*:Edit")
	//删除友情链接
	beego.Router("/admin/link/delete", &admin.LinkController{}, "*:Delete")

	//-----------------------------用户管理-----------------------------------
	//添加用户
	beego.Router("/admin/user/add", &admin.UserController{}, "*:Add")
	//用户列表
	beego.Router("/admin/user/list", &admin.UserController{}, "*:List")//   /admin/user/list?page=3
	//编辑用户
	beego.Router("/admin/user/edit", &admin.UserController{}, "*:Edit")
	//删除用户
	beego.Router("/admin/user/delete", &admin.UserController{}, "*:Delete")

	//-----------------------------文章管理-------------------------------------
	//文章列表
	beego.Router("/admin/article/list", &admin.ArticleController{}, "*:List")
	//添加文章(跳转到添加文章的页面)
	beego.Router("/admin/article/add", &admin.ArticleController{}, "*:Add")
	//保存文章
	beego.Router("/admin/article/save", &admin.ArticleController{}, "*:Save")
	//编辑文章(跳转到文章编辑页面)
	beego.Router("/admin/article/edit", &admin.ArticleController{}, "*:Edit")
	//更新文章
	beego.Router("/admin/article/update", &admin.ArticleController{}, "*:Update")
	//删除文章
	beego.Router("/admin/article/delete", &admin.ArticleController{}, "*:Delete")
	//批量处理
	beego.Router("/admin/article/batch", &admin.ArticleController{}, "*:Batch")

	//-----------------------------标签管理--------------------------------------
	beego.Router("/admin/tag", &admin.TagController{}, "*:Index")

}
