package httpz

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

// Windows 获取代理设置
func GetProxy() (string, bool, string, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.QUERY_VALUE)
	if err != nil {
		return "", false, "", fmt.Errorf("error opening registry key: %v", err)
	}
	defer key.Close()

	proxyEnable, _, err := key.GetIntegerValue("ProxyEnable")
	if err != nil {
		return "", false, "", fmt.Errorf("error reading ProxyEnable: %v", err)
	}

	if proxyEnable == 0 {
		return "", false, "", nil
	}

	proxyServer, _, err := key.GetStringValue("ProxyServer")
	if err != nil {
		return "", false, "", fmt.Errorf("error reading ProxyServer: %v", err)
	}

	proxyOverride, _, err := key.GetStringValue("ProxyOverride")
	if err != nil {
		return "", false, "", fmt.Errorf("error reading ProxyOverride: %v", err)
	}

	return proxyServer, true, proxyOverride, nil
}

// Windows 还原代理设置
func ResetProxy() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("error opening registry key: %v", err)
	}
	defer key.Close()

	err = key.SetDWordValue("ProxyEnable", 0)
	if err != nil {
		return fmt.Errorf("error disabling ProxyEnable: %v", err)
	}

	err = key.SetStringValue("ProxyServer", "")
	if err != nil {
		return fmt.Errorf("error clearing ProxyServer: %v", err)
	}

	err = key.SetStringValue("ProxyOverride", "")
	if err != nil {
		return fmt.Errorf("error clearing ProxyOverride: %v", err)
	}

	fmt.Println("Proxy settings reset successfully.")
	return nil
}

// Windows 代理设置
func SetProxy(proxyAddress, skipDomains string) error {
	proxyAddress = strings.TrimSpace(proxyAddress)
	skipDomains = strings.TrimSpace(skipDomains)
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("error opening registry key: %v", err)
	}
	defer key.Close()

	var enableProxy uint32 = 0
	if proxyAddress != "" {
		enableProxy = 1
	} else {
		enableProxy = 0
		skipDomains = "" // 禁用代理时清空不使用代理的地址
	}

	err = key.SetDWordValue("ProxyEnable", enableProxy)
	if err != nil {
		return fmt.Errorf("error setting ProxyEnable: %v", err)
	}

	err = key.SetStringValue("ProxyOverride", skipDomains)
	if err != nil {
		return fmt.Errorf("error setting ProxyOverride: %v", err)
	}

	if proxyAddress != "" {
		err = key.SetStringValue("ProxyServer", proxyAddress)
		if err != nil {
			return fmt.Errorf("error setting ProxyServer: %v", err)
		}
	} else {
		err = key.SetStringValue("ProxyServer", "")
		if err != nil {
			return fmt.Errorf("error clearing ProxyServer: %v", err)
		}
	}

	fmt.Println("Proxy settings updated successfully.")
	return nil
}
