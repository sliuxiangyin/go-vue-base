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

// ChunkLevel 语义切分等级
type ChunkLevel int

const (
	ChunkLevel1 ChunkLevel = 1 // 粗颗粒度：长语块
	ChunkLevel2 ChunkLevel = 2 // 中颗粒度：3-6词适中语块
	ChunkLevel3 ChunkLevel = 3 // 细颗粒度：最小朗读单位
)

// AnalyzeSemanticChunks 使用通义千问模型分析英文句子的语义意群
// level: 切分等级 (1=粗颗粒度, 2=中颗粒度, 3=细颗粒度)
func (r *SemanticRepo) AnalyzeSemanticChunks(ctx context.Context, sentence string, level ChunkLevel) ([]string, error) {
	client, err := r.openai.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get openai client: %w", err)
	}

	// 根据 level 选择对应的 prompt
	prompt := r.buildPrompt(sentence, level)

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

// buildPrompt 根据 level 构建对应的 prompt
func (r *SemanticRepo) buildPrompt(sentence string, level ChunkLevel) string {
	switch level {
	case ChunkLevel1:
		return r.buildLevel1Prompt(sentence)
	case ChunkLevel2:
		return r.buildLevel2Prompt(sentence)
	case ChunkLevel3:
		return r.buildLevel3Prompt(sentence)
	default:
		return r.buildLevel2Prompt(sentence) // 默认使用 Level 2
	}
}

// buildLevel1Prompt LEVEL 1（粗颗粒度：长语块）
func (r *SemanticRepo) buildLevel1Prompt(sentence string) string {
	return fmt.Sprintf(`你是一个英文语义意群切分器（English Semantic Chunker）。

任务：根据 LEVEL 1（粗颗粒度：长语块）将英文句子切分为语义自然、朗读节奏良好的意群短语。

切分时禁止按标点切分，必须按语义结构与发音节奏切分。

LEVEL 1 规则如下（请严格遵守）：

1. 每个意群必须表达一个完整、较大的语义概念。
2. 每段允许较长，只要朗读自然即可。
3. 禁止产生只有 2–3 个词的短碎片。
4. 语义紧密相关的词组必须合并为较大的意群。
5. 不按标点切分，只按语义逻辑切分。

示例（错误 ❌）
["Youth is not", "a time", "of life"]

示例（正确 ✅）
["Youth is not a time of life"]

通用输出要求：
1. 只能输出 JSON 数组。
2. 不要返回任何解释性文字或说明。
3. 切分必须遵守 LEVEL 1 的全部硬性规则。

英文句子：%s

请只返回 JSON 数组，不要其他解释。`, sentence)
}

// buildLevel2Prompt LEVEL 2（中颗粒度：3–6 词适中语块）
func (r *SemanticRepo) buildLevel2Prompt(sentence string) string {
	return fmt.Sprintf(`你是一个英文语义意群切分器（English Semantic Chunker）。

任务：根据 LEVEL 2（中颗粒度：3–6 词适中语块）将英文句子切分为语义自然、朗读节奏良好的意群短语。

切分时禁止按标点切分，必须按语义结构与发音节奏切分。

LEVEL 2 规则如下（请严格遵守）：

1. 每个意群应保持 3–6 个英文单词。
2. 要方便一般学习者朗读，不要太长，不要太短。
3. 较长的语义块需要拆成 2–3 段。
4. 禁止出现超过 8 个词的长片段。
5. 不允许出现只有 1–2 词的极短碎片。
6. 按自然语音停顿点切分，例如：主谓结构、介词短语等。

示例（错误 ❌ 太长）
["it is the freshness of the deep springs of life"]

示例（错误 ❌ 太短）
["the deep", "springs of", "life"]

示例（正确 ✅）
["it is the freshness", "of the deep springs of life"]

通用输出要求：
1. 只能输出 JSON 数组。
2. 不要返回任何解释性文字或说明。
3. 切分必须遵守 LEVEL 2 的全部硬性规则。

英文句子：%s

请只返回 JSON 数组，不要其他解释。`, sentence)
}

