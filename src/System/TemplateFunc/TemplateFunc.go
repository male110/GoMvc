package TemplateFunc

import (
	"fmt"
	. "html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
	return Get(strUrl, r.Host, strCookies)
}
func Get(url, host, cookies string) HTML {
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
	if response.StatusCode != 200 && err != nil {
		return HTML(err.Error())
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
	for i, j := 0, len(arrCookies); i < j; i++ {
		c := arrCookies[i]
		strCookie = c.Name + "=" + c.Value + ";"
	}
	strCookie = strings.Trim(strCookie, ";")
	return strCookie
}
func GetUrl(r *http.Request) string {
	strRefer := strings.ToLower(r.Referer())
	strUrl := ""
	if strings.Index(strRefer, "https://") == -1 {
		strUrl = "http://"
	} else {
		strUrl = "https://"
	}
	strUrl = strUrl + r.Host
	strUrl = strings.Trim(strUrl, "/")
	return strUrl
}
