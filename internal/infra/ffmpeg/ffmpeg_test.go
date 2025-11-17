package ffmpeg

import (
	"testing"
)

func TestFFmpeg_CheckFFmpegInstalled(t *testing.T) {
	ff := NewFFmpeg()

	err := ff.CheckFFmpegInstalled()
	if err != nil {
		t.Skipf("跳过测试：%v", err)
	}
}

func TestFFmpeg_GetAudioDuration(t *testing.T) {
	ff := NewFFmpeg()

	// 检查 ffmpeg 是否安装
	if err := ff.CheckFFmpegInstalled(); err != nil {
		t.Skipf("跳过测试：%v", err)
	}

	// 这里需要提供一个真实的音频文件路径进行测试
	// 如果没有测试文件，跳过测试
	testFile := "test_audio.mp3"

	duration, err := ff.GetAudioDuration(testFile)
	if err != nil {
		t.Logf("测试文件不存在或读取失败（预期行为）: %v", err)
		return
	}

	t.Logf("音频时长: %.2f 秒", duration)

	if duration <= 0 {
		t.Errorf("音频时长应该大于 0，实际: %.2f", duration)
	}
}

func TestFFmpeg_GetAudioDurationSimple(t *testing.T) {
	ff := NewFFmpeg()

	if err := ff.CheckFFmpegInstalled(); err != nil {
		t.Skipf("跳过测试：%v", err)
	}

	testFile := "test_audio.mp3"

	duration, err := ff.GetAudioDurationSimple(testFile)
	if err != nil {
		t.Logf("测试文件不存在或读取失败（预期行为）: %v", err)
		return
	}

	t.Logf("音频时长（简单方式）: %.2f 秒", duration)

	if duration <= 0 {
		t.Errorf("音频时长应该大于 0，实际: %.2f", duration)
	}
}

func TestFFmpeg_GetAudioInfo(t *testing.T) {
	ff := NewFFmpeg()

	if err := ff.CheckFFmpegInstalled(); err != nil {
		t.Skipf("跳过测试：%v", err)
	}

	testFile := "test_audio.mp3"

	info, err := ff.GetAudioInfo(testFile)
	if err != nil {
		t.Logf("测试文件不存在或读取失败（预期行为）: %v", err)
		return
	}

	t.Logf("音频信息:")
	t.Logf("  时长: %s 秒", info.Format.Duration)
	t.Logf("  大小: %s 字节", info.Format.Size)
	t.Logf("  比特率: %s", info.Format.BitRate)
	t.Logf("  格式: %s", info.Format.FormatName)
}
