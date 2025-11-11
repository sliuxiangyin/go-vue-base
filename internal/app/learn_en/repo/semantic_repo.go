package repo

import (
	"context"
	"databaseAi/internal/infra/ai"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
)

type SemanticRepo struct {
	openai *ai.Openai
}

func NewSemanticRepo(openai *ai.Openai) *SemanticRepo {
	return &SemanticRepo{
		openai: openai,
	}
}

// AnalyzeSemanticChunks 使用通义千问模型分析英文句子的语义意群
func (r *SemanticRepo) AnalyzeSemanticChunks(ctx context.Context, sentence string) ([]string, error) {
	client, err := r.openai.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get openai client: %w", err)
	}

	// 构造 prompt
	prompt := fmt.Sprintf(`你是一个英文语义分析器。

请将以下英文句子按"意群"分割，每个意群是一组逻辑上紧密关联的词或短语。

不要按标点或固定数量切分，而要按语义关系。

输出 JSON 数组格式，每个元素是一个意群字符串。

英文句子：%s

请只返回 JSON 数组，不要其他解释。`, sentence)

	// 调用通义千问 API
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model: "qwen-plus", // 使用通义千问模型
	})

	if err != nil {
		return nil, fmt.Errorf("failed to call openai api: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("no response from openai")
	}

	// 解析响应
	response := chatCompletion.Choices[0].Message.Content

	// 提取 JSON 数组
	var chunks []string
	// 尝试直接解析
	err = json.Unmarshal([]byte(response), &chunks)
	if err != nil {
		// 如果直接解析失败，尝试提取 JSON 部分
		response = strings.TrimSpace(response)
		if strings.HasPrefix(response, "```json") {
			response = strings.TrimPrefix(response, "```json")
			response = strings.TrimSuffix(response, "```")
			response = strings.TrimSpace(response)
		} else if strings.HasPrefix(response, "```") {
			response = strings.TrimPrefix(response, "```")
			response = strings.TrimSuffix(response, "```")
			response = strings.TrimSpace(response)
		}

		err = json.Unmarshal([]byte(response), &chunks)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response: %w, response: %s", err, response)
		}
	}

	return chunks, nil
}
