package main

import (
	"context"
	"log"
	"time"

	proto "github.com/ciphersmaug/verifiable-process/zeebe-go-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProveCarbonEmissionRequest(VariableA float64, VariableB float64, Operation string, ImageId []uint32) (float64, []uint32, []byte) {

	var request = proto.ProveRequest{
		VariableA: float64(VariableA),
		VariableB: float64(VariableB),
		Operation: Operation,
		ImageId:   []uint32{1, 2, 3},
	}

	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//log.Printf("{}", conn)

	c := proto.NewVerifiableProcessingServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.ProveExecution(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	conn.Close()
	return r.ResponseValue, r.ImageId, r.Receipt
}
