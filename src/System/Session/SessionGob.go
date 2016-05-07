package Session

import (
	. "System/Log"
	"bytes"
	"encoding/gob"
	"reflect"
)

type SessionGob struct {
	registerType map[string]bool
}

/*map[string]interface{},interface需要注册，要么没办法反序列化*/
func (this *SessionGob) registerItem(item interface{}) {
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			AppLog.Add("in SessionGob.registerItem()," + err.Error())
		}
	}()
	switch item.(type) {
	case int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64:
		return
	case uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64:
		return
	case float32, *float32, float64, *float64, string, *string:
		return
	default:
		//非内置类型，需要注册
		rv := reflect.ValueOf(item)
		strTypeName := rv.Type().String()
		_, ok := this.registerType[strTypeName]
		if ok {
			return
		}
		for rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		//直接注册
		gob.Register(rv.Interface())
		this.registerType[strTypeName] = true
	}
}

/*如果map的值有slice,array，其元素不能是指针类型*/
func (this *SessionGob) Encode(m map[string]interface{}) ([]byte, error) {
	for _, v := range m {
		this.registerItem(v)
	}
	
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(m)
	if err != nil {
		AppLog.Add("in SessionGob.Encode()," + err.Error())
		return nil, err
	}
	return buf.Bytes(), err
}
func (this *SessionGob) Decode(data []byte) (map[string]interface{}, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	out := make(map[string]interface{})
	err := dec.Decode(&out)
	
	return out, err
}

var GobSerialize = &SessionGob{make(map[string]bool)}
