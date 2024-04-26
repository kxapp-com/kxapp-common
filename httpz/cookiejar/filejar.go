package cookiejar

import (
	"bytes"
	"encoding/gob"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"strings"
)

/*
*
gob encode all cookies in this jar
*/
func EncodeJar(jar *cookiejar.Jar) []byte {
	return EncodeCookies(GetCookiesByDomain("*", jar))
}

/*
new jar , and add all cookies in this  gob data  to this jar
*/
func DecodeJar(data []byte) *cookiejar.Jar {
	jar, _ := cookiejar.New(nil)
	cookies := DecodeCookies(data)
	for domain, cookies := range cookies {
		u, e := url.Parse(domain)
		if e != nil {
			continue
		}
		jar.SetCookies(u, cookies)
	}
	return jar
}

func EncodeCookies(cookieMap map[string][]*http.Cookie) []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	encoder.Encode(cookieMap)
	return buffer.Bytes()
}
func DecodeCookies(data []byte) map[string][]*http.Cookie {
	savedCookies := make(map[string][]*http.Cookie)
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(&savedCookies)
	return savedCookies
}
func GetAllCookies(jar *cookiejar.Jar) map[string][]*http.Cookie {
	return GetCookiesByDomain("*", jar)
}
func GetDomains(jar *cookiejar.Jar) map[string]bool {
	result := make(map[string]bool)
	jarStruct := reflect.ValueOf(jar)
	entries := jarStruct.Elem().FieldByName("entries")
	jarIterator := entries.MapRange()
	for jarIterator.Next() {
		//jarId := jarIterator.Key()
		jarEntry := jarIterator.Value()
		cookieIterator := jarEntry.MapRange()
		//fmt.Printf("cookie id %v jarEntry \n", jarId)
		for cookieIterator.Next() {
			cookieEntry := cookieIterator.Value()
			domain := cookieEntry.FieldByName("Domain")
			sec := cookieEntry.FieldByName("Secure")
			result[domain.String()] = sec.Bool()
			//fmt.Printf("domain %v sec %v ", domain, sec)
		}
	}
	return result
}
func GetURLList(jar *cookiejar.Jar) []string {
	result := make([]string, 0)
	domains := GetDomains(jar)
	for domain, _ := range domains {
		/*
			httpURL := "http://" + domain
			httpCookies := GetCookies(httpURL, jar)
			if httpCookies != nil {
				result = append(result, httpURL)
			}*/
		httpsURL := "https://" + domain
		httpsCookies := GetCookiesByURL(httpsURL, jar)
		if httpsCookies != nil {
			result = append(result, httpsURL)
		}
	}
	return result
}
func GetCookiesByDomain(domainReg string, jar *cookiejar.Jar) map[string][]*http.Cookie {
	domainReg = strings.TrimSpace(domainReg)
	saveJars := make(map[string][]*http.Cookie)
	domainSec := GetDomains(jar)
	for domain, _ := range domainSec {
		//reg, e := regexp.Compile(domainReg)
		//if e == nil && reg.MatchString(domain) {
		var match = domainReg == domain
		if domainReg == "*" {
			match = true
		} else if strings.HasPrefix(domainReg, "*.") {
			domainBody := strings.Replace(domainReg, "*.", "", 1)
			if strings.HasSuffix(domain, domainBody) {
				match = true
			}
		}
		if match {
			httpsURL := "https://" + domain
			httpsCookies := GetCookiesByURL(httpsURL, jar)
			if httpsCookies != nil {
				saveJars[httpsURL] = httpsCookies
			}
		}
	}
	return saveJars
}
func GetCookiesList(urls []string, jar *cookiejar.Jar) map[string][]*http.Cookie {
	saveJars := make(map[string][]*http.Cookie)
	for _, domain := range urls {
		cookies := GetCookiesByURL(domain, jar)
		if cookies != nil {
			saveJars[domain] = cookies
		}
	}
	return saveJars
}
func GetCookiesByURL(urlString string, jar *cookiejar.Jar) []*http.Cookie {
	urlString = strings.TrimSpace(urlString)
	httpURL, e := url.Parse(urlString)
	if e == nil {
		httpCookies := jar.Cookies(httpURL)
		if len(httpCookies) > 0 {
			return httpCookies
		} else {
			return nil
		}
	}
	return nil
}
