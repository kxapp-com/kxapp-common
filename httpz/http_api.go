package httpz

import (
	"net"
	"net/http"
	"time"
)

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
func GetFastestURL(urls ...string) string {
	type Result struct {
		URL      string
		Duration time.Duration
	}
	getResponseTime := func(url string, ch chan Result) {
		start := time.Now()
		resp, err := http.Head(url)
		if err != nil {
			ch <- Result{URL: url, Duration: time.Duration(0)}
			return
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		duration := time.Since(start)
		ch <- Result{URL: url, Duration: duration}
	}

	ch := make(chan Result, len(urls))
	for _, url := range urls {
		go getResponseTime(url, ch) // 启动并发请求
	}

	var fastest Result
	for i := 0; i < len(urls); i++ {
		result := <-ch
		if fastest.Duration == 0 || result.Duration < fastest.Duration {
			fastest = result
		}
	}

	return fastest.URL
}
