package tool

import "fmt"

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry { // 注册新的工具注册表
	return &Registry{
		make(map[string]Tool),
	}
}

func (r *Registry) Register(t Tool) {
	r.tools[t.Name()] = t // 将工具存储到注册表中，键为工具的名称
}

func (r *Registry) Execute(name, input string) (string, error) {
	t, ok := r.tools[name] // 通过名称调用已注册的工具
	if !ok {
		return "", fmt.Errorf("未找到工具: %s", name) // 输出工具调用失败的原因
	}
	return t.Run(input) // 返回执行结果
}

// 单个工具自己执行叫Run
// 注册器根据名称调度工具叫Execute

func (r *Registry) ListTools() []Tool { // 列出注册表中所有已注册的工具
	result := make([]Tool, 0, len(r.tools))
	for _, t := range r.tools {
		result = append(result, t)
	}
	return result
}
