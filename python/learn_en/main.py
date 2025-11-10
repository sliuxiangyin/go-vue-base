from concurrent import futures
import time
import grpc

from handler.audio import Audio
from proto import  audio_pb2_grpc

# 实现 Greeter 服务

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    audio_pb2_grpc.add_AudioServicer_to_server(Audio(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server started on port 50051")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()