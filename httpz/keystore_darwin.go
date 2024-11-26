package httpz

/*
#cgo LDFLAGS: -framework Security -framework CoreFoundation
#include <stdlib.h>
#include <stdio.h>
#include <Security/SecCertificate.h>
#include <Security/SecItem.h>
#include <sys/sysctl.h>
#include <errno.h>

int uninstall_cert(char *authorityName)
{
	CFStringRef str = CFStringCreateWithCString(NULL,
		authorityName, kCFStringEncodingUTF8);
	if (!str) {
		fprintf(stderr, "Error CFStringCreateWithCString");
		return 1;
	}

	const int nkeys = 3;
	CFStringRef keys[nkeys] = {
		kSecClass,
		kSecMatchSubjectWholeString,
		kSecMatchLimit,
	};
	CFTypeRef values[nkeys] = {
		kSecClassCertificate,
		str,
		kSecMatchLimitAll,
	};
	CFDictionaryRef query = CFDictionaryCreate(NULL,
		(const void **)keys, (const void **)values, nkeys,
		NULL, NULL);
	if (!query) {
		fprintf(stderr, "Error CFDictionaryCreate");
		CFRelease(str);
		return 1;
	}

	OSStatus s = SecItemDelete(query);

	CFRelease(query);
	CFRelease(str);
	return s;
}
*/
import "C"

var re = regexp.MustCompile(`(?m)^\([0-9]+\)(.+)`)
var relativeProfilesPath = "Firefox/Profiles"

func VerifyCert(certPath string) (err error) {
	// TODO: use Go API once supported
	// current api doesn't support reloading
	// https://github.com/golang/go/issues/46287
	_, err = runCommand("security", "verify-cert", "-c", certPath)
	return err
}

// Supported whether auto configuration
// is supported for this build
func Supported() bool {
	return true
}

func getProxyStatusForService(service, autoURL string) (ProxyStatus, error) {
	out, err := runCommand("networksetup", "-getautoproxyurl", service)
	if err != nil {
		return ProxyStatusConflict, err
	}

	url, enabled := parseGetAutoURL(string(out))
	if !enabled {
		return ProxyStatusNone, nil
	}

	if url == "" {
		return ProxyStatusNone, nil
	}

	if !equalURL(url, autoURL) {
		return ProxyStatusConflict, nil
	}

	return ProxyStatusInstalled, nil
}

func UninstallAutoProxy(autoURL string) {
	services, err := getNetworkServices()
	if err != nil {
		return
	}

	for _, service := range services {
		status, err := getProxyStatusForService(service, autoURL)
		if err != nil {
			continue
		}

		if status == ProxyStatusInstalled {
			_, _ = runCommand("networksetup", "-setautoproxystate", service, "off")
		}
	}
}

func InstallAutoProxy(autoURL string) error {
	services, err := getNetworkServices()
	if err != nil {
		return fmt.Errorf("failed reading network services")
	}

	var lastErr error
	installed := 0

	for _, service := range services {
		status, err := getProxyStatusForService(service, autoURL)
		if err != nil {
			lastErr = fmt.Errorf("failed checking proxy status")
			continue
		}
		if status == ProxyStatusConflict {
			lastErr = fmt.Errorf("auto configuration failed your OS has existing proxy settings")
		}
		if status != ProxyStatusNone {
			continue
		}
		if _, err = runCommand("networksetup", "-setautoproxyurl", service, autoURL); err != nil {
			lastErr = fmt.Errorf("failed configuring proxy make sure your user account has permissions to change proxy settings")
			continue
		}
		if _, err = runCommand("networksetup", "-setautoproxystate", service, "on"); err != nil {
			lastErr = fmt.Errorf("failed changing proxy state for service %s", service)
			continue
		}

		installed++
	}

	if installed > 0 {
		return nil
	}

	return lastErr
}

func runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

func getNetworkServices() ([]string, error) {
	cmd := exec.Command("networksetup", "-listnetworkserviceorder")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed reading network services: %v", err)
	}

	var services []string
	for _, match := range re.FindAllStringSubmatch(string(out), -1) {
		if len(match) != 2 {
			continue
		}
		services = append(services, strings.TrimSpace(match[1]))
	}

	return services, nil
}

func InstallCert(certPath string) error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	kp := path.Join(dir, "Library/Keychains/login.keychain")
	out, err := runCommand("security", "add-trusted-cert",
		"-p", "basic", "-p", "ssl", "-k", kp, certPath)

	if err == nil {
		return nil
	}

	if _, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("failed installing cert: %s", string(out))
	}

	return fmt.Errorf("failed installing cert: %v", err)
}

func UninstallCert(certPath string) error {
	_, err := runCommand("security", "remove-trusted-cert", certPath)
	return err
}

// parses networksetup -getautoproxyurl <service>
func parseGetAutoURL(txt string) (url string, enabled bool) {
	lines := strings.Split(strings.TrimSpace(txt), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "URL:") {
			url = strings.TrimSpace(line[4:])
			continue
		}

		if strings.HasPrefix(line, "Enabled:") {
			line = strings.TrimSpace(line[8:])
			enabled = line != "No"
			continue
		}
	}

	return
}

const certMode = syscall.S_IFREG | 0644

func Command(command string, args ...string) *exec.Cmd {
	return exec.Command(command, args...)
}

func InstallCertV2(certPath string) error {
	var s syscall.Stat_t
	err := syscall.Lstat(certPath, &s)
	if err != nil {
		return err
	}

	if s.Mode != certMode {
		return fmt.Errorf("cert incorrect mode: %o", s.Mode)
	}

	cmd := Command("/usr/bin/security",
		"add-trusted-cert",
		"-d",
		"-k", "/Library/Keychains/System.keychain",
		"-r", "trustRoot",
		certPath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func UninstallCertV2(authorityName string) error {
	cAuthorityName := C.CString(authorityName)
	defer C.free(unsafe.Pointer(cAuthorityName))

	ret := C.uninstall_cert(cAuthorityName)
	if ret == C.errSecItemNotFound {
		return fmt.Errorf("sec item not found")
	} else if ret != 0 {
		return fmt.Errorf("uninstall_cert: %d", ret)
	}
	return nil
}
