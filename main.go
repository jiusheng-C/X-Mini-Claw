package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin) // 读取数据的缓冲读取器

	fmt.Println("Mini Code Agent 已启动")
	fmt.Println("输入‘帮助’查看命令")

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("输入错误：", err)
			continue
		}

		input = strings.TrimSpace(input) // 去除字符串input首尾所有的空白字符

		if input == "" {
			continue
		}

		switch {
		case input == "退出":
			fmt.Println("再见!")
			return

		case input == "帮助":
			printHelp()

		case input == "列表":
			listdir(".") // 当前目录

		case strings.HasPrefix(input, "读 "):
			filename := strings.TrimSpace(strings.TrimPrefix(input, "读 "))
			readFile(filename)

		case strings.HasPrefix(input, "写 "):
			handleWrite(input)

		case strings.HasPrefix(input, "运行 "):
			command := strings.TrimSpace(strings.TrimPrefix(input, "运行 "))
			runcommand(command)

		default:
			fmt.Println("未知命令:", input)
		}
	}

}

func printHelp() {
	fmt.Println("命令:")
	fmt.Println("  列表")
	fmt.Println("  读取 <文件名>")
	fmt.Println("  写入 <文件名> <内容>")
	fmt.Println("  运行 <命令>")
	fmt.Println("  帮助")
	fmt.Println("  退出")
}

func listdir(path string) { // 其实就是为了读取目录并打印文件名
	files, err := os.ReadDir(path) // 读取目录
	if err != nil {
		fmt.Println("列表目录错误:", err)
		return
	}

	for _, file := range files { // 判断当前文件是目录还是文件
		if file.IsDir() {
			fmt.Println("[目录] ", file.Name())
		} else {
			fmt.Println("[文件]", file.Name())
		}
	}
}

// readFile 函数用于读取文件内容并打印
// 参数 filename: 要读取的文件名
func readFile(filename string) {
	// 使用 os.ReadFile 读取文件内容
	// content: 读取到的文件内容
	content, err := os.ReadFile(filename)
	// 检查读取文件时是否发生错误
	if err != nil {
		fmt.Println("读取文件错误:", err)
		return
	}

	// 将读取到的内容转换为字符串并打印
	fmt.Println(string(content))
}

// handleWrite 处理写入命令，将输入内容写入指定文件
// handleWrite 处理写入文件的请求
// 参数 input: 包含命令和参数的输入字符串
func handleWrite(input string) {
	// 使用 SplitN 函数将输入字符串按空格分割成最多三部分
	// 这是为了分离命令、文件名和内容
	parts := strings.SplitN(input, " ", 3) // 将输入字符串按空格分割成三部分
	// 检查分割后的部分数量是否足够
	// 如果少于3部分，说明输入格式不正确
	if len(parts) < 3 {
		fmt.Println("用法: 写 <文件名> <内容>")
		return
	}

	// 从分割后的部分中提取文件名和内容
	// parts[0] 是命令本身，parts[1] 是文件名，parts[2] 是要写入的内容
	filename := parts[1]
	content := parts[2]

	// 尝试将内容写入指定文件
	// 0644 是文件权限，表示所有者可读写，其他用户只读
	err := os.WriteFile(filename, []byte(content), 0644)
	// 检查是否有错误发生
	if err != nil {
		fmt.Println("写入文件错误:", err)
		return
	}

	// 写入成功后输出提示信息
	fmt.Println("写入成功:", filename)
}

// runcommand 函数用于执行系统命令并输出结果
// 参数:
//
//	command: 要执行的shell命令字符串
func runcommand(command string) {
	// 使用exec.Command创建一个命令对象，使用sh -c来执行命令
	cmd := exec.Command("sh", "-c", command)

	// 执行命令并获取标准输出和标准错误合并后的输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 如果命令执行出错，打印错误信息
		fmt.Println("命令错误:", err)
	}

	// 打印命令执行结果
	fmt.Println(string(output))
}
