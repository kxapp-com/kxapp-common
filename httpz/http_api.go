package httpz

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kxapp-com/kxapp-common/cryptoz"
	"io"
	"net"
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
func CallAuApiService(urls []string, funcName string, params map[string]any, encodePassword string, decodePassword string) (any, error) {
	params["tm"] = time.Now().UnixMilli()
	jsonBytes, e1 := json.Marshal(params)
	if e1 != nil {
		return nil, e1
	}
	fmt.Printf("request data %v \n", string(jsonBytes))
	basedParams := cryptoz.EncryptAndEncode(jsonBytes, encodePassword)
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
	bodyData, e5 := cryptoz.DecodeAndDecrypt(string(responseData), decodePassword)
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

/*
*
获取本机ip地址，如果有对外ip，则会优先对外ip
没有对外ip，则找到私有ip地址中的一个返回
如果一个私有ip也没有，或者发生错误则返回127.0.0.1
*/
func GetIpAddress() string {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err == nil && conn != nil { //通过对外请求来找到真实的对外的ip地址，这个地址是最准确的地址，避免虚拟机之类的地址
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		ip := localAddr.IP.String()
		conn.Close()
		return ip
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.IsPrivate() && ipnet.IP.To4() != nil {
			return ipnet.IP.String()
		}
	}
	return "127.0.0.1"
}

/*
*
获取mac地址，如果出错则返回"",如果找到了出口ip，则根据出口ip找到对应的mac地址，如果没找到出口ip，则随便找一个私有ip对应的mac地址
*/
func GetMacAddress() string {
	is, e := net.Interfaces()
	if e != nil {
		return ""
	}
	ip := GetIpAddress()
	var macOK string
	for _, i := range is {
		mac := i.HardwareAddr.String()
		if mac == "" {
			continue
		}
		if ip == "127.0.0.1" { //mac 不为空，ip为空，随便获取一个mac返回
			return mac
		}
		addrs, e2 := i.Addrs()
		if e2 == nil {
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && ipnet.IP.To4() != nil {
					if ip == ipnet.IP.To4().String() { //找到ip地址是出口地址的mac地址返回
						return mac
					}
					if !ipnet.IP.IsLoopback() && ipnet.IP.IsPrivate() && ipnet.IP.To4() != nil {
						macOK = mac
					}
				}
			}
		}
	}
	return macOK
}
