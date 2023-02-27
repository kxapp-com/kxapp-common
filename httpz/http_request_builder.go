package httpz

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"io"
	"net"
	"net/http/cookiejar"
	"regexp"
	"time"

	//"log"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

const UserAgent_GoogleChrome = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"
const UserAgent_AKD = "akd/1.0 CFNetwork/808.2.16 Darwin/16.3.0"
const UserAgent_XCode = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko)"
const UserAgent_XCode_Simple = "Xcode"
const ContentType_JSON = "application/json"
const ContentType_Plist = "text/x-xml-plist"
const ContentType_Form_URL = "application/x-www-form-urlencoded"
const ContentType_VND_JSON = "application/vnd.api+json"
const AcceptType_JSON = "application/json, text/plain, */*"

/*
http请求过程中发生了网络错误，各种网络错误都归为状态码604
*/

// var ProxyServiceURL = ""
var PrintDebug = true

/*
*
返回结果类，未对数据进行解码，是[]byte类型
*/
type HttpResponse struct {
	Status         int
	Body           []byte
	Header         http.Header
	Error          error //error类型是网络请求之类，io读取之类的错误，不包含逻辑错误
	responseClosed *http.Response
}

func (response *HttpResponse) HasError() bool {
	return response.Error != nil
}
func (response *HttpResponse) CookieValue(name string) string {
	c := response.Cookie(name, false)
	if c != nil {
		return c.Value
	}
	return ""
}
func (response *HttpResponse) Cookie(name string, reg bool) *http.Cookie {
	if !response.HasError() {
		for _, cookie := range response.responseClosed.Cookies() {
			if !reg {
				if name == cookie.Name {
					return cookie
				}
			} else {
				if mt, e := regexp.MatchString(name, cookie.Name); e == nil && mt {
					return cookie
				}
			}
		}
	}
	return nil
}
func defaultTransportDialContext(dialer *net.Dialer) func(context.Context, string, string) (net.Conn, error) {
	return dialer.DialContext
}
func NewHttpClient(jar http.CookieJar) *http.Client {
	if jar == nil {
		jar, _ = cookiejar.New(nil)
	}
	tra2 := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: defaultTransportDialContext(&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}),
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   30 * time.Second,
		ExpectContinueTimeout: 10 * time.Second,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	//tra := http.DefaultTransport
	return &http.Client{
		//Timeout:   180 * time.Second,
		Jar:       jar,
		Transport: tra2,
	}
}

