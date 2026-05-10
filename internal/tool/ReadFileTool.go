package tool

import (
	"os"
	"strings"
)

type ReadFileTool struct{}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Description() string {
	return "读取一个文件,输入应为一个文件路径。"
}

func (t *ReadFileTool) Run(input string) (string, error) {
	filename := strings.TrimSpace(input) // 去掉前后空格

	context, err := os.ReadFile(filename) // 读取文件内容

	if err != nil {
		return "", err
	}

	return string(context), nil
}
