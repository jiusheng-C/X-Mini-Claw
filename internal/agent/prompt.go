package agent

import (
	"fmt"
	"strings"

	"github.com/jiusheng-C/X-Mini-Claw/internal/tool"
)

// BuildToolDescription 构建工具描述字符串
// 输入参数:
//   - registry: 工具注册表指针，包含所有可用的工具信息
//
// 返回值:
//   - string: 格式化后的工具描述字符串，每个工具占一行，格式为"工具名: 工具描述"
func BuildToolDescription(registry *tool.Registry) string {
	var builder strings.Builder // 使用strings.Builder高效构建字符串

	// 遍历注册表中的所有工具
	for _, t := range registry.ListTools() {
		// 为每个工具添加格式化的描述信息
		builder.WriteString(fmt.Sprintf("- %s: %s\n", t.Name(), t.Description()))
	}

	return builder.String() // 返回构建完成的工具描述字符串
}

func BuildSystemPrompt(toolDescriptions string) string {
	return fmt.Sprintf(`
你是一个本地代码智能体。

你可以使用以下工具：

%s

你必须只返回 JSON，不要输出 JSON 之外的任何内容。

JSON 格式如下：
{
  "tool": "tool_name",
  "input": "tool_input"
}

当你已经获得足够信息时，使用：
{
  "tool": "final_answer",
  "input": "你的最终回答"
}
`, toolDescriptions)
}
