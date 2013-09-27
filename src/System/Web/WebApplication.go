package Web

import (
	"System/Config"
	"System/Log"
	"System/Session"
	"System/ViewEngine"
	"System/fsnotify"
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"
)

type WebApplication struct {
	Log             *Log.Logger
	watcher         *fsnotify.Watcher //对文件修改进行监控，如配置文件
	Configs         *Config.Config
	controllers     *ControllersCollection
	SessionProvider Session.ISession
	ViewEngine      ViewEngine.IViewEngine
	IsInit          bool
}

func (this *WebApplication) Init() {
	this.Log = Log.AppLog
	this.Configs = Config.AppConfig
	this.SessionProvider = Session.NewSession(this.Configs.SessionType)
	this.controllers = NewControllersCollection()
	this.ViewEngine = ViewEngine.NewDefualtEngine()
	this.initWatcher()
	this.IsInit = true
}

//初始化文件监控对像
func (this *WebApplication) initWatcher() {
	var err error
	this.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		this.Log.AddError(err)
		return
	}
	err = this.watcher.Watch("web.config")
	if err != nil {
		this.Log.AddError(err)
		return
	}
	go this.watchModify()
}
func (this *WebApplication) RegisterController(ctl interface{}) {
	/*if !App.IsInit {
		App.Init()
	}*/
	err := this.controllers.Add(ctl)
	if err != nil {
		this.Log.Add("注册Controllor时出错：" + err.Error())
	}
}
func (this *WebApplication) GetController(routeData map[string]interface{}) (reflect.Value, error) {
	return this.controllers.GetController(routeData)
}
func (this *WebApplication) Run() error {
	go this.SessionGC()
	handler := new(HttpHandler)
	strPort := fmt.Sprintf(":%v", App.Configs.ListenPort)
	err := http.ListenAndServe(strPort, handler)
	return err
}

//监控文件修改的消息
func (this *WebApplication) watchModify() {
	for {
		select {
		case ev := <-this.watcher.Event:
			if ev.IsModify() {
				strFilePath := path.Clean(strings.Replace(ev.Name, "\\", "/", -1))
				switch strFilePath {
				case "web.config":
					Config.LoadConfig(this.Configs)
				}
			}
		}
	}
}
func (this *WebApplication) SessionGC() {
	for {
		n := time.Duration(this.Configs.MemFreeInterval) * time.Second
		timer := time.After(n)
		<-timer
		if this.SessionProvider != nil {
			this.SessionProvider.GC(this.Configs.SessionTimeOut, this.Configs.SessionLocation)
		}
	}
}

var App = new(WebApplication)

func init() {
	//if !App.IsInit {
	App.Init()
	//}
}
