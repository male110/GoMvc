package Session

import (
	"System/Config"
	. "System/Log"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"
	"strings"
)

var SessionKeyName = "GoMvcSession"

type ISession interface {
	/*应该在SessionStart里修改Session最后一次的访问时间，并返回Session数据，map[string]interface{}*/
	StartSession(w http.ResponseWriter, r *http.Request, location string) (map[string]interface{}, error)
	/*在请求处理结束时调用，与StartSession相对应,用来把Session数据存回存储介质中*/
	EndSession(data map[string]interface{}, location string, r *http.Request) error
	/*定时对Session进行清理，timeOut是Session超期时间，单位分钟*/
	GC(timeOut int, location string)
}

type SessionBase struct {
	/*用来实际获取Session*/
	readSession func(sid string, location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error)
	/*产生一个新的Session*/
	newSession func(location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error)
}

func (this *SessionBase) StartSession(w http.ResponseWriter, r *http.Request, location string) (map[string]interface{}, error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err := e.(error)
			AppLog.Add("In Session.StartSession:\t" + err.Error())
		}
	}()
	sid := this.getSessionId(r)
	if sid == "" {
		//cookie不存在，产生一个新的Session
		s, err := this.newSession(location, w, r)
		if err != nil {
			AppLog.Add("in fnc StartSession," + err.Error())
		}
		return s, err
	} else {
		s, err := this.readSession(sid, location, w, r)
		if err != nil {
			AppLog.Add("in fnc StartSession," + err.Error())
		}
		return s, err
	}
}
func (this *SessionBase) getSessionId(r *http.Request) string {
	arrCookies := r.Cookies()
	var strValue string
	/*取最一个，该Cookies可能会被重置，最后一个才是最新SID*/
	for j := len(arrCookies) - 1; j >= 0; j-- {
		c := arrCookies[j]
		if c.Name == SessionKeyName {
			strValue = strings.TrimSpace(c.Value)
			break
		}
	}
	str, err := url.QueryUnescape(strValue)
	if err != nil {
		return strValue
	}
	return str
}

func (this *SessionBase) newSid() string {
	u := make([]byte, 16)
	rand.Read(u)
	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F
	return hex.EncodeToString(u)
}

func (this *SessionBase) setSessionName(sid string, w http.ResponseWriter, r *http.Request) {
	sid = url.QueryEscape(sid)
	cookie := &http.Cookie{
		Name:     SessionKeyName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
	}
	if Config.AppConfig.CookieDomain != "" {
		cookie.Domain = Config.AppConfig.CookieDomain
	}
	http.SetCookie(w, cookie)
	/*产生新的cookie后同时更新Request中的Session值，要么会出问题*/
	r.AddCookie(cookie)
}
func NewSession(ntype int) ISession {
	if ntype == 0 {
		return nil
	}
	switch ntype {
	case 1:
		return NewFileSession()
	case 2:
		return NewMemSession()
	case 3:
		return NewMysqlSession()
	default:
		return NewFileSession()
	}
}
