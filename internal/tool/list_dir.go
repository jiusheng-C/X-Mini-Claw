package tool

import (
	"os"
	"strings"
)

type ListDirTool struct{}

func (t *ListDirTool) Name() string {
	return "list_dir"
}

func (t *ListDirTool) Description() string {
	return "列出文件和目录，输入应为目录路径。"
}

func (t *ListDirTool) Run(input string) (string, error) {
	path := strings.TrimSpace(input) // 去除输入的前后空格
	if path == "" {                  // 如果path为空，则默认为当前目录
		path = "."
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}

	var builder strings.Builder // 创建一个字符串构建器

	for _, file := range files {
		if file.IsDir() {
			builder.WriteString("[目录] ") // 如果是目录，则在文件名前加上"[目录]"
		} else {
			builder.WriteString("[文件] ") // 如果是文件，则在文件名前加上"[文件]"
		}

		builder.WriteString(file.Name()) // 写入文件名
		builder.WriteRune('\n')
	}

	return builder.String(), nil
}
