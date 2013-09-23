package Session

import (
	. "System/Log"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

/*
CREATE TABLE `session` (
	`session_id` CHAR(32) NULL,
	`session_data` BLOB NULL,
	`lastupdatetime` DATETIME NULL,
	PRIMARY KEY (`session_id`)
)
COLLATE='utf8_general_ci';*/
type MysqlSession struct {
	SessionBase
	gcing bool
}

func (this *MysqlSession) readSessionFromDB(sid string, location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	db, err := sql.Open("mysql", location)
	m := make(map[string]interface{})
	if err != nil {
		AppLog.Add("MysqlSession.readSession,打开数据连接时出错：" + err.Error() + ",连接字符串为：" + location)
		return m, err
	}
	defer db.Close()
	row, err := db.Query("select `session_data` from `session` where session_id=?", sid)
	if err != nil {
		AppLog.Add("MysqlSession.readSession,取数据时出错：" + err.Error())
		return m, err
	}
	if row.Next() {
		var buf []byte
		err = row.Scan(&buf)
		if err != nil {
			AppLog.Add("MysqlSession.readSession,从行里取数据时出错：" + err.Error())
			return m, err
		}
		if buf != nil && len(buf) > 0 {
			m, err = GobSerialize.Decode(buf)
			if err != nil {
				AppLog.Add("MysqlSession.readSession,gob.Decode时出错：" + err.Error())
			}
		}
		return m, err
	} else {
		return this.newSession(location, w, r)
	}
}

func (this *MysqlSession) CreateNewSession(location string, w http.ResponseWriter, r *http.Request) (map[string]interface{}, error) {
	sid := this.newSid()
	this.setSessionName(sid, w, r)
	m := make(map[string]interface{})

	db, err := sql.Open("mysql", location)
	if err != nil {
		AppLog.Add("MysqlSession.newSession,打开数据连接时出错：" + err.Error() + ",连接字符串为：" + location)
		return m, err
	}
	defer db.Close()
	_, err = db.Exec("insert into `session` (`session_id`,`session_data`,`lastupdatetime`) values(?,null,now())", sid)
	if err != nil {
		AppLog.Add("MysqlSession.newSession,插入数据时出错：" + err.Error())
	}
	return m, err
}
func (this *MysqlSession) EndSession(data map[string]interface{}, location string, r *http.Request) error {
	buf, err := GobSerialize.Encode(data)
	if err != nil {
		AppLog.Add("MysqlSession.EndSession,gob.Encode时出错：" + err.Error())
		return err
	}
	db, err := sql.Open("mysql", location)
	if err != nil {
		AppLog.Add("MysqlSession.newSession,打开数据连接时出错：" + err.Error() + ",连接字符串为：" + location)
		return err
	}
	defer db.Close()
	sid := this.getSessionId(r)
	if sid == "" {
		AppLog.Add("无法获取SessionID")
		return nil
	}
	_, err = db.Exec("update `session` set `session_data`=? where `session_id`=?", buf, sid)
	if err != nil {
		AppLog.Add("MysqlSession.newSession,插入数据时出错：" + err.Error())
	}
	return err
}
func (this *MysqlSession) GC(timeOut int, location string) {

	if this.gcing {
		return
	}
	this.gcing = true
	defer func() { this.gcing = false }()
	db, err := sql.Open("mysql", location)
	if err != nil {
		AppLog.Add("MysqlSession.GC,打开数据连接时出错：" + err.Error() + ",连接字符串为：" + location)
		return
	}
	strSql := "delete from `session` where  date_add(`lastupdatetime`, interval " + fmt.Sprintf("%v", timeOut) + " minute)" + "<=now()"
	_, err = db.Exec(strSql)
	if err != nil {
		AppLog.Add("MysqlSession.GC," + err.Error())
	}
}
func NewMysqlSession() *MysqlSession {
	s := &MysqlSession{gcing: false}
	s.readSession = s.readSessionFromDB
	s.newSession = s.CreateNewSession
	return s
}
