package Controllers

import (
	"Model"
	"System/Web"
	"fmt"
)

type Home struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(Home{})
}
func (this *Home) OnLoad() {
	//如果在OnLoad里调用了Response.Write需选设置Content-Type头，设置文档类型
	//this.Response.Header().Add("Content-Type", "text/html;charset=utf-8")
	//this.Response.Write([]byte("在 OnLoad函数里面"))
}

func (this *Home) Index() *Web.ViewResult {
	this.ViewData["Title"] = "欢迎使用GoMvc"
	this.Session["aaa"] = "aaaaa"
	return this.View()
}
func (this *Home) Config() *Web.ViewResult {
	this.ViewData["Title"] = "配置文件"
	this.ViewData["aa"] = this.Session["aaa"]
	return this.View()
}
func (this *Home) Route() *Web.ViewResult {
	this.ViewData["Title"] = "路由"
	return this.View()
}
func (this *Home) TemplateFunc() *Web.ViewResult {
	this.ViewData["Title"] = "模板函数"

	return this.View()
}
func (this *Home) Binder(u Model.User) *Web.ViewResult {
	this.ViewData["Title"] = "参数绑定"
	if this.Request.Method == "POST" {
		this.ViewData["ShowLogin"] = false
		this.ViewData["User"] = fmt.Sprintf("%v", u)
	} else {
		this.ViewData["ShowLogin"] = true
	}
	return this.View()
}
func (this *Home) Footer() *Web.ViewResult {
	this.ViewData["Copyright"] = "© Company 2013"
	return this.View()
}
func (this *Home) TestScript() *Web.JavaScriptResult {
	return this.JavaScript("alert('OK!');", "utf-8")
}
func (this *Home) TestJson() *Web.JsonResult {
	this.ViewData["UserName"] = "张三"
	this.ViewData["AGe"] = "30"
	return this.Json(this.ViewData, "utf-8")
}

type User struct {
	UserName string
	Age      int
}

func (this *Home) TestXml() *Web.XmlResult {
	u := User{"张三", 19}
	return this.Xml(u, "utf-8")
}
func (this *Home) Action() *Web.ViewResult {
	return this.View()
}
