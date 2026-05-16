package agent

type ToolCall struct {
	Tool  string `json:"tool"`
	Input string `json:"input"`
}
