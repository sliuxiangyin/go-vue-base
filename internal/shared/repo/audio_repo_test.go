package repo

import (
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/grpc"
	"databaseAi/internal/infra/storage"
	"testing"
	"time"
)

func TestAudioRepo_Transcribe(t *testing.T) {
	// 创建配置
	cfg := config.NewConfig("dev")

	// 创建 gRPC 工厂
	grpcFactory := grpc.NewGrpc(cfg.GrpcDNS)

	// 创建文件存储
	fileStorage := storage.NewFileStorage(cfg.AppEnv)

	// 创建 AudioRepo
	repo := NewAudioRepo(grpcFactory, cfg, fileStorage)

	// 测试音频文件路径（需要先有一个测试音频文件）
	audioPath := "D:\\code\\databaseAi\\backend\\测试.mp3"

	// 调用转录
	result, err := repo.Transcribe(audioPath, "en", "small")
	if err != nil {
		t.Fatalf("Transcribe failed: %v", err)
	}

	// 验证结果
	if result.Text == "" {
		t.Error("Transcribed text is empty")
	}

	if len(result.Words) == 0 {
		t.Error("No word timestamps returned")
	}

	t.Logf("Transcribed text: %s", result.Text)
	t.Logf("Language: %s", result.Language)
	t.Logf("Word count: %d", len(result.Words))

	// 打印前 5 个单词的时间戳
	for i, word := range result.Words {
		if i >= 5 {
			break
		}
		t.Logf("Word[%d]: %s (%.2f-%.2f, confidence: %.2f)",
			i, word.Word, word.Start, word.End, word.Confidence)
	}
}

func TestAudioRepo_Transcribe_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	cfg := config.NewConfig("dev")
	grpcFactory := grpc.NewGrpc(cfg.GrpcDNS)
	fileStorage := storage.NewFileStorage(cfg.AppEnv)
	repo := NewAudioRepo(grpcFactory, cfg, fileStorage)

	// 先生成一个音频文件
	t.Log("Step 1: Generating audio...")
	audioPath, err := repo.Synthesizer(
		"qwen3-tts-flash",
		"Jennifer",
		"Hello world, this is a test.",
		"English",
	)
	if err != nil {
		t.Fatalf("Synthesizer failed: %v", err)
	}
	t.Logf("Audio saved to: %s", audioPath)

	// 等待文件写入完成
	time.Sleep(1 * time.Second)

	// 转录该音频
	t.Log("Step 2: Transcribing audio...")
	result, err := repo.Transcribe(audioPath, "en", "small")
	if err != nil {
		t.Fatalf("Transcribe failed: %v", err)
	}

	t.Logf("Transcribed text: %s", result.Text)
	t.Logf("Words: %d", len(result.Words))
}
