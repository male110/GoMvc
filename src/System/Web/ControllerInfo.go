package Web

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ControllerNotExist error = errors.New("Controller不存在")
var ActionNotExist error = errors.New("Action不存在")

type ControllerInfo struct {
	ControllerName string
	ControllerType reflect.Type
	Methods        map[string]string
	AreaName       string
}

type ControllersCollection struct {
	Controllers    map[string]*ControllerInfo
	AreaController map[string]map[string]*ControllerInfo
}

func NewControllerInfo(controllerName string, t reflect.Type) *ControllerInfo {
	c := new(ControllerInfo)
	c.ControllerName = controllerName
	c.ControllerType = t
	c.Methods = make(map[string]string)
	return c
}
func NewControllersCollection() *ControllersCollection {
	c := &ControllersCollection{Controllers: make(map[string]*ControllerInfo)}
	c.AreaController = make(map[string]map[string]*ControllerInfo)
	return c
}

func (this *ControllersCollection) Add(c interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	rt := reflect.TypeOf(c)
	/*取类型名和类型*/
	typeName, t := this.getTypeNameAndType(rt)
	ctlinfo := NewControllerInfo(typeName, t)
	//如果原来是变为非指针,变为指针类型，
	if rt.Kind() != reflect.Ptr {
		rt = reflect.PtrTo(rt)
	}
	/*获取该类型下的函数*/
	this.getTypeMethod(rt, ctlinfo)
	area := this.GetPath(t)

	if area == "" {
		this.Controllers[strings.ToLower(typeName)] = ctlinfo
	} else {
		ctlinfo.AreaName = area
		strLowerArea := strings.ToLower(area)
		_, ok := this.AreaController[strLowerArea]
		if !ok {
			this.AreaController[strLowerArea] = make(map[string]*ControllerInfo)
		}
		this.AreaController[strLowerArea][strings.ToLower(typeName)] = ctlinfo
	}
	return nil
}
func (this *ControllersCollection) GetPath(rt reflect.Type) string {
	strPath := rt.PkgPath()
	if strPath == "" || len(strPath) < 7 {
		return ""
	}
	strErrMsg := strPath + "/" + rt.Name() + "目录结构错误，Arear的目录结构为Area\\域名\\Controller"
	arrPath := []rune(strPath)
	strPrifix := strings.ToLower(string(arrPath[0:6]))
	if strPrifix != "areas/" {
		return ""
	}
	strPath = string(arrPath[6:])

	i := strings.IndexAny(strPath, "/")
	if i == -1 {
		panic(strErrMsg)
		return ""
	}
	arrPath = []rune(strPath)
	area := string(arrPath[0:i])
	return area
}
func (this *ControllersCollection) getTypeMethod(t reflect.Type, ctlinfo *ControllerInfo) {
	for i, j := 0, t.NumMethod(); i < j; i++ {
		m := t.Method(i)
		strMethodName := strings.ToLower(m.Name)
		ctlinfo.Methods[strMethodName] = m.Name
	}
}
func (this *ControllersCollection) getTypeNameAndType(t reflect.Type) (string, reflect.Type) {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	strName := t.Name()
	if strings.HasSuffix(strName, "Controller") && len(strName) > len("Controller") {
		strName = strings.TrimRight(strName, "Controller")
	}
	return strName, t
}

func (this *ControllersCollection) GetController(routeData map[string]interface{}) (reflect.Value, error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err := e.(error)
			App.Log.Add("In ControllerInfo.GetController:\t" + fmt.Sprintf("%v", err.Error()))
		}
	}()
	var result reflect.Value
	controllerName := strings.ToLower(fmt.Sprintf("%v", routeData["controller"]))
	actionName := strings.ToLower(fmt.Sprintf("%v", routeData["action"]))
	var controllers map[string]*ControllerInfo
	area, ok := routeData["area"]
	//如果是area,则取area下的controller,否则取默认的
	if ok {
		strArea := strings.ToLower(fmt.Sprintf("%v", area))
		controllers, ok = this.AreaController[strArea]
		if !ok {
			return result, ControllerNotExist
		}
	} else {
		controllers = this.Controllers
	}
	ctlinfo, ok := controllers[controllerName]
	if !ok {
		return result, ControllerNotExist
	}
	methodName, ok := ctlinfo.Methods[actionName]
	if !ok {
		return result, ActionNotExist
	}
	routeData["controller"] = ctlinfo.ControllerName
	routeData["action"] = methodName
	routeData["area"] = ctlinfo.AreaName
	result = reflect.New(ctlinfo.ControllerType)
	return result, nil
}
