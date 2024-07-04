package main

import (
	"context"
	"log"
	"time"

	proto "github.com/ciphersmaug/verifiable-process/zeebe-go-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//var (
//	addr = flag.String("addr", "localhost:50051", "the address to connect to")
//)

func VerifyCarbonEmissionRequest(verification_value float64, image_id []uint32, receipt []byte) bool {

	var request = proto.VerifyRequest{
		VerificationValue: float64(verification_value),
		ImageId:           []uint32(image_id),
		Receipt:           []byte(receipt),
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

	r, err := c.VerifyExecution(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	conn.Close()
	return r.IsValidExecuted
}
