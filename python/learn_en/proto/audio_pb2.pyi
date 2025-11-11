from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SynthesizerRequest(_message.Message):
    __slots__ = ()
    API_KEY_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    VOICE_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    LANGUAGE_TYPE_FIELD_NUMBER: _ClassVar[int]
    api_key: str
    model: str
    voice: str
    text: str
    language_type: str
    def __init__(self, api_key: _Optional[str] = ..., model: _Optional[str] = ..., voice: _Optional[str] = ..., text: _Optional[str] = ..., language_type: _Optional[str] = ...) -> None: ...

class SynthesizerReply(_message.Message):
    __slots__ = ()
    URL_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    url: str
    message: str
    code: int
    def __init__(self, url: _Optional[str] = ..., message: _Optional[str] = ..., code: _Optional[int] = ...) -> None: ...

class TranscribeRequest(_message.Message):
    __slots__ = ()
    AUDIO_PATH_FIELD_NUMBER: _ClassVar[int]
    LANGUAGE_FIELD_NUMBER: _ClassVar[int]
    MODEL_SIZE_FIELD_NUMBER: _ClassVar[int]
    audio_path: str
    language: str
    model_size: str
    def __init__(self, audio_path: _Optional[str] = ..., language: _Optional[str] = ..., model_size: _Optional[str] = ...) -> None: ...

class WordTimestamp(_message.Message):
    __slots__ = ()
    WORD_FIELD_NUMBER: _ClassVar[int]
    START_FIELD_NUMBER: _ClassVar[int]
    END_FIELD_NUMBER: _ClassVar[int]
    CONFIDENCE_FIELD_NUMBER: _ClassVar[int]
    word: str
    start: float
    end: float
    confidence: float
    def __init__(self, word: _Optional[str] = ..., start: _Optional[float] = ..., end: _Optional[float] = ..., confidence: _Optional[float] = ...) -> None: ...

class TranscribeReply(_message.Message):
    __slots__ = ()
    TEXT_FIELD_NUMBER: _ClassVar[int]
    WORDS_FIELD_NUMBER: _ClassVar[int]
    LANGUAGE_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    text: str
    words: _containers.RepeatedCompositeFieldContainer[WordTimestamp]
    language: str
    code: int
    message: str
    def __init__(self, text: _Optional[str] = ..., words: _Optional[_Iterable[_Union[WordTimestamp, _Mapping]]] = ..., language: _Optional[str] = ..., code: _Optional[int] = ..., message: _Optional[str] = ...) -> None: ...
