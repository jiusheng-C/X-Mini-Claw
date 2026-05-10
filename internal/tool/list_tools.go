package tool

import (
	"fmt"
	"strings"
)

type ListTools struct {
	registry *Registry
}

func NewListTools(registry *Registry) *ListTools {
	return &ListTools{
		registry: registry,
	}
}

func (t *ListTools) Name() string {
	return "list_tools"
}

func (t *ListTools) Description() string {
	return "打印一个列表, 包含全部现有工具。"
}

func (t *ListTools) Run(input string) (string, error) {
	var builder strings.Builder
	builder.WriteString("已注册工具:\n")

	for _, tool := range t.registry.ListTools() {
		builder.WriteString(fmt.Sprintf("  %s - %s\n", tool.Name(), tool.Description()))
	}

	return builder.String(), nil
}
