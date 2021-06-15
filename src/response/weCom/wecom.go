package weCom

import (
	"net/http"
)

type ResponseWeCom struct {
	ErrCode int    `json:"errcode"`
	ErrMSG  string `json:"errmsg"`
}

func (res *ResponseWeCom) GetBody() *http.ResponseWriter {
	return nil
}
func (res *ResponseWeCom) GetHeaders() *http.ResponseWriter {
	return nil
}

func (res *ResponseWeCom) GetStatusCode() int {
	return 200
}
