package Web

import (
	"net/http"
)

type JavaScriptResult struct {
	Script   string
	Response http.ResponseWriter
	CharSet  string
}

func (this *JavaScriptResult) ExecuteResult() error {
	if this.CharSet == "" {
		this.Response.Header().Add("Content-Type", "application/x-javascript;charset=utf-8")
	} else {
		this.Response.Header().Add("Content-Type", "application/x-javascript;charset="+this.CharSet)
	}

	this.Response.Write([]byte(this.Script))
	return nil
}
