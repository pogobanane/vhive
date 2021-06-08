// MIT License
//
// Copyright (c) 2021 Michal Baczun and EASE lab
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"

	pb "tests/chained-functions-serving/proto"
)

type consumerServer struct {
	pb.UnimplementedProducerConsumerServer
}

func (s *consumerServer) ConsumeString(ctx context.Context, str *pb.ConsumeStringRequest) (*pb.ConsumeStringReply, error) {
	log.Printf("[consumer] Consumed %v\n", str.Value)
	return &pb.ConsumeStringReply{Value: true}, nil
}

func (s *consumerServer) ConsumeStream(stream pb.ProducerConsumer_ConsumeStreamServer) error {
	for {
		str, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.ConsumeStringReply{Value: true})
		}
		if err != nil {
			return err
		}
		log.Printf("[consumer] Consumed %v\n", str.Value)
	}
}

func main() {
	portFlag := flag.Int("p", 3030, "Port")
	flag.Parse()

	//get client port (from env of flag)
	var port int
	portStr, ok := os.LookupEnv("PORT")
	if !ok {
		port = *portFlag
	} else {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatalf("[producer] PORT_CLIENT env variable is not int: %v", err)
		}
	}

	//set up log file
	file, err := os.OpenFile("log/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	//set up server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("[consumer] failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := consumerServer{}
	pb.RegisterProducerConsumerServer(grpcServer, &s)

	log.Println("[consumer] Server Started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[consumer] failed to serve: %s", err)
	}

}
