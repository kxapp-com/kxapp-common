package errorz

import (
	"fmt"
	//"github.com/kxapp/apple-service/pkg/httpz"
	"net/http"
	"strings"
)

const (
	StatusNetworkError   = 6000
	StatusParseDataError = 6001
	StatusInternalError  = 6002
	StatusParamError     = 6003
	StatusIOError        = 6004
	StatusExpiredError   = 6005
	StatusSuccess        = 0
)
const StatusNameSuccess = "success"

func StatusText(code int) string {
	text := http.StatusText(code)
	if text != "" {
		text = strings.ReplaceAll(text, " ", "")
		return text
	}
	switch code {
	case StatusNetworkError:
		return "NetworkError"
	case StatusInternalError:
		return "InternalError"
	case StatusParseDataError:
		return "ReadDataError"
	case StatusParamError:
		return "StatusParamError"
	case StatusIOError:
		return "StatusIOError"
	case StatusExpiredError:
		return "StatusExpiredError"
	case StatusSuccess:
		return StatusNameSuccess
	default:
		return "UnknownError"
	}
}

type StatusError struct {
	/**
	用于标记错误编码，可能是来自于http的状态码，也可能是来源于业务逻辑返回结果中的状态码
	*/
	Status int `json:"status"`
	/**
	错误提醒消息
	*/
	//Message string `json:"message"`
	/**
	对于网络请求，为了方便，包含整个body内容
	*/
	Body string `json:"body"`
	//Body []byte `json:"body"`
}

func (e *StatusError) AsStatusResult() *StatusResult {
	return &StatusResult{Status: e.Status, Body: e.Body, StatusName: StatusText(e.Status)}
}
func (e *StatusError) StatusName() string {
	return StatusText(e.Status)
}
func (e *StatusError) Error() string {
	message := StatusText(e.Status)
	return fmt.Sprintf("status:%d message:%s body: %s", e.Status, message, string(e.Body))
	//return fmt.Sprintf("status:%d message:%s body: %s", e.Status, e.Message, string(e.Body))
}

func NewParseDataError(e ...error) *StatusError {
	for _, err := range e {
		if err != nil {
			return &StatusError{Status: StatusParseDataError, Body: err.Error()}
		}
	}
	return nil
}
func NewNetworkError(e error) *StatusError {
	if e == nil {
		return nil
	}
	return &StatusError{Status: StatusNetworkError, Body: e.Error()}
}
func NewInternalError(e string) *StatusError {
	return &StatusError{Status: StatusInternalError, Body: e}
}
func NewParamError(e string) *StatusError {
	return &StatusError{Status: StatusParamError, Body: e}
}
func NewIOError(e error) *StatusError {
	return &StatusError{Status: StatusIOError, Body: e.Error()}
}
func NewUnauthorizedError(body string) *StatusError {
	return &StatusError{Status: http.StatusUnauthorized, Body: body}
}

type StatusResult struct {
	Status     int    `json:"status"`      //0表示成功，其他表示失败
	StatusName string `json:"status_name"` //对status的解析说明的提示消息
	Body       any    `json:"body"`        //成功存放结果数据，失败存放消息提醒
}

func SuccessStatusResult(body any) *StatusResult {
	return &StatusResult{Status: 0, StatusName: StatusNameSuccess, Body: body}
}
