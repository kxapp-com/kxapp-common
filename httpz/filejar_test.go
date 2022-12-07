package httpz

//
//func TestEncodeCookies(t *testing.T) {
//	jar, _ := cookiejar.New(nil)
//	client := &http.Client{
//		Jar: jar,
//		Transport: &http.Transport{
//			//Proxy:           GetProxy(config.Proxy),
//			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
//		},
//	}
//	client.Get("https://www.baidu.com")
//	data := EncodeJar(jar)
//	jar2 := DecodeJar(data)
//	fmt.Printf("%v\n", jar)
//	fmt.Printf("%v\n", jar2)
//	cookies := GetCookiesByDomain("*", jar)
//	cookies2 := GetCookiesByDomain("*", jar2)
//	fmt.Printf("cook1 %v\n", cookies)
//	fmt.Printf("cook2 %v\n", cookies2)
//}
