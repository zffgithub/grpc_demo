import logging

import grpc

import cyber_pb2
import cyber_pb2_grpc

def run():
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = cyber_pb2_grpc.CyberManagerStub(channel)
        responses = stub.CyberControl(cyber_pb2.Request(cyber_id=2))
        for response in responses:
            print("调用成功: {}!".format(response))

if __name__ == '__main__':
    logging.basicConfig()
    run()