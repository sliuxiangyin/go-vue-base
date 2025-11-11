# coding=utf-8
import dashscope
import whisper_timestamped as whisper
import os
from proto import audio_pb2_grpc, audio_pb2

class Audio(audio_pb2_grpc.AudioServicer):
    def __init__(self):
        # Whisper 模型缓存
        self.whisper_models = {}
        self.models_dir = "./models"
    
    def Synthesizer(self, request, context):
        # 设置 API Key
        dashscope.api_key = request.api_key
        
        # 设置基础 URL（北京地域）
        dashscope.base_http_api_url = 'https://dashscope.aliyuncs.com/api/v1'
        
        try:
            # 使用 MultiModalConversation 调用 TTS
            response = dashscope.MultiModalConversation.call(
                model=request.model,
                api_key=request.api_key,
                text=request.text,
                voice=request.voice,
                language_type=request.language_type if request.language_type else "English",
                stream=False
            )
            
            # 获取音频 URL
            audio_url = response.output.audio.url
            
            print(f'[Success] 生成音频 URL: {audio_url}')
            
            return audio_pb2.SynthesizerReply(url=audio_url, code=0, message="success")
            
        except Exception as e:
            error_msg = f'生成音频失败: {str(e)}'
            print(f'[Error] {error_msg}')
            return audio_pb2.SynthesizerReply(url="", code=1, message=error_msg)
    
    def Transcribe(self, request, context):
        """Whisper 音频转录，返回逐词时间戳"""
        try:
            audio_path = request.audio_path
            language = request.language if request.language else "en"
            model_size = request.model_size if request.model_size else "small"
            
            print(f'[Transcribe] audio_path={audio_path}, language={language}, model={model_size}')
            
            # 检查音频文件是否存在
            if not os.path.exists(audio_path):
                return audio_pb2.TranscribeReply(
                    code=1,
                    message=f"Audio file not found: {audio_path}"
                )
            
            # 加载或获取缓存的模型
            if model_size not in self.whisper_models:
                print(f'[Transcribe] Loading Whisper model: {model_size}')
                self.whisper_models[model_size] = whisper.load_model(
                    model_size,
                    download_root=self.models_dir
                )
            
            model = self.whisper_models[model_size]
            
            # 加载音频
            audio = whisper.load_audio(audio_path)
            
            # 转录并获取逐词时间戳
            print(f'[Transcribe] Transcribing audio...')
            result = whisper.transcribe(
                model,
                audio,
                language=language
            )
            
            # 提取完整文本
            full_text = result.get('text', '')
            detected_language = result.get('language', language)
            
            # 提取逐词时间戳
            word_timestamps = []
            if 'segments' in result:
                for segment in result['segments']:
                    if 'words' in segment:
                        for word_info in segment['words']:
                            word_timestamps.append(audio_pb2.WordTimestamp(
                                word=word_info.get('text', ''),
                                start=word_info.get('start', 0.0),
                                end=word_info.get('end', 0.0),
                                confidence=word_info.get('confidence', 0.0)
                            ))
            
            print(f'[Transcribe] Success: {len(word_timestamps)} words, text="{full_text[:50]}..."')
            
            return audio_pb2.TranscribeReply(
                text=full_text,
                words=word_timestamps,
                language=detected_language,
                code=0,
                message="success"
            )
            
        except Exception as e:
            error_msg = f'转录失败: {str(e)}'
            print(f'[Error] {error_msg}')
            import traceback
            traceback.print_exc()
            return audio_pb2.TranscribeReply(
                code=1,
                message=error_msg
            )


