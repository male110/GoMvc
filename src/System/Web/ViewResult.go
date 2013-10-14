package Web

import (
	"System/ViewEngine"
	"net/http"
)

type IActionResult interface {
	ExecuteResult() error
}
type ViewResult struct {
	ViewData       map[string]interface{}
	ViewEngine     ViewEngine.IViewEngine
	Response       http.ResponseWriter
	ActionName     string
	ControllerName string
	Theme          string //主题
	Area           string //区
}

func (this *ViewResult) ExecuteResult() error {
	this.Response.Header().Add("Content-Type", "text/html;charset=utf-8")
	return this.ViewEngine.RenderView(this.Area, this.ControllerName, this.ActionName, this.Theme, this.ViewData, this.Response)
}
