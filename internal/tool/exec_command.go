package tool

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ExecCommandTool struct {
	Timeout time.Duration // 持续时间
}

func (t *ExecCommandTool) Name() string {
	return "exec_command"
}

func (t *ExecCommandTool) Description() string {
	return "执行一个shell命令, 输入应为命令字符串。"
}

func (t *ExecCommandTool) Run(input string) (string, error) {
	command := strings.TrimSpace(input)
	if command == "" {
		return "", fmt.Errorf("命令为空")
	}

	if isDangerousCommand(command) {
		return "", fmt.Errorf("危险指令已拦截: %s", command)
	}
	// ------- 以上为危险命令拦截部分

	timeout := t.Timeout
	if timeout <= 0 { // 如果没有设置或者设置不合理的超时时间
		timeout = 10 * time.Second // 默认超时10秒，time.second 表示1秒
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout) // 创建一个带有超时的上下文
	defer cancel()                                                    // 在函数退出时调用cancel函数，释放资源

	cmd := exec.CommandContext(ctx, "sh", "-c", command) // 将上下文与命令绑定
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", "[Console]::OutputEncoding = [Text.UTF8Encoding]::UTF8; "+command)
	}
	output, err := cmd.CombinedOutput() // 执行命令并获取输出

	if ctx.Err() == context.DeadlineExceeded { // 如果命令执行时间超过设定的超时时间
		return "", fmt.Errorf("命令已超时: %s", timeout)
	}
	// ------ 以上为超时控制部分

	if err != nil { // 如果命令执行失败
		return string(output), err
	}

	return string(output), nil // 输出命令执行结果
}

func isDangerousCommand(command string) bool {
	dangerousPatterns := []string{
		"rm -rf /",
		"rm -rf /*",
		"shutdown",
		"reboot",
		"mkfs",
		":(){:|:&};:",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(command, pattern) { // 如果在command中能找到pattern
			return true
		}
	}

	return false
}
