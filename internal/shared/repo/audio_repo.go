package repo

import (
	"context"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/grpc"
	"databaseAi/internal/infra/storage"
	protos "databaseAi/internal/shared/proto"
	"fmt"
	"time"
)

type AudioRepo struct {
	grpcFactory *grpc.GrpcFactory
	apiKey      string
	fileStorage *storage.FileStorage
}

func NewAudioRepo(grpcFactory *grpc.GrpcFactory, config *config.Config, fileStorage *storage.FileStorage) *AudioRepo {
	return &AudioRepo{
		grpcFactory: grpcFactory,
		apiKey:      config.OpenaiKey,
		fileStorage: fileStorage,
	}
}

// Synthesizer 调用 gRPC 服务进行语音合成，返回音频 URL
func (r *AudioRepo) Synthesizer(model, voice, text, languageType string) (string, error) {
	// 获取 gRPC 连接
	conn, err := r.grpcFactory.Get()
	if err != nil {
		return "", fmt.Errorf("failed to get grpc connection: %w", err)
	}

	// 创建 Audio 客户端
	client := protos.NewAudioClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用 Synthesizer 方法
	resp, err := client.Synthesizer(ctx, &protos.SynthesizerRequest{
		ApiKey:       r.apiKey,
		Model:        model,
		Voice:        voice,
		Text:         text,
		LanguageType: languageType,
	})
	if err != nil {
		return "", fmt.Errorf("grpc call failed: %w", err)
	}

	// 检查响应状态
	if resp.Code != 0 {
		return "", fmt.Errorf("synthesizer failed: %s (code: %d)", resp.Message, resp.Code)
	}

	// 从 URL 下载文件并保存到本地
	relativePath, err := r.fileStorage.DownloadFromURL(resp.Url)
	if err != nil {
		return "", fmt.Errorf("failed to save audio file: %w", err)
	}

	return relativePath, nil
}

// Transcribe 调用 gRPC 服务进行音频转录，返回逐词时间戳
func (r *AudioRepo) Transcribe(audioPath, language, modelSize string) (*protos.TranscribeReply, error) {

	audioPath = r.fileStorage.GetFullPath(audioPath)
	// 获取 gRPC 连接
	conn, err := r.grpcFactory.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: %w", err)
	}

	// 创建 Audio 客户端
	client := protos.NewAudioClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 调用 Transcribe 方法
	resp, err := client.Transcribe(ctx, &protos.TranscribeRequest{
		AudioPath: audioPath,
		Language:  language,
		ModelSize: modelSize,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc call failed: %w", err)
	}

	// 检查响应状态
	if resp.Code != 0 {
		return nil, fmt.Errorf("transcribe failed: %s (code: %d)", resp.Message, resp.Code)
	}

	return resp, nil
}
