import grpc

from protos import helloworld_pb2_grpc, helloworld_pb2


def run():
    # 连接服务器
    channel = grpc.insecure_channel('localhost:50051')
    stub = helloworld_pb2_grpc.GreeterStub(channel)

    # 调用 RPC 方法
    response = stub.SayHello(helloworld_pb2.HelloRequest(name="ChatGPT"))
    print("Greeter client received:", response.message)

if __name__ == "__main__":
    run()
