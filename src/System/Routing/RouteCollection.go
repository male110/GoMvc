package Routing

import (
	"errors"
	"strings"
)

type RouteCollection struct {
	/*这里把名称分开存放，主要是因为map是无序的，而在路由的匹配时，我想按添加时的顺序来*/
	routeNames map[string]bool
	routes     []*Route
}
type RouteItem struct {
	Name        string
	Url         string
	Defaults    map[string]interface{}
	Constraints map[string]string
}

/*添加路由,该函数接受三个参数，第一个参数url，为路由的url字符串，如：{controller}/{action}
 *第二个参数为默认值，为map[string]interface{}，
 *第三个参数为约束，为map[string]string，map的值为是一个正则表达式，用来对参数进行验证，如"^(\\d)+$"
 */
func (this *RouteCollection) Add(name, url string, arr ...interface{}) (*Route, error) {
	var defaults map[string]interface{}
	var constraints map[string]string
	if this.routeNames == nil {
		this.routeNames = make(map[string]bool)
	}
	_, ok := this.routeNames[name]
	if ok {
		return nil, errors.New("名为[" + name + "]路由已经存在，重复定义")
	}
	if len(arr) > 0 {
		defaults = arr[0].(map[string]interface{})
	}
	if len(arr) > 1 {
		constraints = arr[1].(map[string]string)
	}

	//去掉开头，结尾的/
	url = strings.Trim(url, "/")
	r, err := NewRoute(url, defaults, constraints)
	if err != nil {
		return nil, err
	}
	this.routeNames[name] = true
	this.routes = append(this.routes, r)
	return r, nil
}
func (this *RouteCollection) AddRote(item *RouteItem) (*Route, error) {
	return this.Add(item.Name, item.Url, item.Defaults, item.Constraints)
}
func (this *RouteCollection) GetRouteData(requestPath string) map[string]interface{} {
	if len(this.routes) > 0 {
		for _, v := range this.routes {
			routeData := v.GetRouteData(requestPath)
			if routeData != nil {
				return routeData
			}
		}
	}
	return nil
}

var RouteTable *RouteCollection = &RouteCollection{routeNames: make(map[string]bool)}
