package agent

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jiusheng-C/X-Mini-Claw/internal/llm"
	"github.com/jiusheng-C/X-Mini-Claw/internal/tool"
)

type Agent struct {
	client   llm.Client     // 负责"问模型"
	registry *tool.Registry // 负责"执行模型选择的工具"
}

func NewAgent(client llm.Client, registry tool.Registry) *Agent {
	return &Agent{
		client:   client,
		registry: &registry,
	}
}

func (a *Agent) Run(userInput string) (string, error) {
	toolDescriptions := BuildToolDescription(a.registry) // 根据工具注册表生成工具说明
	systemPrompt := BuildSystemPrompt(toolDescriptions)  // 把工具说明塞进系统提示词里

	messages := []llm.Message{
		{Role: "system", Content: systemPrompt}, // 系统消息，告诉模型行为规则
		{Role: "user", Content: userInput},      // 用户消息，表示用户真实输入
	}

	const maxSteps = 5 // Agent最多和模型交互5轮

	for step := 0; step < maxSteps; step++ {
		response, err := a.client.Chat(messages)
		if err != nil {
			return "", err
		}
		// 把当前消息列表发给模型

		cleaned := stripThink(response)

		var call ToolCall
		if err := json.Unmarshal([]byte(cleaned), &call); err != nil {
			return "", fmt.Errorf("无效的模型响应: %w; raw=%s", err, response)
		}
		// 解析模型的返回，如果模型没有返回合法JSON，就会报错

		if call.Tool == "final_answer" {
			return call.Input, nil
		}
		// 如果模型选择的工具是"final_answer"，说明它已经就最终答案了，返回call.Input
		// 这也就是最终回答的内容

		result, err := a.registry.Execute(call.Tool, call.Input)
		if err != nil {
			result = "tool error: " + err.Error()
		}
		// 如果不是最终答案，就把模型制定的工具名和输入交给工具注册表去执行

		messages = append(messages, llm.Message{
			Role:    "assistant",
			Content: response,
		})
		// 如果工具执行失败，不直接终止，而是把错误作为工具结果返回给模型

		messages = append(messages, llm.Message{
			Role:    "tool",
			Content: result,
		})
		// 把模型刚才的JSON回复加入上下文
	}

	return "", fmt.Errorf("agent reached max steps")
}

var thinkBlock = regexp.MustCompile(`(?s)<think>.*?</think>`)

func stripThink(s string) string {
	s = thinkBlock.ReplaceAllString(s, "")
	return strings.TrimSpace(s)
}

// 有时候模型会多说几句(里面会含think<>)，但是协议需要返回纯JSON，
// 所以在agent层做一次标准化
