package Web

import (
	"encoding/xml"
	"net/http"
)

type XmlResult struct {
	XmlText  string
	Response http.ResponseWriter
	CharSet  string
	Data     interface{}
}

func (this *XmlResult) ExecuteResult() error {
	if this.CharSet == "" {
		this.CharSet = "utf-8"
	}
	this.Response.Header().Add("Content-Type", "text/xml;charset="+this.CharSet)
	if this.XmlText == "" {
		if this.Data != nil {
			buf, err := xml.Marshal(this.Data)
			if err != nil {
				return err
			}
			this.XmlText = string(buf)
		}
	}
	this.Response.Write([]byte(this.XmlText))
	return nil
}
