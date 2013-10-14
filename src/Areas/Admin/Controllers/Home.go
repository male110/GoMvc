package Controllers

import (
	"System/Web"
	"net/http"
)

type Home struct {
	Web.Controller
}

//注册Controller
func init() {
	Web.App.RegisterController(Home{})
}
func (this *Home) OnLoad() {
	_, ok := this.Session["UserName"]
	strActionName := this.RouteData["action"]
	if !ok && strActionName != "Login" {
		http.Redirect(this.Response, this.Request, "/Admin/Home/Login", http.StatusFound)
		this.ResponseEnd()
	}
}

func (this *Home) Login() *Web.ViewResult {
	if this.Request.Method == "POST" {
		this.Session["UserName"] = this.Form["txtUserName"]

		http.Redirect(this.Response, this.Request, "/Admin/Home/Index", http.StatusFound)
		this.ResponseEnd()
		return this.View()
	} else {
		return this.View()
	}

}
func (this *Home) Top() *Web.ViewResult {
	this.ViewData["UserName"] = this.Session["UserName"]
	return this.View()
}
func (this *Home) Index() *Web.ViewResult {
	return this.View()
}
func (this *Home) Left() *Web.ViewResult {
	return this.View()
}
func (this *Home) Data() *Web.ViewResult {
	return this.View()
}
func (this *Home) Center() *Web.ViewResult {
	return this.View()
}
func (this *Home) Footer() *Web.ViewResult {
	return this.View()
}
func (this *Home) TestGlobal() *Web.ViewResult {
	return this.View()
}
