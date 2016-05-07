package Web

import (
	"System/ViewEngine"
	"fmt"
	"net/http"
	"strconv"
	"runtime/debug"
	"html/template"
	"strings"
	"System/Config"
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
	SetReferer(string)
	SetIsAjax(bool)
}
type WebParameter map[string]string //前台提交的Form或Url参数值
//取int值，如果不存在，或不是整型数据返回0
func(this *WebParameter) Int(strParamKey string)int{
	str,ok:=(*this)[strParamKey]
	if !ok{
		return 0//不存在默认返回0
	}
	v,_:=strconv.Atoi(str)
	return v
}
//取Int64值，如果不存在，或不是整型数据返回0
func(this *WebParameter) Int64(strParamKey string)int64{
	str,ok:=(*this)[strParamKey]
	if !ok{
		return 0//不存在默认返回0
	}
	v,_:=strconv.ParseInt(str,10,64)
	return v
}
//取float32值，如果不存在，或不是整型数据返回0
func(this *WebParameter) Float(strParamKey string)float32{
	str,ok:=(*this)[strParamKey]
	if !ok{
		return 0//不存在默认返回0
	}
	v,_:=strconv.ParseFloat(str,32)
	return float32(v)
}
//取float64值，如果不存在，或不是整型数据返回0
func(this *WebParameter) Float64(strParamKey string) float64{
	str,ok:=(*this)[strParamKey]
	if !ok{
		return 0//不存在默认返回0
	}
	v,_:=strconv.ParseFloat(str,64)
	return v
}
func(this *WebParameter) String(strParamKey string)string{
	str,ok:=(*this)[strParamKey]
	if !ok{
		fmt.Println(strParamKey+"不存在")
		return ""//不存在默认返回""
	}
	return str
}
type Controller struct {
	Request       *http.Request
	Response      http.ResponseWriter
	ViewData      map[string]interface{}
	Session       map[string]interface{}
	RouteData     map[string]interface{}
	QueryString   WebParameter
	Form          WebParameter
	ViewEngine    ViewEngine.IViewEngine
	IsEndResponse bool   //是否停止当前请求的执行
	Theme         string //主题名称
	Cookies       map[string]string
	DefaultBinder *Binder
	IsPost        bool //如果是Post提交的数据为true,否则为false
	Referer string //来源网址
	IsAjax bool //如果是ajax则为true,否则为false
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
func (this *Controller) SetReferer(strRefer string){
	this.Referer=strRefer
}
func (this *Controller) SetIsAjax(isAjax bool){
	this.IsAjax=isAjax
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
//页面跳转
func (this *Controller) Redirect(strUrl string) {
	http.Redirect(this.Response, this.Request, strUrl, http.StatusFound)
	this.ResponseEnd()
}
func (this *Controller) BindModel(model interface{}) error {
	return this.DefaultBinder.BindModel(model)
}
//清除Session
func (this *Controller) ClearSession() {
	this.Session=make(map[string]interface{})
}
func(this *Controller) View404()*ViewResult{	
	result := &ViewResult{
		ViewData:       this.ViewData,
		ViewEngine:     this.ViewEngine,
		Response:       this.Response,
		ActionName:     "404",
		ControllerName: "",
		Theme:          App.Configs.Theme}
	return result
}
func(this *Controller) View505(err error)*ViewResult{
	var strErr string
	if Config.AppConfig.ShowErrors{
		if err!=nil{
			strErr=err.Error()
		}
		strErr+="<p>"+string(debug.Stack())+"</p>"
		this.ViewData["ErrMsg"]=template.HTML(strErr)
	}		
	
	result := &ViewResult{
		ViewData:       this.ViewData,
		ViewEngine:     this.ViewEngine,
		Response:       this.Response,
		ActionName:     "505",
		ControllerName: "",
		Theme:          App.Configs.Theme}
	return result
}
//消息提示,title标题,strMsg要提示的消息内容，可以是string或template.HTML类型，strUrl为要跳转的URL地址,waitSecond跳转等待时间 单位为秒

func(this *Controller) Msg(title string,msg interface{},url string,waitSecond int)*ViewResult{
    if	waitSecond<1{
		waitSecond=3
	}
	this.ViewData["title"]=title
	this.ViewData["msg"]=msg
	this.ViewData["url"]=url
	this.ViewData["waitSecond"]=waitSecond
	result := &ViewResult{
		ViewData:       this.ViewData,
		ViewEngine:     this.ViewEngine,
		Response:       this.Response,
		ActionName:     "Msg",
		ControllerName: "",
		Theme:          App.Configs.Theme}
	return result
}
//是否搜索引擎的抓取,true是,false不是
func(this *Controller) IsCrawler() bool{
	agent,ok:=this.Request.Header["User-Agent"]
	if !ok||len(agent)==0{
		return false
	}
	strAgent:=agent[0]
	crawlers:=[]string{"googlebot","baiduspider","sogou","360spider","bingbot","spider","bot","Spider","Bot"}//最后加了spider,bot，所有包含这个都认为是搜索引擎
	for _,crw:=range crawlers{
		if strings.Index(strAgent,crw)!=-1 {
				return true			
		}
	}
	
	return false
}
//获取客户端IP
func(this *Controller) GetClientIp() string{
	strIp:=this.Request.Header.Get("HTTP_CLIENT_IP")
	if strIp!=""{
		return strIp
	}
	strIp=this.Request.Header.Get("HTTP_X_FORWARDED_FOR")
	if strIp!=""{
		return strIp
	}
	strIp=this.Request.Header.Get("REMOTE_ADDR")
	if strIp!=""{
		return strIp
	}
	return strings.Split(this.Request.RemoteAddr,":")[0]
}

//向客户端输出文本内容，并结束请求
func(this *Controller) ResponseText(strText string)*ViewResult{
	this.Response.Write([]byte(strText))
    this.ResponseEnd()
	return this.View()
}
//向客户端输出HTML内容，并结束请求
func(this *Controller) ResponseHTML(strText string)*ViewResult{
	html:=template.HTML(strText)
	this.Response.Write([]byte(html))
    this.ResponseEnd()
	return this.View()
}
//判断是否在微信下，如果是则返回true，如果不是则返回false
func(this *Controller)IsInWeiXin() bool{
	agent:=this.Request.Header.Get("User-Agent")
	if(agent==""||strings.Index(agent,"MicroMessenger")==-1){
		return false
	}else{
		return true
	}	
}