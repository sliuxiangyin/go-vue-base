# coding=utf-8
"""测试 dashscope 连接"""
import os
import requests
import dashscope
# 设置 API Key
dashscope.api_key = "sk-4008472696c146f18096f2dab05dea85"
#  DashScope SDK 版本不低于 1.24.6


# 以下为北京地域url，若使用新加坡地域的模型，需将url替换为：https://dashscope-intl.aliyuncs.com/api/v1
dashscope.base_http_api_url = 'https://dashscope.aliyuncs.com/api/v1'

text = "Life was like a box of chocolates, you never know what you're gonna get."
# SpeechSynthesizer接口使用方法：dashscope.audio.qwen_tts.SpeechSynthesizer.call(...)
response = dashscope.MultiModalConversation.call(
    model="qwen3-tts-flash",
    # 新加坡和北京地域的API Key不同。获取API Key：https://help.aliyun.com/zh/model-studio/get-api-key
    # 若没有配置环境变量，请用百炼API Key将下行替换为：api_key = "sk-xxx"
    api_key="sk-4008472696c146f18096f2dab05dea85",
    text=text,
    voice="Jennifer",
    language_type="English", # 建议与文本语种一致，以获得正确的发音和自然的语调。
    stream=False
)
audio_url = response.output.audio.url
save_path = "downloaded_audio.wav"  # 自定义保存路径

try:
    response = requests.get(audio_url)
    response.raise_for_status()  # 检查请求是否成功
    with open(save_path, 'wb') as f:
        f.write(response.content)
    print(f"音频文件已保存至：{save_path}")
except Exception as e:
    print(f"下载失败：{str(e)}")
