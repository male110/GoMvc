package Routing

import (
	"regexp"
)

type Route struct {
	url         string
	urlParser   *PathParser
	Defaults    map[string]interface{}
	constraints map[string]*regexp.Regexp
}

func (this *Route) Parse(url string) error {
	routParse := new(RouteParser)
	var err error
	this.urlParser, err = routParse.ParseUrl(url)
	return err
}
func (this *Route) AddDefault(paramName string, defaultValue interface{}) {
	if this.Defaults == nil {
		this.Defaults = make(map[string]interface{})
	}
	this.Defaults[paramName] = defaultValue
}
func (this *Route) AddConstraint(paramName string, regex string) error {
	reg, err := regexp.Compile(regex)
	if err != nil {
		return err
	}
	if this.constraints == nil {
		this.constraints = make(map[string]*regexp.Regexp)
	}
	this.constraints[paramName] = reg
	return nil
}

/*func (this *Route) GetRouteData(request *http.Request) map[string]interface{} {
requestPath := request.URL.Path*/
func (this *Route) GetRouteData(requestPath string) map[string]interface{} {
	routData := this.urlParser.Match(requestPath, this.Defaults)
	if routData == nil {
		return nil
	}
	//进行约束检查
	if !this.ProcessConstraints(routData) {
		return nil
	}
	return routData
}

func (this *Route) ProcessConstraints(routData map[string]interface{}) bool {
	if this.constraints == nil || len(this.constraints) == 0 {
		return true
	}
	for k, v := range routData {
		reg, ok := this.constraints[k]
		if !ok {
			continue
		}
		str, ok := v.(string)
		//如果不是string类型，就是取的默认值,默认值不需做验证
		if !ok {
			continue
		}
		if !reg.MatchString(str) {
			return false
		}
	}
	return true
}
func NewRoute(url string, defaults map[string]interface{}, constraints map[string]string) (*Route, error) {
	r := &Route{Defaults: defaults}
	err := r.Parse(url)
	if err != nil {
		return nil, err
	}
	if constraints != nil {
		for k, v := range constraints {
			err = r.AddConstraint(k, v)
			if err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}
