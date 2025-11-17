package ffmpeg

import "github.com/google/wire"

// ProviderSet Wire 提供者集合
var ProviderSet = wire.NewSet(
	NewFFmpeg,
)
