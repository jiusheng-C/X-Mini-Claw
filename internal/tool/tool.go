package tool

type Tool interface {
	Name() string                     // 工具名称
	Description() string              // 工具描述
	Run(input string) (string, error) // 工具执行逻辑
}
