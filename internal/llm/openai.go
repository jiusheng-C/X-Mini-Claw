package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAIClient struct {
	APIKey     string
	BaseURL    string
	Model      string
	APIStyle   string
	HTTPClient *http.Client
}

func NewOpenAIClient(apiKey, baseURL, model, apiStyle string) *OpenAIClient {
	apiStyle = strings.TrimSpace(apiStyle)
	if apiStyle == "" {
		apiStyle = "responses"
	}

	return &OpenAIClient{
		APIKey:   apiKey,
		BaseURL:  strings.TrimRight(baseURL, "/"),
		Model:    model,
		APIStyle: apiStyle,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *OpenAIClient) Chat(messages []Message) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("api key为空")
	}
	if c.BaseURL == "" {
		return "", fmt.Errorf("base url为空")
	}
	if c.Model == "" {
		return "", fmt.Errorf("model为空")
	}

	switch c.APIStyle {
	case "responses":
		return c.chatWithResponses(messages)
	case "chat_completions":
		return c.chatWithCompletions(messages)
	default:
		return "", fmt.Errorf("unsupported api_style: %s", c.APIStyle)
	}
}

func (c *OpenAIClient) chatWithResponses(messages []Message) (string, error) {
	input := make([]responsesInputItem, 0, len(messages))
	for _, m := range messages {
		input = append(input, responsesInputItem{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	reqBody := responsesRequest{
		Model: c.Model,
		Input: input,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,        // POST请求方法
		c.BaseURL+"/responses", // 请求地址
		bytes.NewReader(body),  // 把刚才生成的JSON当做请求体
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req) // 发送HTTP请求
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // 响应读完后关闭资源

	respBody, err := io.ReadAll(resp.Body) // 读取服务器返回的完整内容，转成字符串返回
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai api error : status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var result responsesResponse
	if err := json.Unmarshal(respBody, &result); err != nil { // 将响应JSON解析到结构体
		return "", err
	}

	if len(result.Output) == 0 {
		return "", fmt.Errorf("responses api returned empty output: %s", string(respBody))
	}

	for _, out := range result.Output {
		for _, content := range out.Content {
			if (content.Type == "output_text" || content.Type == "text") && content.Text != "" {
				return content.Text, nil
			}
		}
	}

	return "", fmt.Errorf("responses api returned no text content: %s", string(respBody))
}

func (c *OpenAIClient) chatWithCompletions(messages []Message) (string, error) {
	reqBody := chatCompletionRequest{
		Model:    c.Model,
		Messages: messages,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.BaseURL+"/chat/completions",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai api error : status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var result chatCompletionResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("chat completions api returned no choices: %s", string(respBody))
	}

	return result.Choices[0].Message.Content, nil
}

type responsesRequest struct {
	Model string               `json:"model"`
	Input []responsesInputItem `json:"input"`
}

type responsesInputItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responsesResponse struct {
	Output []struct {
		Type    string `json:"type"`
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

type chatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}
