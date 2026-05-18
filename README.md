# X-Mini-Claw

X-Mini-Claw 是一个使用 Go 编写的本地命令行 Mini Code Agent。它可以通过大模型理解用户输入，并在本地调用已注册的工具完成文件读取、文件写入、目录查看、命令执行等操作。

项目目标是用较小的代码量实现一个简化版的代码智能体，帮助理解 Agent、LLM 调用、本地工具注册与工具执行循环的基本原理。

## 功能特性

- 命令行交互式运行
- 支持直接执行本地命令
- 支持读取文件内容
- 支持写入文件内容
- 支持查看目录列表
- 支持查看当前已注册工具
- 支持将普通自然语言请求交给大模型处理
- 支持 OpenAI Responses API 和 Chat Completions API 两种调用风格
- 命令执行带超时控制和基础危险命令拦截

## 项目结构

```text
.
├── cmd/
│   └── mini-agent/
│       └── main.go              # 程序入口，负责命令行交互和工具注册
├── internal/
│   ├── agent/
│   │   ├── agent.go             # Agent 主循环：请求模型、解析工具调用、执行工具
│   │   ├── prompt.go            # 构建系统提示词和工具描述
│   │   └── types.go             # Agent 相关类型定义
│   ├── config/
│   │   └── config.go            # 配置文件读取
│   ├── llm/
│   │   ├── client.go            # LLM 客户端接口定义
│   │   ├── fake.go              # 测试或模拟用客户端
│   │   └── openai.go            # OpenAI 风格 API 客户端实现
│   └── tool/
│       ├── Registry.go          # 工具注册表
│       ├── tool.go              # 工具接口定义
│       ├── ReadFileTool.go      # 读取文件工具
│       ├── write_file.go        # 写入文件工具
│       ├── list_dir.go          # 查看目录工具
│       ├── exec_command.go      # 执行 shell 命令工具
│       └── list_tools.go        # 查看工具列表工具
├── config.example.json          # 配置文件模板
├── config.json                  # 本地实际配置文件
└── go.mod                       # Go 模块文件
```

## 环境要求

- Go 1.25 或更高版本
- 一个兼容 OpenAI 接口的大模型 API Key

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/jiusheng-C/X-Mini-Claw.git
cd X-Mini-Claw
```

### 2. 准备配置文件

复制配置模板：

```bash
cp config.example.json config.json
```

然后编辑 `config.json`：

```json
{
  "api_key": "your-api-key",
  "base_url": "https://api.openai.com/v1",
  "model": "gpt-4o-mini",
  "api_style": "responses",
  "workspace": ".",
  "command_timeout_seconds": 10
}
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `api_key` | 大模型 API Key |
| `base_url` | API 基础地址，例如 `https://api.openai.com/v1` |
| `model` | 使用的模型名称 |
| `api_style` | API 风格，可选 `responses` 或 `chat_completions` |
| `workspace` | 工作目录配置，目前预留 |
| `command_timeout_seconds` | 命令超时时间，单位是秒 |

> 注意：`config.json` 中包含 API Key，不建议提交到公开仓库。

### 3. 运行项目

```bash
go run ./cmd/mini-agent
```

启动成功后会看到类似输出：

```text
配置加载成功，模型： gpt-4o-mini
Mini Code Agent 已启动
输入‘帮助’查看命令
>
```

## 内置命令

程序启动后，可以直接输入以下命令：

| 命令 | 作用 |
| --- | --- |
| `帮助` | 查看命令帮助 |
| `查看当前列表` | 查看当前目录下的文件和文件夹 |
| `查看列表 <目录路径>` | 查看指定目录下的文件和文件夹 |
| `读取 <文件名>` | 读取指定文件内容 |
| `写入 <文件名> <内容>` | 将内容写入指定文件 |
| `运行 <命令>` | 执行一条 shell 命令 |
| `工具列表` | 查看当前注册的全部工具 |
| `退出` | 退出程序 |

示例：

```text
> 查看当前列表
> 读取 go.mod
> 写入 hello.txt 你好，Mini Agent
> 运行 go test ./...
> 工具列表
> 退出
```

## 自然语言 Agent 模式

如果输入内容不匹配内置命令，程序会把用户输入交给 Mini Agent 处理。

Mini Agent 会：

1. 根据当前工具注册表生成工具说明
2. 构造系统提示词
3. 请求大模型返回 JSON 格式的工具调用
4. 执行对应工具
5. 将工具执行结果继续交给模型
6. 最终通过 `final_answer` 返回答案

模型返回的 JSON 格式如下：

```json
{
  "tool": "tool_name",
  "input": "tool_input"
}
```

当模型已经获得足够信息时，需要返回：

```json
{
  "tool": "final_answer",
  "input": "最终回答内容"
}
```

Agent 最多执行 5 轮工具调用，超过后会返回 `agent reached max steps` 错误。

## 当前已注册工具

| 工具名 | 功能 |
| --- | --- |
| `read_file` | 读取一个文件 |
| `write_file` | 写入一个文件 |
| `list_dir` | 列出目录内容 |
| `exec_command` | 执行 shell 命令 |
| `list_tools` | 输出当前工具列表 |
| `final_answer` | Agent 最终回答，不是真实注册工具，由 Agent 协议使用 |

## 实现原理简述

### 工具接口

所有本地工具都实现统一的 `Tool` 接口：

```go
type Tool interface {
    Name() string
    Description() string
    Run(input string) (string, error)
}
```

这样 Agent 不需要关心每个工具内部如何实现，只需要通过工具名找到工具并执行。

### 工具注册表

`Registry` 负责保存和调度工具：

- `Register`：注册工具
- `Execute`：根据工具名执行工具
- `ListTools`：列出所有已注册工具

### Agent 循环

Agent 的核心流程在 `internal/agent/agent.go` 中：

- 先把工具说明写入系统提示词
- 要求模型只返回 JSON
- 根据模型返回的 `tool` 字段调用本地工具
- 将工具结果放回上下文
- 如果模型返回 `final_answer`，则结束循环

### LLM 客户端

`internal/llm/openai.go` 实现了 OpenAI 风格的 HTTP 调用，支持：

- `responses`：请求 `/responses`
- `chat_completions`：请求 `/chat/completions`

## 安全说明

本项目包含本地命令执行能力，因此需要谨慎使用。

当前 `exec_command` 工具已经做了基础危险命令拦截，例如：

- `rm -rf /`
- `rm -rf /*`
- `shutdown`
- `reboot`
- `mkfs`
- fork bomb

但这只是基础保护，不能覆盖所有危险命令。实际使用时请避免让 Agent 操作重要目录或执行不可信命令。

## 开发与测试

格式化代码：

```bash
gofmt -w .
```

运行测试：

```bash
go test ./...
```

运行程序：

```bash
go run ./cmd/mini-agent
```

## 后续可改进方向

- 让 `workspace` 配置真正限制工具可访问目录
- 增加命令执行日志文件输出
- 增加更完善的命令安全策略
- 增加单元测试
- 增加多轮对话历史保存
- 支持更标准的 tool calling 协议
- 增加文件编辑、代码搜索等更实用的工具

## 许可证

当前项目暂未声明许可证。