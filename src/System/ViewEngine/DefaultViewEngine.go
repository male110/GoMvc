package ViewEngine

import (
	"System/TemplateFunc"
	"errors"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

var PageNoFind error = errors.New("View不存在")

type IViewEngine interface {
	RenderView(controllerName, actionName, theme string, viewData map[string]interface{}, writer http.ResponseWriter) error
}

/*全局模板*/
type GlobalTemplate struct {
	lastReadTime      time.Time //最后一次更新时间
	gloableTplContent string    //全局模板的内容
}

//默认的模板引擎
type DefaultViewEngine struct {
	SearchLocation []string                   //模板的搜位置
	Extension      string                     //扩展名
	globalTpl      map[string]*GlobalTemplate //按照不同的主题进行存放
	mutex          sync.RWMutex
}

func NewDefualtEngine() *DefaultViewEngine {
	e := &DefaultViewEngine{Extension: ".ghtm"}
	e.SearchLocation =
		[]string{"Views/{theme}/{controller}/{action}",
			"Views/{theme}/_Global/{action}",
			"{area}/Views/{theme}/{controller}/{action}",
			"{area}/Views/{theme}/_Global/{action}"}
	return e
}

//展示
func (this *DefaultViewEngine) RenderView(controllerName, actionName, theme string, viewData map[string]interface{}, writer http.ResponseWriter) error {
	locations := this.getViewLocation(controllerName, actionName, theme)

	//在指定位置搜索模板
	var strTplPath string
	for _, l := range locations {
		_, err := os.Stat(l)
		if err != nil {
			continue
		}
		strTplPath = l
		break
	}
	//模板不存在
	if strTplPath == "" {
		return PageNoFind
	}
	//取全局模板
	glbTpl, err := this.getGlobalTemplate(theme)
	if err != nil {
		return err
	}
	tpl := template.New("view").Funcs(TemplateFunc.TemplatFuncs)

	buf, err := ioutil.ReadFile(strTplPath)
	if err != nil {
		return err
	}
	strTplContent := glbTpl + string(buf)
	tpl, err = tpl.Parse(strTplContent)
	if err != nil {
		return err
	}
	err = tpl.Execute(writer, viewData)
	return err
}
func (this *DefaultViewEngine) getViewLocation(controllerName, actionName, theme string) []string {
	locations := make([]string, len(this.SearchLocation))
	i := 0
	for _, v := range this.SearchLocation {
		str := strings.Replace(v, "{controller}", controllerName, -1)
		str = strings.Replace(str, "{action}", actionName, -1)
		str = strings.Replace(str, "{theme}", theme, -1) + this.Extension
		str = strings.Replace(str, "//", "/", -1)
		locations[i] = str
		i++
	}
	return locations
}

func (this *DefaultViewEngine) getGlobalTemplate(theme string) (string, error) {
	if this.globalTpl == nil {
		this.mutex.Lock()
		this.globalTpl = make(map[string]*GlobalTemplate)
		this.mutex.Unlock()
	}
	this.mutex.RLock()
	globalItem, ok := this.globalTpl[theme]
	this.mutex.RUnlock()
	//为了减小IO操作，每隔一分钟，才对Global进行一次更新
	if !ok || time.Now().Sub(globalItem.lastReadTime).Minutes() > 1 {
		isChange, files, err := this.isGlobalChanged(theme)

		if err != nil {
			return "", err
		}
		if isChange {
			tplContent, err := this.ReadFiles(files...)
			if err != nil {
				return "", err
			}
			if !ok {
				//this.globalTpl中没有，新建一个
				globalItem = &GlobalTemplate{lastReadTime: time.Now(), gloableTplContent: tplContent}
				this.mutex.Lock()
				this.globalTpl[theme] = globalItem
				this.mutex.Unlock()

			} else {
				//this.globalTpl中已经存在，更新模板的内容
				globalItem.gloableTplContent = tplContent
				globalItem.lastReadTime = time.Now()
			}

		}

	}
	return globalItem.gloableTplContent, nil
}

/*判断全局文件是否改变,并返回所有的全局模板路径*/
func (this *DefaultViewEngine) isGlobalChanged(theme string) (bool, []string, error) {
	var files []string
	isChange := false
	strGlobalDir := path.Join("Views", theme, "_Global")

	fs, err := ioutil.ReadDir(strGlobalDir)
	if err != nil {
		return false, nil, nil
	}
	globalItem, ok := this.globalTpl[theme]
	if !ok {
		isChange = true
	}
	tplExt := strings.ToLower(this.Extension)
	for _, v := range fs {
		if v.IsDir() {
			continue
		}
		//对扩展名进行判断,不区分大小写
		strExt := path.Ext(v.Name())
		if strings.ToLower(strExt) != tplExt {

			continue
		}
		if isChange == false {
			if globalItem.lastReadTime.Before(v.ModTime()) {
				isChange = true
			}
		}
		files = append(files, path.Join(strGlobalDir, v.Name()))
	}

	return isChange, files, nil
}

func (this *DefaultViewEngine) ReadFiles(strFileName ...string) (string, error) {
	strContent := ""
	if len(strFileName) == 0 {
		return strContent, nil
	}
	for _, fileName := range strFileName {
		buf, err := ioutil.ReadFile(fileName)
		if err != nil {
			return strContent, err
		}
		strContent += string(buf)
	}
	return strContent, nil
}
