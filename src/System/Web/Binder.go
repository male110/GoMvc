package Web

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Binder struct {
	Request   map[string]string
	Post      map[string]string
	RouteData map[string]interface{}
}

func NewBinder(request *http.Request, routeData map[string]interface{}) *Binder {
	binder := new(Binder)
	//初始化变量
	binder.Request = make(map[string]string)
	binder.Post = make(map[string]string)
	binder.RouteData = make(map[string]interface{})
	querys := request.URL.Query()
	forms := request.PostForm
	//取URL地址栏的参数，对于同名的，以最后一次的为准,key转换为小写，取值时不区分大小写
	if querys != nil {
		for k, v := range querys {
			k = strings.ToLower(k)
			binder.Request[k] = v[len(v)-1]
		}
	}
	//取form的值，对于同名的，以最后一次的为准，key转换为小写，取值时不区分大小写
	if forms != nil {
		for k, v := range forms {
			k = strings.ToLower(k)
			binder.Post[k] = v[len(v)-1]
		}
	}
	//key转换为小写，取值时不区分大小写
	if routeData != nil {
		for k, v := range routeData {
			k = strings.ToLower(k)
			binder.RouteData[k] = v
		}
	}
	return binder
}
func (this *Binder) BindModel(data interface{}) error {
	var rv reflect.Value
	var rt reflect.Type
	//传进来的值，有可能是用reflect.New创建的
	switch dataType := data.(type) {
	case reflect.Value:
		rv = dataType
		rt = dataType.Type()
	default:
		rv = reflect.ValueOf(data)
		rt = reflect.TypeOf(data)
	}

	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	//只能对Struct结构体进行绑定
	if rt.Kind() != reflect.Struct {
		return errors.New("data参数必须是结构体")
	}
	numField := rt.NumField()
	for i := 0; i < numField; i++ {
		ft := rt.Field(i)
		fv := rv.Field(i)
		value := fv.Interface()
		if fv.Kind() == reflect.Struct {
			//如果字段是一个结构体，不做处理
			continue
		} else {
			rfvaleType := reflect.TypeOf(value)
			//取字段名对应的值，顺序是Post,query,route
			paramValue := this.getValue(ft.Name)
			//针对不同的类型进行类型转换
			switch value.(type) {
			case string:
				strValue := this.asString(paramValue)
				fv.SetString(strValue)
			case int, int8, int16, int32, int64:
				strValue := this.asString(paramValue)
				intValue, err := strconv.ParseInt(strValue, 10, rfvaleType.Bits())
				if err != nil {
					return nil
				}
				fv.SetInt(intValue)
			case uint, uint8, uint16, uint32, uint64:
				strValue := this.asString(paramValue)
				uintValue, err := strconv.ParseUint(strValue, 10, rfvaleType.Bits())
				if err != nil {
					return nil
				}
				fv.SetUint(uintValue)
			case float32, float64:
				strValue := this.asString(paramValue)
				floatValue, err := strconv.ParseFloat(strValue, rfvaleType.Bits())
				if err != nil {
					return nil
				}
				fv.SetFloat(floatValue)
			case bool:
				strValue := this.asString(paramValue)
				boolValue, err := strconv.ParseBool(strValue)
				if err != nil {
					return nil
				}
				fv.SetBool(boolValue)
			case interface{}:
				fv.Set(reflect.ValueOf(paramValue))
			}
		}
	}
	return nil
}

func (this *Binder) getValue(paramName string) interface{} {
	paramName = strings.ToLower(paramName)
	//先从post中取
	str, ok := this.Post[paramName]
	if ok {
		return str
	}
	//再从url地址栏的参数中取
	str, ok = this.Request[paramName]
	if ok {
		return str
	}
	//最后从RouteData中取
	data, ok := this.RouteData[paramName]
	if ok {
		return data
	}
	return nil
}

func (this *Binder) asString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	return fmt.Sprintf("%v", src)
}
