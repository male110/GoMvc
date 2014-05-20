package Web

import (
	"System/ViewEngine"
	"fmt"
	"net/http"
)

/*定义一些通用的接口*/
type IController interface {
	SetResponse(rw http.ResponseWriter)
	SetRequest(request *http.Request)
	SetRouteData(routeData map[string]interface{})
	SetViewData(viewData map[string]interface{})
	SetSession(session map[string]interface{})
	IsEnd() bool //是否停止当前请求的执行过程
	SetTheme(string)
	SetViewEngin(ViewEngine.IViewEngine)
	SetCookies(map[string]string)
	SetBinder(binder *Binder)
	SetQueryString(map[string]string)
	SetForm(map[string]string)
	SetIsPost(bool)
}
type Controller struct {
	Request       *http.Request
	Response      http.ResponseWriter
	ViewData      map[string]interface{}
	Session       map[string]interface{}
	RouteData     map[string]interface{}
	QueryString   map[string]string
	Form          map[string]string
	ViewEngine    ViewEngine.IViewEngine
	IsEndResponse bool   //是否停止当前请求的执行
	Theme         string //主题名称
	Cookies       map[string]string
	DefaultBinder *Binder
	IsPost        bool //如果是Post提交的数据为true,否则为false
}

func (this *Controller) SetResponse(rw http.ResponseWriter) {
	this.Response = rw
}
func (this *Controller) SetRequest(request *http.Request) {
	this.Request = request
}
func (this *Controller) SetRouteData(routeData map[string]interface{}) {
	this.RouteData = routeData
}
func (this *Controller) SetViewData(viewData map[string]interface{}) {
	this.ViewData = viewData
}
func (this *Controller) SetSession(session map[string]interface{}) {
	this.Session = session
}
func (this *Controller) SetTheme(theme string) {
	this.Theme = theme
}
func (this *Controller) SetViewEngin(engine ViewEngine.IViewEngine) {
	this.ViewEngine = engine
}
func (this *Controller) IsEnd() bool {
	return this.IsEndResponse
}
func (this *Controller) SetCookies(cookie map[string]string) {
	this.Cookies = cookie
}
func (this *Controller) SetQueryString(queryString map[string]string) {
	this.QueryString = queryString
}
func (this *Controller) SetForm(form map[string]string) {
	this.Form = form
}
func (this *Controller) SetBinder(binder *Binder) {
	this.DefaultBinder = binder
}
func (this *Controller) ResponseEnd() {
	this.IsEndResponse = true
}
func (this *Controller) UpdateModel(data interface{}) error {
	return this.DefaultBinder.BindModel(data)
}
func (this *Controller) SetIsPost(isPost bool) {
	this.IsPost = isPost
}

/*该函数接受两个参数，第一个是要输出的javascript脚本内容,第二个参数是字符集，可变参数，可以省略不传，默认是utf-8编码*/
func (this *Controller) JavaScript(js string, charSet ...string) *JavaScriptResult {
	result := &JavaScriptResult{Script: js, Response: this.Response}
	if charSet != nil && len(charSet) > 0 {
		result.CharSet = charSet[0]
	}
	return result
}

/*该函数可以接受两个参数，第一个参数为json的内容，可以是一个Json字符串，也可以是一个对像，第二个为charSet字符集，第二个参数都可以省略，
第一个参数省略时或为""时，会把ViewData的内容转换为Json字符串，字符集省略时，默认为utf-8*/
func (this *Controller) Json(data interface{}, args ...string) *JsonResult {
	result := &JsonResult{Response: this.Response}
	strJson, ok := data.(string)
	if ok {
		result.JsonText = strJson
	} else {
		result.Data = data
	}
	if args != nil {
		if len(args) > 0 {
			result.CharSet = args[0]
		}
	}
	return result
}

/*该函数可以接受两个参数，第一个参数为xml的内容或struct结构体,不支持map，第二个为charSet字符集，第二个参数都可以省略，
第一个参数省略时或为""时，会把ViewData的内容转换为Json字符串，字符集省略时，默认为utf-8*/
func (this *Controller) Xml(data interface{}, charSet ...string) *XmlResult {
	result := &XmlResult{Response: this.Response}
	strXml, ok := data.(string)
	if ok {
		result.XmlText = strXml
	} else {
		result.Data = data
	}
	if charSet != nil {
		if len(charSet) > 0 {
			result.CharSet = charSet[0]
		}
	}
	return result
}

/*该函数可以接受两个参数，第一个参数为模板的名称，第二个为主题的名称，两个参数都可以省略，
第一个参数省略时或为""时，Action做为模板名称，第二个参数省略时为默认的主题*/
func (this *Controller) View(args ...string) *ViewResult {
	result := &ViewResult{Response: this.Response, ViewData: this.ViewData, Theme: this.Theme, ViewEngine: this.ViewEngine}
	actionName := fmt.Sprintf("%v", this.RouteData["action"])
	controllerName := fmt.Sprintf("%v", this.RouteData["controller"])
	themeName := this.Theme
	if args != nil {
		if len(args) > 0 {
			actionName = args[0]
		}
		if len(args) > 1 {
			themeName = args[1]
		}
	}
	result.ActionName = actionName
	result.ControllerName = controllerName
	result.Theme = themeName
	areaName, ok := this.RouteData["area"]
	if ok {
		result.Area = fmt.Sprintf("%v", areaName)
	}
	return result
}

func (this *Controller) Redirect(strUrl string) {
	http.Redirect(this.Response, this.Request, strUrl, http.StatusFound)
	this.ResponseEnd()
}
func (this *Controller) BindModel(model interface{}) error {
	return this.DefaultBinder.BindModel(model)
}
