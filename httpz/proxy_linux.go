package httpz

func GetProxy() (string, bool, string, error) {
	proxy := os.Getenv("http_proxy")
	if proxy == "" {
		return "", false, "", nil
	}

	noProxy := os.Getenv("no_proxy")
	return proxy, true, noProxy, nil
}

func ResetProxy() error {
	cmd := "unset http_proxy && unset https_proxy && unset no_proxy"
	err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("error resetting proxy: %v", err)
	}
	fmt.Println("Proxy settings reset successfully.")
	return nil
}
func SetProxy(proxyAddress, skipDomains string) error {
	proxyAddress = strings.TrimSpace(proxyAddress)
	skipDomains = strings.TrimSpace(skipDomains)
	var cmd string
	if proxyAddress != "" {
		cmd = fmt.Sprintf("export http_proxy=http://%s && export https_proxy=http://%s", proxyAddress, proxyAddress)
	} else {
		cmd = "unset http_proxy && unset https_proxy"
	}

	err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("error setting proxy: %v", err)
	}

	// 处理不使用代理的地址
	if skipDomains != "" {
		cmd = fmt.Sprintf("export no_proxy=%s", skipDomains)
		err = executeCommand(cmd)
		if err != nil {
			return fmt.Errorf("error setting no_proxy: %v", err)
		}
	}

	fmt.Println("Proxy settings updated successfully.")
	return nil
}
