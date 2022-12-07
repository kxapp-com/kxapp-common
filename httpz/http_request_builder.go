package httpz

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
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

/*
http请求过程中发生了网络错误，各种网络错误都归为状态码604
*/

var ProxyServiceURL = ""
var PrintDebug = true

/*
*
返回结果类，未对数据进行解码，是[]byte类型
*/
type HttpResponse struct {
	Status int
	Body   []byte
	Error  error //error类型是网络请求之类，io读取之类的错误，不包含逻辑错误
}

func (response *HttpResponse) HasError() bool {
	return response.Error != nil
}

// func NewHttpClient(jar http.CookieJar, proxyURL string) *http.Client {
func NewHttpClient(jar http.CookieJar) *http.Client {
	if ProxyServiceURL == "" {
		//return http.DefaultClient
		return &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	proxyURLUrl, err := url.Parse(ProxyServiceURL)
	if err != nil {
		log.Println(err)
	}
	return &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(proxyURLUrl),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
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

	request, _ := http.NewRequest(builder.method, builder.url, reader)
	for k, v := range builder.headers {
		k = strings.ReplaceAll(k, "\n", "")
		v = strings.ReplaceAll(v, "\n", "")
		request.Header.Add(k, v)
		//request.Header.Set(k, v)
	}
	for _, cookie := range builder.cookies {
		request.AddCookie(cookie)
	}
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
		b := request.Header.Get("scnt") == "AAAA-jhEMkUyNzQyRkQ5RDY4RUEwMEUyNzQyMjBDNjY1Q0M3OTZGN0U1RTkwNEFDMzMwNTI4NEVFMENFM0JGRUFDNDlGNEUzM0U5MDA3QUY3RDE4MzkyN0NFRjEzNEMwMTkyRkZEQjA5NUEzRDMwRjIxODkzRTdEQjZFRjI2NTBEM0U0MDk4NzU4QTU5NDEyNTQ1NkM4RDQ0NDVEMzk2QUNFMkNERjg0OTQ2NDAxQkY2QjE3NDVDREQyMDY1NDlCQ0Y5MkQwRjIwRDNEOUQ2RjIwQzFGRTNFRkU5OUI2MTM0NUMzQjM5OUM1NDVEREM0Q0FCN3wxAAABhOraKFyZV7hqPGvMsUnGVIOws4bzm5fkTKtPS971pkOiwhRQhX5Xzwsp_SHyAAggu_p2IIWNKnCyMeAWSixB3cahYCRsLKeB7ccNLHN9jx3zGI_HTw"
		v := request.Header.Get("X-Apple-ID-Session-Id") == "8D2E2742FD9D68EA00E274220C665CC796F7E5E904AC3305284EE0CE3BFEAC49F4E33E9007AF7D183927CEF134C0192FFDB095A3D30F21893E7DB6EF2650D3E4098758A594125456C8D4445D396ACE2CDF84946401BF6B1745CDD206549BCF92D0F20D3D9D6F20C1FE3EFE99B61345C3B399C545DDC4CAB7"
		log.Printf("scnt %v   session %v\n", b, v)
		log.Debugf("\n request--------------- %s %s \n request headers:%s\n", request.Method, request.URL, string(headers))
		if builder.body != nil {
			if b, ok := builder.body.([]byte); ok {
				log.Debugf("Body %s\n", string(b))
			} else {
				log.Debugf("Body %+v\n", builder.body)
			}
		}
	}
	response, httpError := httpClient.Do(request)
	if response == nil || httpError != nil {
		if PrintDebug {
			log.Debugf("response nil error %v  \n", httpError)
		}
		return &HttpResponse{Error: httpError}
	}
	if PrintDebug && response != nil {
		log.Debugf("response status %v  headers:%v\n", response.StatusCode, response.Header)
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	if builder.afterAction != nil {
		builder.afterAction(response)
	}

	if response.Body == nil {
		return &HttpResponse{Status: response.StatusCode}
	} else {
		body, ioReadError := io.ReadAll(response.Body)
		if PrintDebug {
			log.Debugf("Body %v \n", string(body))
		}
		if ioReadError != nil {
			return &HttpResponse{Error: ioReadError}
		}
		return &HttpResponse{Status: response.StatusCode, Body: body}
	}
}