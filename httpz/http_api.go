package httpz

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kxapp-com/kxapp-common/cryptoz"
	"io"
	"net/http"
	"net/url"
	"time"
)

/*
*
urls包括多个接口服务器地址，循环请求，直到有一个成功
返回的结果如果不是json，则返回解码错误的信息,如果是json，但是json内status不是0，则返回json内的错误提示作为error的消息
否则返回json的data字段作为any
*/
func CallAuApiService(urls []string, funcName string, params map[string]any) (any, error) {
	params["tm"] = time.Now().UnixMilli()
	jsonBytes, e1 := json.Marshal(params)
	if e1 != nil {
		return nil, e1
	}
	fmt.Printf("request data %v \n", string(jsonBytes))
	basedParams := cryptoz.EncryptAndEncode(jsonBytes, "text/plain, */*")
	formdata := make(url.Values)
	formdata["func"] = []string{funcName}
	formdata["data"] = []string{basedParams}

	httpClient := http.DefaultClient
	var resp *http.Response
	var e3 error
	for _, ur := range urls {
		resp, e3 = httpClient.PostForm(ur, formdata)
		if resp != nil {
			defer resp.Body.Close()
		}
		if e3 == nil {
			break
		}
	}
	if e3 != nil {
		return nil, e3
	}
	responseData, e4 := io.ReadAll(resp.Body)
	if e4 != nil {
		return nil, e4
	}
	bodyData, e5 := cryptoz.DecodeAndDecrypt(string(responseData), "application/json,")
	if e5 != nil {
		return nil, errors.New(string(responseData))
	}
	var bodyJson = make(map[string]any)
	e := json.Unmarshal(bodyData, &bodyJson)
	if e != nil {
		return bodyData, e
	}
	if bodyJson["status"].(float64) != 0 {
		return nil, errors.New(bodyJson["data"].(string))
	} else {
		return bodyJson["data"], nil
	}
}
