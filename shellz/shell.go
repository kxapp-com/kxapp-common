package shellz

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

// Import necessary Windows API functions
var (
	modshell32       = syscall.NewLazyDLL("shell32.dll")
	procShellExecute = modshell32.NewProc("ShellExecuteW")
)

// ShellExecuteW is a wrapper for the Windows ShellExecute function
func ShellExecute(hwnd uintptr, lpOperation, lpFile, lpParameters, lpDirectory *uint16, nShowCmd int32) uintptr {
	ret, _, _ := procShellExecute.Call(
		hwnd,
		uintptr(unsafe.Pointer(lpOperation)),
		uintptr(unsafe.Pointer(lpFile)),
		uintptr(unsafe.Pointer(lpParameters)),
		uintptr(unsafe.Pointer(lpDirectory)),
		uintptr(nShowCmd),
	)
	return ret
}

func ShellExecuteAdmin(script string) error {
	lpOperation, _ := syscall.UTF16PtrFromString("runas")
	lpFile, _ := syscall.UTF16PtrFromString("cmd")
	lpParameters, _ := syscall.UTF16PtrFromString("/C " + script)
	lpDirectory, _ := syscall.UTF16PtrFromString("")

	// Execute the command with elevated privileges
	ret := ShellExecute(0, lpOperation, lpFile, lpParameters, lpDirectory, syscall.SW_HIDE)
	if ret > 32 {
		return nil
	}
	return errors.New(getErrorMessage(ret))
}
func getErrorMessage(code uintptr) string {
	switch code {
	case 0:
		return "The operating system is out of memory or resources."
	case 2:
		return "The specified file was not found."
	case 3:
		return "The specified path was not found."
	case 5:
		return "Access is denied."
	case 8:
		return "Not enough memory to complete the operation."
	case 10:
		return "The environment is incorrect."
	case 11:
		return "The .exe file is invalid (non-Win32 .exe or error in .exe image)."
	case 26:
		return "A sharing violation occurred."
	case 27:
		return "The filename association is incomplete or invalid."
	case 28:
		return "The DDE transaction could not be completed because the request timed out."
	case 29:
		return "The DDE transaction failed."
	case 30:
		return "The DDE transaction could not be completed because other DDE transactions were being processed."
	case 31:
		return "There is no application associated with the given file extension."
	case 32:
		return "The specified dynamic-link library was not found."
	default:
		return fmt.Sprintf("Unknown error code: %d", code)
	}
}

// 以管理员身份运行创建link
func CreateLinkWindows(src string, dest string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	script := fmt.Sprintf("mklink  %s %s", dest, src)
	if info.IsDir() {
		script = fmt.Sprintf("mklink /c %s %s", dest, src)
	}
	return ShellExecuteAdmin(script)
}
