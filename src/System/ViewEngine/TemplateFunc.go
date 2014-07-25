package ViewEngine

import (
	"System/Config"
	"System/Log"
	"bytes"
	"fmt"
	. "html/template"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var TemplatFuncs FuncMap = make(FuncMap)

func init() {
	TemplatFuncs["RanderAction"] = RanderAction
	TemplatFuncs["Equal"] = Equal
	TemplatFuncs["Greater"] = Greater
	TemplatFuncs["GreaterOrEqual"] = GreaterOrEqual
	TemplatFuncs["Less"] = Less
	TemplatFuncs["LessOrEqual"] = LessOrEqual
	TemplatFuncs["SubString"] = SubString
	TemplatFuncs["Trim"] = Trim
	TemplatFuncs["TrimSpace"] = TrimSpace
	TemplatFuncs["IsOddNumber"] = IsOddNumber
	TemplatFuncs["Mod"] = Mod
	TemplatFuncs["RenderView"] = RenderView
	TemplatFuncs["FormatTime"] = FormatTime
	TemplatFuncs["AddValue"] = AddValue
	TemplatFuncs["RandomMetroCSS"] = RandomMetroCSS
}

//等于
func Equal(a, b interface{}) bool {
	result := a == b
	//fmt.Println(a, "\t", b, "\t", result)
	return result
}

//大于
func Greater(a, b interface{}) bool {
	strB := fmt.Sprintf("%v", b)
	strA := fmt.Sprintf("%v", a)
	switch at := a.(type) {
	case int, int8, int16, int32, int64:
		intB, _ := strconv.ParseInt(strB, 10, 64)
		intA, _ := strconv.ParseInt(strA, 10, 64)
		return intA > intB
	case uint, uint8, uint16, uint32, uint64:
		uintB, _ := strconv.ParseUint(strB, 10, 64)
		uintA, _ := strconv.ParseUint(strA, 10, 64)
		return uintA > uintB
	case float32, float64:
		floatB, _ := strconv.ParseFloat(strB, 64)
		floatA, _ := strconv.ParseFloat(strA, 64)
		return floatA > floatB
	case string:
		return at > strB
	}
	return false
}

//大于等于
func GreaterOrEqual(a, b interface{}) bool {
	strB := fmt.Sprintf("%v", b)
	strA := fmt.Sprintf("%v", a)
	switch at := a.(type) {
	case int, int8, int16, int32, int64:
		intB, _ := strconv.ParseInt(strB, 10, 64)
		intA, _ := strconv.ParseInt(strA, 10, 64)
		return intA >= intB
	case uint, uint8, uint16, uint32, uint64:
		uintB, _ := strconv.ParseUint(strB, 10, 64)
		uintA, _ := strconv.ParseUint(strA, 10, 64)
		return uintA >= uintB
	case float32, float64:
		floatB, _ := strconv.ParseFloat(strB, 64)
		floatA, _ := strconv.ParseFloat(strA, 64)
		return floatA >= floatB
	case string:
		return at >= strB
	}
	return false
}

//小于
func Less(a, b interface{}) bool {
	strB := fmt.Sprintf("%v", b)
	strA := fmt.Sprintf("%v", a)
	switch at := a.(type) {
	case int, int8, int16, int32, int64:
		intB, _ := strconv.ParseInt(strB, 10, 64)
		intA, _ := strconv.ParseInt(strA, 10, 64)
		return intA < intB
	case uint, uint8, uint16, uint32, uint64:
		uintB, _ := strconv.ParseUint(strB, 10, 64)
		uintA, _ := strconv.ParseUint(strA, 10, 64)
		return uintA < uintB
	case float32, float64:
		floatB, _ := strconv.ParseFloat(strB, 64)
		floatA, _ := strconv.ParseFloat(strA, 64)
		return floatA < floatB
	case string:
		return at < strB
	}
	return false
}

//小于等于
func LessOrEqual(a, b interface{}) bool {
	strB := fmt.Sprintf("%v", b)
	strA := fmt.Sprintf("%v", a)
	switch at := a.(type) {
	case int, int8, int16, int32, int64:
		intB, _ := strconv.ParseInt(strB, 10, 64)
		intA, _ := strconv.ParseInt(strA, 10, 64)
		return intA <= intB
	case uint, uint8, uint16, uint32, uint64:
		uintB, _ := strconv.ParseUint(strB, 10, 64)
		uintA, _ := strconv.ParseUint(strA, 10, 64)
		return uintA <= uintB
	case float32, float64:
		floatB, _ := strconv.ParseFloat(strB, 64)
		floatA, _ := strconv.ParseFloat(strA, 64)
		return floatA <= floatB
	case string:
		return at <= strB
	}
	return false
}

//字符串截取
func SubString(str string, start, length int) string {
	end := start + length
	arr := []rune(str)
	if end > len(arr) {
		end = len(arr)
	}
	if end < start {
		t := end
		end = start
		start = t
	}
	str = string(arr[start:end])
	return str
}
func Trim(str, cutset string) string {
	return strings.Trim(str, cutset)
}
func TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

func RanderAction(controller, action, param string, r *http.Request) HTML {
	strUrl := "/" + controller + "/" + action
	param = strings.TrimSpace(param)
	if param != "" {
		if strings.HasPrefix(param, "?") {
			strUrl = strUrl + param
		} else {
			strUrl = strUrl + "?" + param
		}
	}
	strCookies := GetCookies(r)
	strUrl = GetUrl(r) + strUrl
	return Get(strUrl, "localhost", strCookies, r)
}
func Get(url, host, cookies string, r *http.Request) HTML {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return HTML(err.Error())
	}
	req.Header.Add("Cookie", cookies)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("User-Agent", "	Mozilla/5.0 (Windows NT 5.1; rv:28.0) Gecko/20100101 Firefox/28.0")
	req.Header.Add("Host", host)
	response, err := client.Do(req)
	if err != nil {
		return HTML(err.Error())
	}
	if response.StatusCode != 200 {
		if err != nil {
			return HTML(err.Error())
		} else {
			Log.AppLog.Add("RenderAction URL:" + url + "\r\nHost:" + host + "\r\nStatusCode:" + strconv.Itoa(response.StatusCode))
			Log.AppLog.Add("In RenderAction request.URL:" + r.URL.String() + "\r\nReferer:" + r.Referer())
		}

	}
	buf, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		return HTML(err.Error())
	}
	return HTML(string(buf))
}
func GetCookies(r *http.Request) string {
	strCookie := ""
	arrCookies := r.Cookies()
	for _, c := range arrCookies {
		strCookie += c.Name + "=" + c.Value + ";"
	}
	strCookie = strings.Trim(strCookie, ";")
	return strCookie
}
func GetUrl(r *http.Request) string {
	strRefer := strings.ToLower(r.URL.String())

	strUrl := ""
	if strings.Index(strRefer, "https://") == -1 {
		strUrl = "http://"
	} else {
		strUrl = "https://"
	}
	strUrl = strUrl + "localhost"
	if Config.AppConfig.ListenPort != 80 {
		strUrl = strUrl + ":" + strconv.Itoa(Config.AppConfig.ListenPort)
	}

	strUrl = strings.Trim(strUrl, "/")
	return strUrl
}

