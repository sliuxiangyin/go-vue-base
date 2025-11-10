# coding=utf-8
import dashscope
from dashscope.audio.tts_v2 import *
from proto import audio_pb2_grpc, audio_pb2

class Audio(audio_pb2_grpc.AudioServicer):
    def Synthesizer(self, request, context):
        print(request)
        # 若没有将API Key配置到环境变量中，需将your-api-key替换为自己的API Key
        dashscope.api_key = request.api_key
        # 模型
        model = request.model
        # 音色
        voice =request.voice
        # 实例化SpeechSynthesizer，并在构造方法中传入模型（model）、音色（voice）等请求参数
        synthesizer = SpeechSynthesizer(model=model, voice=voice,language_hints=["en" ])
        # 发送待合成文本，获取二进制音频
        audio = synthesizer.call(request.text)
        # 首次发送文本时需建立 WebSocket 连接，因此首包延迟会包含连接建立的耗时
        print('[Metric] requestId为：{}，首包延迟为：{}毫秒'.format(
            synthesizer.get_last_request_id(),
            synthesizer.get_first_package_delay()))
        # 将音频保存至本地
        # with open('../output.mp3', 'wb') as f:
        #     f.write(audio)
        return audio_pb2.SynthesizerReply(data=audio,code=0)


