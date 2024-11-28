package httpz

import (
	"os/exec"
	"strings"
)

// 执行命令
func executeCommand(command string) error {
	cmd := strings.Split(command, " ")
	return exec.Command(cmd[0], cmd[1:]...).Run()
}
