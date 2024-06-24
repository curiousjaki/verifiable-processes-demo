package main

import (
	"context"
	"flag"
	"log"
	"time"

	zeebe_go_client "github.com/ciphersmaug/verifiable-process/zeebe-go-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

var request = zeebe_go_client.EmissionProofRequest{
	Consumption:    2,
	EmissionFactor: 10,
}

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//log.Printf("{}", conn)

	c := zeebe_go_client.NewCarbonEmissionClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.ProveCarbonEmission(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Co2Emissions)
}
