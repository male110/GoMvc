package Log

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

type Logger struct {
	nowDate         string //当前时间
	fileIndex       int
	maxFileSize     int64 //单个文件的大小
	fileSize        int64 //文件的当前大小
	lock            sync.Mutex
	logDir          string //日志存放的目录
	currentFileName string //当前日志文件的名称
	currentPath     string //当前路径
}

/*
strLocation:日志保存的位置
maxFileSize:文件大小限制,单位MB
*/
func New(strLocation string, maxLogSize float64) *Logger {

	log := new(Logger)
	log.SetLocation(strLocation)
	log.SetMaxSize(maxLogSize)
	return log
}

//生成日志的文件路径
func (this *Logger) generateFileName() error {
	if this.logDir == "" {
		this.logDir = "Log/"
	}
	date := time.Now()
	//取当前时间，年月日形式
	strDate := date.Format("2006-01-02")
	//判断日志保存路径是否改变
	if this.currentFileName != "" && !strings.HasPrefix(this.currentFileName, this.logDir) {
		this.currentFileName = ""
		this.fileIndex = 0
		this.fileSize = 0
		this.currentPath = ""
	}
	//nowDate==当前日期
	if strDate == this.nowDate && this.currentFileName != "" {
		//判断文件大小是否大于最大限制，如果大于等于最大限制，产生一个新的日志
		if this.fileSize >= this.maxFileSize {
			this.fileIndex += 1
			strFileName := this.getNewFileName()
			this.currentFileName = path.Join(this.currentPath, strFileName)
		}
		return nil
	}
	////nowDate!=当前日期,产生一个新的日志文件
	this.nowDate = strDate
	strYear := fmt.Sprintf("%v", date.Year())
	strMonth := fmt.Sprintf("%v", int(date.Month()))
	//以年/月的目录结构来存放日志
	strPath := path.Join(this.logDir, strYear, strMonth)
	if this.currentPath != strPath {
		err := os.MkdirAll(strPath, os.ModePerm)
		if err != nil {
			return err
		}
		this.currentPath = strPath
	}
	this.currentFileName = this.getNewFileName()
	for {
		fileinfo, err := os.Stat(this.currentFileName)
		if err != nil {
			break
		}
		if fileinfo.Size() >= this.maxFileSize {
			this.fileIndex++
			this.currentFileName = this.getNewFileName()
		} else {
			break
		}
	}

	this.fileIndex = 0
	this.fileSize = 0
	return nil
}

func (this *Logger) getNewFileName() string {
	if this.fileIndex > 0 {
		return path.Join(this.currentPath, this.nowDate+"_"+fmt.Sprintf("%v", this.fileIndex)+".log")
	} else {
		return path.Join(this.currentPath, this.nowDate+".log")
	}

}

/*写文件，并更新文件大小*/
func (this *Logger) writeFile(content string) error {
	//打开文件
	file, err := os.OpenFile(this.currentFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	//写入日志
	_, err = file.Write([]byte(content))
	if err != nil {
		return err
	}
	//取文件的大小
	var fileInfo os.FileInfo
	fileInfo, err = file.Stat()
	if err != nil {
		return err
	}
	this.fileSize = fileInfo.Size()
	return nil
}

/*添加一条日志*/
func (this *Logger) Add(content string) error {
	content = time.Now().Format("2006-01-02 15:04:05") + "\t" + content + "\r\n"
	this.lock.Lock()
	defer this.lock.Unlock()
	err := this.generateFileName()
	if err != nil {
		return err
	}
	err = this.writeFile(content)
	return err
}

/*添加一条错误信息*/
func (this *Logger) AddError(err error) error {
	content := time.Now().Format("2006-01-02 15:04:05") + "\t" + err.Error() + "\r\n"
	this.lock.Lock()
	defer this.lock.Unlock()
	err1 := this.generateFileName()
	if err1 != nil {
		return err1
	}
	err1 = this.writeFile(content)
	return err1
}

/*设置日志的存放位置*/
func (this *Logger) SetLocation(location string) {
	if location == "" {
		location = "Log" //默认位置
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.logDir = location
}

/*设置日志文件的大小设置,单位MB*/
func (this *Logger) SetMaxSize(maxSize float64) {
	size := int64(maxSize * 1024 * 1024)
	if size <= 0 {
		size = 5 * 1024 * 1024 //默认大小5M
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	this.maxFileSize = size

}

/*创建一个日志对像，记录在默认位置，默认大小,在配置文件加载后，重新设置保存位置和大小*/
var AppLog = New("", 0)
