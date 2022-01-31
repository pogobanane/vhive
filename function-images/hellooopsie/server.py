from concurrent import futures
import logging
import random
import time
import threading

import grpc

import helloworld_pb2
import helloworld_pb2_grpc

responses = ["record_response", "replay_response"]

ERROR_NEEDED = random.randrange(0, 10)

class Greeter(helloworld_pb2_grpc.GreeterServicer):

    def SayHello(self, request, context):
        global ERROR_NEEDED
        if request.name == "record":
            msg = 'Hello, %s!' % responses[0]
        elif request.name == "replay":
            msg = 'Hello, %s!' % responses[1]
        elif request.name == "satan":
            print("ERROR: oopsie!", flush=True)
            msg = 'Hello, oopsie!'
        else:
            if ERROR_NEEDED == 9:
                ERROR_NEEDED = 0
                print("ERROR: random error", flush=True)
                msg = 'Hello, RandomError(%s)!' % request.name
            else:
                msg = 'Hello, %s!' % request.name

        ERROR_NEEDED += 1
        return helloworld_pb2.HelloReply(message=msg)


class PeriodicError(threading.Thread):
    def __init__(self):
        threading.Thread.__init__(self)

    def run(self):
        time.sleep(1);
        while True:
            print("ERROR: some periodic error", flush=True);
            time.sleep(10);


def serve():
    errors = PeriodicError()
    errors.start()
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
    helloworld_pb2_grpc.add_GreeterServicer_to_server(Greeter(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()
    errors.join()


if __name__ == '__main__':
    print("ERROR: some virtual startup error", flush=True);
    logging.basicConfig()
    serve()
