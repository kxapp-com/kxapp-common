package httpz

// GetProxy 获取 macOS 系统代理设置
func GetProxy() (string, bool, string, error) {
	// 获取当前活动的网络服务
	cmd := exec.Command("networksetup", "-listallnetworkservices")
	output, err := cmd.Output()
	if err != nil {
		return "", false, "", fmt.Errorf("error getting network services: %v", err)
	}

	// 查找活动网络服务
	services := strings.Split(string(output), "\n")
	var activeService string
	for _, service := range services {
		if strings.Contains(service, "Wi-Fi") { // 假设 Wi-Fi 是活动的网络
			activeService = strings.TrimSpace(service)
			break
		}
	}

	if activeService == "" {
		return "", false, "", fmt.Errorf("no active network service found")
	}

	// 获取该服务的代理设置
	cmd = exec.Command("networksetup", "-getwebproxy", activeService)
	output, err = cmd.Output()
	if err != nil {
		return "", false, "", fmt.Errorf("error getting proxy settings: %v", err)
	}

	// 解析代理设置
	proxyEnabled := false
	proxyServer := ""
	proxyOverride := ""

	// 检查代理启用状态
	if strings.Contains(string(output), "Enabled: Yes") {
		proxyEnabled = true
		// 查找代理服务器
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Server:") {
				proxyServer = strings.TrimSpace(strings.Split(line, ":")[1])
			}
			if strings.Contains(line, "Bypass") {
				proxyOverride = strings.TrimSpace(strings.Split(line, ":")[1])
			}
		}
	}

	return proxyServer, proxyEnabled, proxyOverride, nil
}

// macOS 还原代理设置
func ResetProxy() error {
	cmd := "networksetup -setwebproxystate Wi-Fi off"
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
	cmd := fmt.Sprintf("networksetup -setwebproxy Wi-Fi %s %s", proxyAddress, "on")
	if proxyAddress == "" {
		cmd = "networksetup -setwebproxystate Wi-Fi off"
	}

	// 设置代理
	err := executeCommand(cmd)
	if err != nil {
		return fmt.Errorf("error setting proxy: %v", err)
	}

	// 设置跳过代理的域名
	if skipDomains != "" {
		cmd := fmt.Sprintf("networksetup -setproxybypassdomains Wi-Fi %s", skipDomains)
		err = executeCommand(cmd)
		if err != nil {
			return fmt.Errorf("error setting proxy bypass domains: %v", err)
		}
	}

	fmt.Println("Proxy settings updated successfully.")
	return nil
}
