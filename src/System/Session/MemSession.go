package Session

import (
	"net/http"
	"sync"
	"time"
)

/*需要实现,ReadSession,newSession,EndSession,GC*/
type MemSession struct {
	SessionBase
	sessions map[string]MemSessionItem
	mutex    sync.RWMutex
	gcing    bool
}
type MemSessionItem struct {
	lastAccessTime time.Time
	sessionData    map[string]interface{}
}

/*该函数被StartSession调用，StartSession在基类SessionBase中实现，以减少子类的代码量*/
func (this *MemSession) GetSession(sid string, location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	sessItem, ok := this.sessions[sid]
	if ok {
		this.mutex.RLock()
		defer this.mutex.RUnlock()
		sessItem.lastAccessTime = time.Now()
		this.sessions[sid] = sessItem
		return sessItem.sessionData, nil
	} else {
		//直接返
		return this.newSession("", w, r)
	}
}

/*产生一个新的session*/
func (this *MemSession) CreateNewSession(location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	sid := this.newSid()
	this.setSessionName(sid, w, r)
	sessItem := MemSessionItem{lastAccessTime: time.Now(), sessionData: make(map[string]interface{})}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.sessions[sid] = sessItem
	return sessItem.sessionData, nil
}

/*在请求处理结束时调用，与StartSession相对应,用来把Session数据存回存储介质中*/
func (this *MemSession) EndSession(data map[string]interface{}, location string, r *http.Request) error {
	/*在ReadSession中，返回的是map，map是引用类型，对他的修改会直接反应SessionItem中，所以直接返回*/
	return nil
}

/*定时对Session进行清理，timeOut是Session超期时间，单位分钟*/
func (this *MemSession) GC(timeOut int, location string) {
	if this.gcing {
		return
	}
	this.gcing = true
	defer func() { this.gcing = false }()
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	for key, item := range this.sessions {
		d := time.Now().Sub(item.lastAccessTime)
		if d.Minutes() >= float64(timeOut) {
			this.mutex.RUnlock()
			this.mutex.Lock()
			delete(this.sessions, key)
			this.mutex.Unlock()
			this.mutex.RLock()
		}
	}
}
//删除session
func (this *MemSession)deleteBySid(sid string,location string) error{
	delete(this.sessions,sid)
	return nil
}
func NewMemSession() *MemSession {
	ms := &MemSession{sessions: make(map[string]MemSessionItem), gcing: false}
	ms.readSession = ms.GetSession
	ms.newSession = ms.CreateNewSession
	ms.deleteSession=ms.deleteBySid
	return ms
}