// buildLevel3Prompt LEVEL 3（细颗粒度：最小朗读单位）
func (r *SemanticRepo) buildLevel3Prompt(sentence string) string {
	return fmt.Sprintf(`你是一个英文语义意群切分器（English Semantic Chunker）。

任务：根据 LEVEL 3（细颗粒度：最小朗读单位）将英文句子切分为语义自然、朗读节奏良好的意群短语。

切分时禁止按标点切分，必须按语义结构与发音节奏切分。

LEVEL 3 规则如下（请严格遵守）：

1. 每个意群尽可能短且自然。
2. 每段最多 3 个英文单词（特殊短语最多允许到 4 个）。
3. 能拆则必须拆，禁止将多个语义片段合并。
4. 严禁超过 4 个词的长意群。
5. 保持自然语音节奏：主语短语、动词短语、介词短语、形容词短语可独立成组。
6. 不按标点切分，只按微语义关系划分。

示例（错误 ❌ 太长）
["Youth is not a time of life"]

示例（正确 ✅）
["Youth is not", "a time", "of life"]

通用输出要求：
1. 只能输出 JSON 数组。
2. 不要返回任何解释性文字或说明。
3. 切分必须遵守 LEVEL 3 的全部硬性规则。

英文句子：%s

请只返回 JSON 数组，不要其他解释。`, sentence)
}

// AlignmentRecord 中英文对齐记录
type AlignmentRecord struct {
	Src        string  `json:"src"`         // 英文词或短语
	SrcIndices []int   `json:"src_indices"` // 英文对应原句的索引列表
	Tgt        string  `json:"tgt"`         // 对应中文词或短语
	TgtIndices []int   `json:"tgt_indices"` // 中文对应索引列表
	Score      float64 `json:"score"`       // 对齐置信度（0-1）
}

// AlignSemanticChunks 中英文语义意群对齐
func (r *SemanticRepo) AlignSemanticChunks(ctx context.Context, englishSentence, chineseSentence string) ([]AlignmentRecord, error) {
	client, err := r.openai.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get openai client: %w", err)
	}

	prompt := r.buildAlignmentPrompt(englishSentence, chineseSentence)

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
	var alignments []AlignmentRecord
	// 尝试直接解析
	err = json.Unmarshal([]byte(response), &alignments)
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

		err = json.Unmarshal([]byte(response), &alignments)
		if err != nil {
			return nil, fmt.Errorf("failed to parse response: %w, response: %s", err, response)
		}
	}

	return alignments, nil
}

// buildAlignmentPrompt 构建中英文对齐的 prompt
func (r *SemanticRepo) buildAlignmentPrompt(englishSentence, chineseSentence string) string {
	return fmt.Sprintf(`你是一个中英文语义对齐专家。

任务：将以下英文句子与中文句子进行语义意群对齐，生成 JSON 输出。

要求：

1. 英文短语（phrase）可与中文短语对齐（例如：dish out -> 分发）。
2. 中文先逐词分割，每个词作为索引单位。
3. 输出 JSON 数组，每条记录包含：
   - "src": 英文词或短语
   - "src_indices": 英文对应原句的索引列表（按空格分词后的索引，从 0 开始）
   - "tgt": 对应中文词或短语
   - "tgt_indices": 中文对应索引列表（按字符分词后的索引，从 0 开始）
   - "score": 对齐置信度（0-1，由你根据语义相关性判断）
4. JSON 可直接用于前端点击高亮。
5. 只返回 JSON 数组，不要任何解释性文字。

示例：
英文："When we dish out food!"
中文："当我们分发食物时！"

输出格式参考：
[
  {
    "src": "When",
    "src_indices": [0],
    "tgt": "当",
    "tgt_indices": [0],
    "score": 0.95
  },
  {
    "src": "we",
    "src_indices": [1],
    "tgt": "我们",
    "tgt_indices": [1, 2],
    "score": 0.98
  },
  {
    "src": "dish out",
    "src_indices": [2, 3],
    "tgt": "分发",
    "tgt_indices": [3, 4],
    "score": 0.92
  },
  {
    "src": "food",
    "src_indices": [4],
    "tgt": "食物",
    "tgt_indices": [5, 6],
    "score": 0.99
  }
]

现在请对齐以下句子：

英文句子：%s

中文句子：%s

请只返回 JSON 数组，不要其他解释。`, englishSentence, chineseSentence)
}
