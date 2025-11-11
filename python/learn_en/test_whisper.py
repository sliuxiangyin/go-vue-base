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

                            )
# 构造 JSON

print(result)
