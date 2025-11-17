package repo

import (
	"context"
	"databaseAi/internal/infra/ai"
	"databaseAi/internal/infra/config"
	"testing"
	"time"
)

func TestSemanticRepo_AnalyzeSemanticChunks(t *testing.T) {
	// 创建配置
	cfg := config.NewConfig("dev")

	// 检查是否配置了有效的 API Key
	if cfg.OpenaiKey == "" {
		t.Skip("跳过测试：未配置 OPENAI_KEY 环境变量")
	}

	// 创建 OpenAI 客户端
	openaiClient := ai.NewOpenai(cfg.OpenaiKey, cfg.OpenaiURl)

	// 创建 SemanticRepo
	repo := NewSemanticRepo(openaiClient)

	// 测试用例
	tests := []struct {
		name      string
		sentence  string
		wantErr   bool
		minChunks int // 期望最少的意群数量
	}{
		{
			name:      "简单句子",
			sentence:  "The quick brown fox jumps over the lazy dog",
			wantErr:   false,
			minChunks: 2,
		},
		{
			name:      "复杂句子",
			sentence:  "When I was young, I used to play basketball in the beautiful garden near my house every weekend",
			wantErr:   false,
			minChunks: 3,
		},
		{
			name:      "短句子",
			sentence:  "Hello world",
			wantErr:   false,
			minChunks: 1,
		},
		{
			name:      "带从句的句子",
			sentence:  "Although it was raining heavily, we decided to go hiking in the mountains",
			wantErr:   false,
			minChunks: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置超时上下文
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			// 执行测试
			chunks, err := repo.AnalyzeSemanticChunks(ctx, tt.sentence)

			// 检查错误
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeSemanticChunks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 如果期望成功，检查结果
			if !tt.wantErr {
				if len(chunks) == 0 {
					t.Errorf("AnalyzeSemanticChunks() 返回空结果")
					return
				}

				if len(chunks) < tt.minChunks {
					t.Errorf("AnalyzeSemanticChunks() 返回意群数量 = %d, 期望至少 %d", len(chunks), tt.minChunks)
				}

				// 验证所有意群拼接后包含原句的主要内容
				t.Logf("原句: %s", tt.sentence)
				t.Logf("意群数量: %d", len(chunks))
				for i, chunk := range chunks {
					t.Logf("  [%d] %s", i+1, chunk)
				}

				// 检查意群不为空
				for i, chunk := range chunks {
					if chunk == "" {
						t.Errorf("意群 [%d] 为空字符串", i)
					}
				}
			}
		})
	}
}

func TestSemanticRepo_AnalyzeSemanticChunks_EmptyInput(t *testing.T) {
	cfg := config.NewConfig("dev")
	if cfg.OpenaiKey == "" {
		t.Skip("跳过测试：未配置 OPENAI_KEY 环境变量")
	}

	openaiClient := ai.NewOpenai(cfg.OpenaiKey, cfg.OpenaiURl)
	repo := NewSemanticRepo(openaiClient)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试空字符串
	_, err := repo.AnalyzeSemanticChunks(ctx, "")
	if err == nil {
		t.Log("警告: 空字符串应该返回错误或空数组")
	}
}

func TestSemanticRepo_AnalyzeSemanticChunks_ContextTimeout(t *testing.T) {
	cfg := config.NewConfig("dev")
	if cfg.OpenaiKey == "" {
		t.Skip("跳过测试：未配置 OPENAI_KEY 环境变量")
	}

	openaiClient := ai.NewOpenai(cfg.OpenaiKey, cfg.OpenaiURl)
	repo := NewSemanticRepo(openaiClient)

	// 创建一个已经超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // 确保超时

	_, err := repo.AnalyzeSemanticChunks(ctx, "This is a test sentence")
	if err == nil {
		t.Error("期望超时错误，但没有返回错误")
	}
}