// func NewHttpClient(jar http.CookieJar, proxyURL string) *http.Client {
func NewHttpClient1(jar http.CookieJar) *http.Client {
	//if ProxyServiceURL == "" {
	//return http.DefaultClient
	return &http.Client{
		//Timeout: 180 * time.Second,
		Jar: jar,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	//}
	//proxyURLUrl, err := url.Parse(ProxyServiceURL)
	//if err != nil {
	//	log.Println(err)
	//}
	//return &http.Client{
	//	Timeout: 120 * time.Second,
	//	Jar:     jar,
	//	Transport: &http.Transport{
	//		Proxy:           http.ProxyURL(proxyURLUrl),
	//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	},
	//}
}

type HttpRequestBuilder struct {
	//request     *http.Request
	afterAction func(response *http.Response)

	url     string
	method  string
	headers map[string]string
	body    any
	cookies []*http.Cookie
}

func Post(url string, headers map[string]string) *HttpRequestBuilder {
	return NewHttpRequestBuilder(http.MethodPost, url).AddHeaders(headers)
}
func Get(url string, headers map[string]string) *HttpRequestBuilder {
	return NewHttpRequestBuilder(http.MethodGet, url).AddHeaders(headers)
	//return NewHttpRequestBuilder(http.MethodPost, url).AddHeaders(headers)
}
func NewHttpRequestBuilder(method string, url string) *HttpRequestBuilder {
	return &HttpRequestBuilder{method: method, url: url, headers: map[string]string{}}
}

/*
*
string 则直接变成2进制，url.Values则编码为values方式提交，其他则先编码为json字符串后提交
*/
func (builder *HttpRequestBuilder) AddBody(structOrMap any) *HttpRequestBuilder {
	builder.body = structOrMap
	return builder
}
func (builder *HttpRequestBuilder) AddHeaders(headers map[string]string) *HttpRequestBuilder {
	for k, v := range headers {
		builder.headers[k] = v
	}
	return builder
}
func (builder *HttpRequestBuilder) ContentType(v string) *HttpRequestBuilder {
	return builder.SetHeader("Content-Type", v)
}
func (builder *HttpRequestBuilder) Accept(v string) *HttpRequestBuilder {
	return builder.SetHeader("Accept", v)
}
func (builder *HttpRequestBuilder) SetHeader(k string, v string) *HttpRequestBuilder {
	builder.headers[k] = v
	return builder
}

func (builder *HttpRequestBuilder) AddCookie(cookie *http.Cookie) *HttpRequestBuilder {
	builder.cookies = append(builder.cookies, cookie)
	return builder
}
func (builder *HttpRequestBuilder) BuildRequest() *http.Request {
	structOrMap := builder.body
	var reader io.Reader = nil

	if v, ok := structOrMap.(url.Values); ok {
		reader = strings.NewReader(v.Encode())
	} else if v, ok := structOrMap.(string); ok {
		reader = strings.NewReader(v)
	} else if v, ok := structOrMap.([]byte); ok {
		reader = bytes.NewReader(v)
	} else {
		if structOrMap != nil {
			buf := new(bytes.Buffer)
			json.NewEncoder(buf).Encode(structOrMap)
			reader = buf
		}
	}

	request, e := http.NewRequest(builder.method, builder.url, reader)
	if e != nil {
		log.Error(e)
		return nil
	}
	for k, v := range builder.headers {
		//k = strings.ReplaceAll(k, "\n", "")
		//	v = strings.ReplaceAll(v, "\n", "")
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		request.Header.Add(k, v)
		//request.Header.Set(k, v)
	}
	for _, cookie := range builder.cookies {
		request.AddCookie(cookie)
	}
	//request.Close = true
	return request
}

/*
*
在request函数返回前进行的处理行为，一般是获取返回的headers，cookies，更新逻辑业务中的一些东西
*/
func (builder *HttpRequestBuilder) BeforeReturn(action func(response *http.Response)) *HttpRequestBuilder {
	builder.afterAction = action
	return builder
}
func (builder *HttpRequestBuilder) Request(httpClient *http.Client) *HttpResponse {
	request := builder.BuildRequest()
	if PrintDebug {

		headers, e := json.Marshal(request.Header)
		if e != nil {
			log.Error("request header error %v", e)
		}
		log.Debugf("\n request--------------- %s %s \n request headers:%s", request.Method, request.URL, string(headers))
		if builder.body != nil {
			if b, ok := builder.body.([]byte); ok {
				if len(b) < 1024 {
					log.Debugf("Body %s\n", string(b))
				} else {
					log.Debugf("Body size >1024 ,log ignore \n")
				}

			} else {
				log.Debugf("Body %+v\n", builder.body)
			}
		}
	}

	response, httpError := httpClient.Do(request)
	if httpError != nil {
		log.Error("http request error %v  \n", httpError)
		if e, ok := httpError.(*url.Error); ok {
			return &HttpResponse{Error: e.Err, responseClosed: response}
		}
		return &HttpResponse{Error: errors.New("NetworkError"), responseClosed: response}
	}
	if PrintDebug && response != nil {
		log.Debugf("response status %v  headers:\n%v\n", response.StatusCode, response.Header)
	}

	if response.Body == nil {
		if builder.afterAction != nil {
			builder.afterAction(response)
		}
		return &HttpResponse{Status: response.StatusCode, Header: response.Header, responseClosed: response}
	} else {
		body, ioReadError := io.ReadAll(response.Body)
		if response != nil && response.Body != nil {
			defer response.Body.Close()
		}
		if builder.afterAction != nil {
			builder.afterAction(response)
		}
		if PrintDebug {
			log.Debugf("Body %v \n", string(body))
		}
		if ioReadError != nil {
			log.Error(ioReadError)
			return &HttpResponse{Error: ioReadError, Header: response.Header, responseClosed: response}
		}
		return &HttpResponse{Status: response.StatusCode, Body: body, Header: response.Header, responseClosed: response}
	}

}