//在模板中嵌入另一个模板文件
func RenderView(strViewName string, viewData map[string]interface{}) HTML {
	var strController, strTheme, strArea string
	if strViewName == "" {
		return HTML("RederView viewName can't empty")
	}
	temp, ok := viewData["Controller"]
	if ok {
		strController = temp.(string)
	}
	temp, ok = viewData["Theme"]
	if ok {
		strTheme = temp.(string)
	}
	temp, ok = viewData["Area"]
	if ok {
		strArea = temp.(string)
	}

	if strController == "" {
		strController = "_Global"
	}
	if strTheme == "" {
		strTheme = "default"
	}

	viewEngine := NewDefualtEngine()
	writer := new(bytes.Buffer)
	err := viewEngine.RenderView(strArea, strController, strViewName, strTheme, viewData, writer)
	if err != nil {
		Log.AppLog.AddErrMsg("RenderView出错 area:" + strArea + ",controller:" + strController + ",viewName:" + strViewName + "\r\n" + err.Error())
	}
	return HTML(writer.String())
}

//取余
func Mod(x, y float64) float64 {
	return math.Mod(x, y)
}

//判断是否是奇数
func IsOddNumber(x int) bool {
	i := Mod(float64(x), 2)
	return i != 0
}
func FormatTime(t time.Time, strFormat string) string {
	return t.Format(strFormat)
}

//函数必须有一个返回值
func AddValue(m map[string]interface{}, key string, value interface{}) string {
	m[key] = value
	return ""
}

//随机返回一个Css样式
func RandomMetroCSS() string {
	arrCss := []string{"amber", "blue", "brown", "cobalt", "crimson", "cyan", "magenta", "lime", "indigo", "green", "emerald", "mango", "mauve", "olive", "orange", "pink", "red", "sienna", "steel", "teal", "violet", "yellow"}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	l := len(arrCss)
	index := r.Int31n(int32(l))
	return arrCss[index]
}
