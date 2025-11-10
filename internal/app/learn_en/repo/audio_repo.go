package repo

import (
	"context"
	protos "databaseAi/internal/app/learn_en/proto"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/grpc"
	"fmt"
	"time"
)

type AudioRepo struct {
	grpcFactory *grpc.GrpcFactory
	apiKey      string
}

func NewAudioRepo(grpcFactory *grpc.GrpcFactory, config *config.Config) *AudioRepo {
	return &AudioRepo{
		grpcFactory: grpcFactory,
		apiKey:      config.OpenaiKey,
	}
}

// Synthesizer 调用 gRPC 服务进行语音合成
func (r *AudioRepo) Synthesizer(model, voice, text string) ([]byte, error) {
	// 获取 gRPC 连接
	conn, err := r.grpcFactory.Get()
	if err != nil {
		return nil, fmt.Errorf("failed to get grpc connection: %w", err)
	}

	// 创建 Audio 客户端
	client := protos.NewAudioClient(conn)

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 调用 Synthesizer 方法
	resp, err := client.Synthesizer(ctx, &protos.SynthesizerRequest{
		ApiKey: r.apiKey,
		Model:  model,
		Voice:  voice,
		Text:   text,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc call failed: %w", err)
	}

	// 检查响应状态
	if resp.Code != 0 {
		return nil, fmt.Errorf("synthesizer failed: %s (code: %d)", resp.Message, resp.Code)
	}

	return resp.Data, nil
}
