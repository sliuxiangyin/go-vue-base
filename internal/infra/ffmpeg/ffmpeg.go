package ffmpeg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// FFmpeg FFmpeg 工具封装
type FFmpeg struct{}

// NewFFmpeg 创建 FFmpeg 实例
func NewFFmpeg() *FFmpeg {
	return &FFmpeg{}
}

// ProbeResult ffprobe 返回的结果结构
type ProbeResult struct {
	Format ProbeFormat `json:"format"`
}

// ProbeFormat 格式信息
type ProbeFormat struct {
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bit_rate"`
	FormatName string `json:"format_name"`
}

// GetAudioDuration 获取音频文件时长（秒）
func (f *FFmpeg) GetAudioDuration(filePath string) (float64, error) {
	// 使用 ffprobe 获取音频信息
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "json",
		filePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to execute ffprobe: %w, output: %s", err, string(output))
	}

	// 解析 JSON 输出
	var result ProbeResult
	if err := json.Unmarshal(output, &result); err != nil {
		return 0, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	// 将 duration 字符串转换为 float64
	duration, err := strconv.ParseFloat(result.Format.Duration, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w", err)
	}

	return duration, nil
}

// GetAudioDurationSimple 获取音频时长（简单方式，解析文本输出）
func (f *FFmpeg) GetAudioDurationSimple(filePath string) (float64, error) {
	// 使用 ffprobe 获取音频信息（文本输出）
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to execute ffprobe: %w, output: %s", err, string(output))
	}

	// 解析输出的数字
	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration: %w, output: %s", err, durationStr)
	}

	return duration, nil
}

// GetAudioInfo 获取音频完整信息
func (f *FFmpeg) GetAudioInfo(filePath string) (*ProbeResult, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_format",
		"-of", "json",
		filePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute ffprobe: %w, output: %s", err, string(output))
	}

	var result ProbeResult
	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse ffprobe output: %w", err)
	}

	return &result, nil
}

// CheckFFmpegInstalled 检查 ffmpeg 和 ffprobe 是否已安装
func (f *FFmpeg) CheckFFmpegInstalled() error {
	// 检查 ffprobe
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return fmt.Errorf("ffprobe not found in PATH, please install ffmpeg")
	}

	// 检查 ffmpeg
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return fmt.Errorf("ffmpeg not found in PATH, please install ffmpeg")
	}

	return nil
}
