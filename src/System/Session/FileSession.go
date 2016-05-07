package Session

import (
	. "System/Log"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

/*需要实现,ReadSession,newSession,EndSession,GC*/
type FileSession struct {
	SessionBase
	location string
	timeOut  int //单位分钟,Session超时时间
	gcing    bool
}

/*因为Session保存位置是配置文件中的，可以动态改变，所以这里要判断一下Session位置是否改变，如果改变删除原来的文件*/
func (this *FileSession) checkLocation(path string) {
	/*判断目录是否改了*/
	if this.location != "" && this.location != path {
		os.RemoveAll(this.location)
	}
	this.location = path
}
func (this *FileSession) getSessionFileName(sid string, location string) string {
	return path.Join(location, string(sid[0]), string(sid[1]), sid)
}
func (this *FileSession) readSessionFromFile(sid string, location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	this.checkLocation(location)
	strSessionFileName := this.getSessionFileName(sid, location)
	_, err := os.Stat(strSessionFileName)

	if err == nil {
		//读文件，取Session数据
		m := make(map[string]interface{})
		buf, err := ioutil.ReadFile(strSessionFileName)
		if err != nil {
			AppLog.Add("in FileSession.readSession，读文件时出错：" + err.Error())
			return m, err
		}
		m, _ = GobSerialize.Decode(buf)
		
		return m, nil
	} else {
		if !os.IsNotExist(err) {
			//未知的错误
			AppLog.Add("in FileSession.readSession，读取Session时出错，" + err.Error())
		}
		/*不存在产生新的Session*/
		return this.newSession(location, w, r)
	}
}

func (this *FileSession) CreateNewSession(location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {

	sid := this.newSid()
	this.setSessionName(sid, w, r)
	sessionFilePath := path.Join(location, string(sid[0]), string(sid[1]))
	_, err := os.Stat(sessionFilePath)
	//创建目录
	if err != nil {

		err := os.MkdirAll(sessionFilePath, os.ModePerm)
		if err != nil {
			AppLog.Add("in FileSession.CreateNewSession,创建目录时出错:" + err.Error())
		}
	}
	m := make(map[string]interface{})
	return m, nil
}
func (this *FileSession) EndSession(data map[string]interface{}, location string, r *http.Request) error {
	this.checkLocation(location)
	sid := this.getSessionId(r)
	if sid == "" {
		AppLog.Add("in FileSession.EndSession sid为空。")
		return nil
	}
	sessionFileName := this.getSessionFileName(sid, location)
	buf, _ := GobSerialize.Encode(data)
	err := ioutil.WriteFile(sessionFileName, buf, os.ModePerm)
	if err != nil {
		AppLog.Add("in FileSession.EndSession,写文件时出错：" + err.Error())
		return err
	}
	return nil
}
//删除Session
func (this *FileSession)deleteBySid(sid string,location string) error{
	strFileName:=this.getSessionFileName(sid,this.location)
	return os.Remove(strFileName)
}

func (this *FileSession) GC(timeOut int, location string) {
	this.timeOut = timeOut
	/*如果正在GC中，直接返回，因为遍历可能比较费时*/
	if this.gcing {
		return
	}
	this.gcing = true
	defer func() { this.gcing = false }()
	if this.location != location {
		this.checkLocation(location)
		return
	}
	filepath.Walk(location, this.checkSessionTimeOut)

}
func (this *FileSession) checkSessionTimeOut(path string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return nil
	}
	d := time.Now().Sub(fi.ModTime())
	if d.Minutes() >= float64(this.timeOut) {
		os.Remove(path)
	}
	return nil
}
func NewFileSession() *FileSession {
	fs := &FileSession{gcing: false}
	fs.readSession = fs.readSessionFromFile
	fs.newSession = fs.CreateNewSession
	fs.deleteSession=fs.deleteBySid
	return fs
}
