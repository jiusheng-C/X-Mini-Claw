package tool

import (
	"fmt"
	"os"
	"strings"
)

type WriteFile struct{}

func (w *WriteFile) Name() string {
	return "write_file"
}

func (w *WriteFile) Description() string {
	return "将内容写入文件，输入格式：<文件名>|<内容>。"
}

func (w *WriteFile) Run(input string) (string, error) {
	parts := strings.SplitN(input, "|", 2) // 将输入按|分割为文件名和内容
	if len(parts) != 2 {
		return "", fmt.Errorf("输入格式无效，预期格式为: <文件名>|<内容>")
	}

	filename := parts[0]
	context := parts[1]

	err := os.WriteFile(filename, []byte(context), 0644) // 将内容写入文件
	if err != nil {
		return "", err
	}

	return "写入成功: " + filename, nil
}
