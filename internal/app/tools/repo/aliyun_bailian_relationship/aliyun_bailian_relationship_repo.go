package aliyun_bailian_relationship

import (
	"context"
	"databaseAi/internal/infra/ai"
	"fmt"
	"github.com/openai/openai-go"
)

// AliYunBaiLianRelationshipRepo todo 应该将各个relationship 存储起来
// AliYunBaiLianRelationshipRepo 阿里云百炼 数据库关系分析
type AliYunBaiLianRelationshipRepo struct {
	openai       *ai.Openai
	baseMessages []openai.ChatCompletionMessageParamUnion // 永久系统/背景消息
	messages     []openai.ChatCompletionMessageParamUnion // 临时会话消息
}

// NewAliYunBaiLianRelationshipRepo 初始化
func NewAliYunBaiLianRelationshipRepo(infraOpenai *ai.Openai) *AliYunBaiLianRelationshipRepo {
	systemPrompt := SystemPrompt() // 你的 system prompt
	return &AliYunBaiLianRelationshipRepo{
		openai:       infraOpenai,
		baseMessages: []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(systemPrompt)},
		messages:     []openai.ChatCompletionMessageParamUnion{},
	}
}

// AddBaseMessage 添加永久存在的消息
// 会在每次 Chat() 调用中一起发送
func (a *AliYunBaiLianRelationshipRepo) AddBaseMessage(messages []string) {
	for i, msg := range messages {
		a.baseMessages = append(a.baseMessages, openai.UserMessage(UserChunkDDlPrompt(i, msg)))
	}
}

// AddMessage 添加临时消息，仅当前对话有效
func (a *AliYunBaiLianRelationshipRepo) AddMessage(messages []string) {
	for _, msg := range messages {
		a.messages = append(a.messages, openai.UserMessage(msg))
	}
}

// ClearAllMessages ClearMessages 清除全部消息
func (a *AliYunBaiLianRelationshipRepo) ClearAllMessages() {
	a.messages = make([]openai.ChatCompletionMessageParamUnion, 0)
	a.baseMessages = make([]openai.ChatCompletionMessageParamUnion, 0)
}

// Chat 执行对话请求
func (a *AliYunBaiLianRelationshipRepo) Chat() (<-chan string, error) {
	var chatSteamChan = make(chan string, 1)

	client, err := a.openai.Get()
	if err != nil {
		return chatSteamChan, err
	}
	// 拼接消息序列
	allMessages := append(a.baseMessages, a.messages...)
	stream := client.Chat.Completions.NewStreaming(
		context.TODO(),
		openai.ChatCompletionNewParams{
			Model:    "qwen3-coder-plus",
			Messages: allMessages,
		},
	)
	go func() {
		defer close(chatSteamChan)
		defer stream.Close()
		for stream.Next() {
			current := stream.Current()
			chatSteamChan <- current.Choices[0].Delta.Content
			fmt.Println(current.Choices[0].Delta.Content)
		}
	}()
	return chatSteamChan, nil
}
