package llm

// Message 定义了消息结构体，用于存储聊天消息的角色和内容
type Message struct {
	Role    string `json:"role"`    // Role 表示消息发送者的角色，如"user"或"assistant"
	Content string `json:"content"` // Conten 表示消息的具体内容
}

// Client 定义了聊天客户端接口，所有聊天客户端都需要实现这个接口
type Client interface {
	// Chat 接收消息列表作为输入，返回聊天回复内容和可能的错误
	Chat(message []Message) (string, error)
}
