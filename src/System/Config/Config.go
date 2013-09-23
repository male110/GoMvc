package Config

import (
	. "System/Log"
	"encoding/xml"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type StaticFile struct {
	Url      string
	FilePath string
}
type Config struct {
	ShowErrors       bool
	CookieDomain     string
	LogPath          string
	LogFileMaxSize   float64 //单位MB，日志文件的大小限制
	DriverName       string
	DriverSourceName string
	StaticDir        []string
	StaticFiles      []StaticFile
	SessionType      int //1,文件,2内存,3数据库
	SessionLocation  string
	SessionTimeOut   int
	MemFreeInterval  int    //单位秒
	Theme            string //当前使用的主题
	ListenPort       int    //Http监听的端口号，该配置改后必须重启程序才能生效
	lastModifyTime   time.Time
	loadTime         time.Time //xml加载时间
}

func NewConfig() *Config {
	c := NewDefault()
	LoadConfig(c)
	return c
}

/*从配置文件加载配置信息*/
func LoadConfig(c *Config) {
	tem := time.Now().Sub(c.loadTime)
	if tem.Seconds() < 1 {
		//离上次加载时间不到一秒，返回
		return
	}
	c.loadTime = time.Now()
	//取文件的最后修改时间，看是否是最新的，不是最新才更新
	fileInfo, ferr := os.Stat("web.config")

	if ferr != nil {
		AppLog.AddError(ferr)
		return
	}
	lastModify := fileInfo.ModTime()
	if !lastModify.After(c.lastModifyTime) {
		//无需更新
		return
	}
	c.lastModifyTime = lastModify
	file, err := os.OpenFile("web.config", os.O_RDONLY, os.ModePerm)
	if err != nil {
		AppLog.AddError(err)

		return
	}
	defer file.Close()
	decoder := xml.NewDecoder(file)
	var token xml.Token
	var tokenName string
	for {
		token, err = decoder.Token()
		if err != nil {
			if err != io.EOF {
				AppLog.AddError(err)
			}
			break
		}
		switch t := token.(type) {
		case xml.StartElement:
			tokenName = t.Name.Local
			if tokenName == "File" {
				var url, strFile string
				//取属性信息
				for _, attr := range t.Attr {
					if attr.Name.Local == "url" {
						url = attr.Value
						if strings.TrimSpace(url) == "" {
							AppLog.Add("静态文件配置错误,url不能为空")
							continue
						}
					}
					if attr.Name.Local == "filePath" {
						strFile = changePathSeparator(attr.Value)
						if strings.TrimSpace(strFile) == "" {
							AppLog.Add("静态文件配置错误，文件名不能为空")
							continue
						}
					}
				}
				//判断文件是否存在
				if !isExist(strFile) {
					AppLog.Add("静态文件配置错误,文件" + strFile + "不存在")
					continue
				}
				//转换为小写
				url = strings.ToLower(url)
				//去掉开头，结尾的"/"
				url = strings.Trim(url, "/")
				f := StaticFile{Url: url, FilePath: strFile}
				c.StaticFiles = append(c.StaticFiles, f)
			}
		case xml.EndElement:
			tokenName = ""
		case xml.CharData:
			if tokenName != "" {
				processXmlTocken(c, tokenName, string([]byte(t)))
			}
		}
	}
	//重新设置日志保存位置
	AppLog.SetLocation(c.LogPath)
	AppLog.SetMaxSize(c.LogFileMaxSize)
}
func processXmlTocken(c *Config, xmlName, data string) {
	data = strings.Trim(data, "\n")
	switch xmlName {
	case "ShowErrors":
		var err error
		c.ShowErrors, err = strconv.ParseBool(data)
		if err != nil {
			AppLog.Add("解析配置文件ShowErrors时出错:" + err.Error() + ",配置错误，只能是true或false")
		}
	case "CookieDomain":
		c.CookieDomain = data
	case "Theme":
		c.Theme = data
	case "LogPath":
		c.LogPath = changePathSeparator(data)
	case "LogFileMaxSize":
		size, err := strconv.ParseFloat(data, 64)
		if err != nil {
			AppLog.Add("解析LogFileMaxSize时出错：" + err.Error() + "，xmlName:" + xmlName + "，data:" + data)
		} else {
			c.LogFileMaxSize = size
		}
	case "DriverName":
		c.DriverName = data
	case "DataSourceName":
		c.DriverSourceName = data
	case "Dir":
		c.StaticDir = append(c.StaticDir, changePathSeparator(data))
	case "SessionType":
		stype, err := strconv.Atoi(data)
		if err != nil {
			AppLog.Add("解析SessionType时出错：" + err.Error() + "，xmlName:" + xmlName + "，data:" + data)
		} else {
			c.SessionType = stype
		}
	case "SessionLocation":
		c.SessionLocation = changePathSeparator(data)
	case "SessionTimeOut":
		timeout, err := strconv.Atoi(data)
		if err != nil {
			AppLog.Add("解析SessionTimeOut时出错：" + err.Error() + "，xmlName:" + xmlName + "，data:" + data)
		} else {
			c.SessionTimeOut = timeout
		}
	case "MemFreeInterval":
		interval, err := strconv.Atoi(strings.TrimSpace(data))
		if err != nil {
			AppLog.Add("解析MemFreeInterval时出错：" + err.Error() + "，xmlName:" + xmlName + "，data:" + data)
		} else {
			c.MemFreeInterval = interval
		}
	case "ListenPort":
		port, err := strconv.Atoi(strings.TrimSpace(data))
		if err != nil {
			AppLog.Add("解析ListenPort时出错：" + err.Error() + "，xmlName:" + xmlName + "，data:" + data)
		} else {
			c.ListenPort = port
		}
	}
}

/*用来把Windows下的路径分隔符\改/,因为Go里定义的分隔附是/*/
func changePathSeparator(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

/*判断文件或目录是否存在*/
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
func NewDefault() *Config {
	c := &Config{Theme: "default", LogPath: "Log", LogFileMaxSize: 5, DriverName: "mysql", SessionType: 1, SessionLocation: "sessions", SessionTimeOut: 30, MemFreeInterval: 60}
	return c
}

var AppConfig *Config = NewConfig()
