package ViewEngine

import (
	"System/Config"
	"System/Function"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"
)

var PageNoFind error = errors.New("View不存在")

type IViewEngine interface {
	RenderView(areaName, controllerName, actionName, theme string, viewData map[string]interface{}, writer io.Writer) error
}

/*全局模板*/
type GlobalTemplate struct {
	lastReadTime      time.Time //最后一次更新时间
	gloableTplContent string    //全局模板的内容
}

//默认的模板引擎
type DefaultViewEngine struct {
	ViewLocation     []string                   //模板的搜位置
	AreaViewLocation []string                   //area模板的位置
	Extension        string                     //扩展名
	globalTpl        map[string]*GlobalTemplate //按照不同的主题进行存放
	mutex            sync.RWMutex
}

func NewDefualtEngine() *DefaultViewEngine {
	e := &DefaultViewEngine{Extension: ".ghtm"}
	e.ViewLocation =
		[]string{"Views/{theme}/{controller}/{action}",
			"Views/{theme}/_Global/{action}"}

	e.AreaViewLocation =
		[]string{"Areas/{area}/Views/{theme}/{controller}/{action}",
			"Areas/{area}/Views/{theme}/_Global/{action}"}
	return e
}

//展示
func (this *DefaultViewEngine) RenderView(areaName, controllerName, actionName, theme string, viewData map[string]interface{}, writer io.Writer) error {
	strTplPath := this.getViewPath(areaName, controllerName, actionName, theme)

	//模板不存在
	if strTplPath == "" {
		return PageNoFind
	}
	//取全局模板
	glbTpl, err := this.getGlobalTemplate(areaName, theme)
	if err != nil {
		return err
	}
	tpl := template.New("view").Funcs(TemplatFuncs)

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
func (this *DefaultViewEngine) getViewPath(areaName, controllerName, actionName, theme string) string {
	if areaName == "" {
		//普通的模板
		for _, v := range this.ViewLocation {
			str := strings.Replace(v, "{controller}", controllerName, -1)
			str = strings.Replace(str, "{action}", actionName, -1)
			str = strings.Replace(str, "{theme}", theme, -1) + this.Extension
			str = strings.Replace(str, "//", "/", -1)
			if Function.FileExist(str) {
				return str
			}
		}
	} else {
		//域模板
		for _, v := range this.AreaViewLocation {
			str := strings.Replace(v, "{area}", areaName, -1)
			str = strings.Replace(str, "{controller}", controllerName, -1)
			str = strings.Replace(str, "{action}", actionName, -1)
			str = strings.Replace(str, "{theme}", theme, -1) + this.Extension
			str = strings.Replace(str, "//", "/", -1)
			if Function.FileExist(str) {
				return str
			}
		}
	}
	return ""
}

func (this *DefaultViewEngine) getGlobalTemplate(area, theme string) (string, error) {
	if this.globalTpl == nil {
		this.mutex.Lock()
		this.globalTpl = make(map[string]*GlobalTemplate)
		this.mutex.Unlock()
	}
	strKeyName := area + theme
	this.mutex.RLock()
	globalItem, ok := this.globalTpl[strKeyName]
	this.mutex.RUnlock()
	//为了减小IO操作，每隔一分钟，才对Global进行一次更新
	if Config.AppConfig.IsDebug || !ok || time.Now().Sub(globalItem.lastReadTime).Minutes() > 1 {
		isChange, files, err := this.isGlobalChanged(area, theme)
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
				this.globalTpl[strKeyName] = globalItem
				this.mutex.Unlock()

			} else {
				//this.globalTpl中已经存在，更新模板的内容
				globalItem.gloableTplContent = tplContent
				globalItem.lastReadTime = time.Now()
			}

		}

	}
	if globalItem == nil {
		return "", nil
	}
	return globalItem.gloableTplContent, nil
}

/*判断全局文件是否改变,并返回所有的全局模板路径*/
func (this *DefaultViewEngine) isGlobalChanged(area, theme string) (bool, []string, error) {
	var files []string
	isChange := false
	var strGlobalDir string
	if area == "" {
		strGlobalDir = path.Join("Views", theme, "_Global")
	} else {
		strGlobalDir = path.Join("Areas/", area, "Views", theme, "_Global")
	}

	fs, err := ioutil.ReadDir(strGlobalDir)
	if err != nil {
		return false, nil, nil
	}
	strKeyName := area + theme
	globalItem, ok := this.globalTpl[strKeyName]
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
