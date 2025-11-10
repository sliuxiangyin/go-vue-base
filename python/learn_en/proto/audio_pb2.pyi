from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class SynthesizerRequest(_message.Message):
    __slots__ = ()
    API_KEY_FIELD_NUMBER: _ClassVar[int]
    MODEL_FIELD_NUMBER: _ClassVar[int]
    VOICE_FIELD_NUMBER: _ClassVar[int]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    api_key: str
    model: str
    voice: str
    text: str
    def __init__(self, api_key: _Optional[str] = ..., model: _Optional[str] = ..., voice: _Optional[str] = ..., text: _Optional[str] = ...) -> None: ...

class SynthesizerReply(_message.Message):
    __slots__ = ()
    DATA_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    CODE_FIELD_NUMBER: _ClassVar[int]
    data: bytes
    message: str
    code: int
    def __init__(self, data: _Optional[bytes] = ..., message: _Optional[str] = ..., code: _Optional[int] = ...) -> None: ...
