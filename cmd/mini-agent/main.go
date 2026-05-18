package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jiusheng-C/X-Mini-Claw/internal/agent"
	"github.com/jiusheng-C/X-Mini-Claw/internal/llm"
	"github.com/jiusheng-C/X-Mini-Claw/internal/tool"
)

type Config struct {
	APIKey                string `json:"api_key"`
	BaseURL               string `json:"base_url"`
	Model                 string `json:"model"`
	APIStyle              string `json:"api_style"`
	Workspace             string `json:"workspace"`
	CommandTimeoutSeconds int    `json:"command_timeout_seconds"`
}

func main() {
	reader := bufio.NewReader(os.Stdin) // 读取数据的缓冲读取器
	registry := newToolRegistry()       // 创建工具注册表

	cfg, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("加载配置失败：", err)
		return
	}
	fmt.Println("配置加载成功，模型：", cfg.Model)

	client := llm.NewOpenAIClient(cfg.APIKey, cfg.BaseURL, cfg.Model, cfg.APIStyle) // 负责和模型通信
	miniAgent := agent.NewAgent(client, *registry)                                  // 负责"让模型决定用哪个工具，并执行工具循环"

	fmt.Println("Mini Code Agent 已启动")
	fmt.Println("输入‘帮助’查看命令")

	for {
		fmt.Print("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("\n再见!")
				return
			}
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

		case input == "查看当前列表":
			printResult(registry.Execute("list_dir", "."))
		case strings.HasPrefix(input, "查看列表 "):
			path := strings.TrimSpace(strings.TrimPrefix(input, "查看列表 "))
			printResult(registry.Execute("list_dir", path))

		case strings.HasPrefix(input, "读取 "):
			filename := strings.TrimSpace(strings.TrimPrefix(input, "读取 "))
			printResult(registry.Execute("read_file", filename))

		case strings.HasPrefix(input, "写入 "):
			handleWrite(registry, input)

		case strings.HasPrefix(input, "运行 "):
			command := strings.TrimSpace(strings.TrimPrefix(input, "运行 "))
			printResult(registry.Execute("exec_command", command))

		case input == "工具列表":
			printResult(registry.Execute("list_tools", ""))

		default:
			output, err := miniAgent.Run(input)
			if err != nil {
				fmt.Println("错误：", err)
				continue
			}
			fmt.Println(output)
			// 如果输入是定义的本地命令（帮助/读取/写入/运行...），走原来的分支;
			// 否则把用户输入交给 miniAgent，让模型决定调用哪个工具，最后输出答案。
		}
	}
}

func loadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func printHelp() {
	fmt.Println("命令:")
	fmt.Println("  查看当前列表")
	fmt.Println("  查看列表 <目录路径>")
	fmt.Println("  读取 <文件名>")
	fmt.Println("  写入 <文件名> <内容>")
	fmt.Println("  运行 <命令>")
	fmt.Println("  工具列表")
	fmt.Println("  帮助")
	fmt.Println("  退出")
}

// handleWrite 处理写入命令，将输入内容写入指定文件
// handleWrite 处理写入文件的请求
// 参数 input: 包含命令和参数的输入字符串
func handleWrite(registry *tool.Registry, input string) {
	parts := strings.SplitN(input, " ", 3)
	if len(parts) < 3 {
		fmt.Println("用法: 写入 <文件名> <内容>")
		return
	}

	printResult(registry.Execute("write_file", parts[1]+"|"+parts[2]))
}

func newToolRegistry() *tool.Registry {
	registry := tool.NewRegistry()
	registry.Register(&tool.ReadFileTool{})
	registry.Register(&tool.WriteFile{})
	registry.Register(&tool.ListDirTool{})
	registry.Register(tool.NewListTools(registry))
	registry.Register(&tool.ExecCommandTool{Timeout: 10 * time.Second})
	return registry
}

func printResult(output string, err error) { // 打印程序的执行结果或错误信息
	if err != nil {
		fmt.Println("错误:", err)
	}
	if output != "" {
		fmt.Println(output)
	}
}
