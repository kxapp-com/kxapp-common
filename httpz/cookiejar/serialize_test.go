package cookiejar

import (
	"crypto/tls"
	"github.com/kxapp-com/kxapp-common/utilz"

	"net/http"
	"testing"
)

func TestEncodeCookies(t *testing.T) {

	utilz.IsInArray[string]("zhang", []string{"li"})
	jar, _ := New(nil)
	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			//Proxy:           GetProxy(config.Proxy),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	client.Get("https://www.baidu.com")
	//js1, _ := json.Marshal(jar.AllEntries())
	//os.WriteFile("all.json", js1, 0666)
	//b, e := jar.MarshalJSON()
	//os.WriteFile("per.json", b, 0666)
	//if e != nil {
	//	fmt.Printf("-----------------err%v", e)
	//}

	/*jar2, _ := New(nil)
	jar2.Unmarshal(b)
	fmt.Printf("newjar %v", jar2)*/
	/*	data := EncodeJar(jar)
		jar2 := DecodeJar(data)
		fmt.Printf("%v\n", jar)
		fmt.Printf("%v\n", jar2)
		cookies := GetCookiesByDomain("*", jar)
		cookies2 := GetCookiesByDomain("*", jar2)
		fmt.Printf("cook1 %v\n", cookies)
		fmt.Printf("cook2 %v\n", cookies2)*/
}
