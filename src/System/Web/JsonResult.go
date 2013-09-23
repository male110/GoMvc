package Web

import (
	"encoding/json"
	"net/http"
)

type JsonResult struct {
	JsonText string
	Response http.ResponseWriter
	CharSet  string
	Data     interface{}
}

func (this *JsonResult) ExecuteResult() error {
	if this.CharSet == "" {
		this.Response.Header().Add("Content-Type", "application/json;charset=utf-8")
	} else {
		this.Response.Header().Add("Content-Type", "application/json;charset="+this.CharSet)
	}

	if this.JsonText == "" {
		if this.Data != nil {
			buf, err := json.Marshal(this.Data)
			if err != nil {
				return err
			}
			this.JsonText = string(buf)
		}
	}

	this.Response.Write([]byte(this.JsonText))
	return nil
}
