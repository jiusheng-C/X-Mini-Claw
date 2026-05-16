package llm

type FakeClient struct{}

// FakeClient 结构体的Chat方法，模拟客户端聊天功能
// 该方法接收一个Message类型的切片作为输入，返回一个字符串和一个错误
func (c *FakeClient) Chat(message []Message) (string, error) {
	// 返回一个固定的JSON字符串作为模拟响应，表示这是一个假客户端的回复
	// JSON格式包含工具名称"final_answer"和输入内容"Fake client response."
	return `{"tool":"final_answer","input":"Fake client response."}`, nil
}
