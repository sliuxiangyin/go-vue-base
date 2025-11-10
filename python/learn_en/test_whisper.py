import whisper_timestamped as whisper
import json

# 加载模型
model = whisper.load_model("small", download_root="./models")  # tiny/base/small/medium/large
# 读取音频
audio_path = "./downloaded_audio.wav"   # 用户录音或标准音频
audio = whisper.load_audio(audio_path)
audio = whisper.pad_or_trim(audio)

# 转录并获取逐词时间戳
result = whisper.transcribe(model,
                            audio,
                            language="en",
                            trust_whisper_timestamps=False,  # 允许自动调整边界
                            refine_whisper_precision=0.1,  # 提高时间戳精度
                            min_word_duration=0.05,  # 保持词时长平滑
                            compute_word_confidence=True,
                            remove_empty_words=True,
                            )
# 构造 JSON
output = []

for segment in result["segments"]:
    # 假设你有一个函数可以把英文翻译成中文，这里示例用空字符串
    zh_text = ""  # TODO: 可调用翻译接口生成中文
    en_text = segment["text"]

    words_list = []
    for w in segment["words"]:
        words_list.append({
            "word": w["text"],
            "start": (w["start"]),
            "end": (w["end"])
        })

    output.append({
        "segment": {
            "zh": zh_text,
            "en": en_text
        },
        "words": words_list
    })
print(output)
