package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	proto "github.com/ciphersmaug/verifiable-process/zeebe-go-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//var (
//	addr = flag.String("addr", "localhost:50051", "the address to connect to")
//)

func VerifyRequest(verification_value float64, image_id []uint32, receipt []byte) bool {

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

func verifyComputationHandler(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()
	log.Println("Processing job", jobKey, "of type", job.Type)

	headers, err := job.GetCustomHeadersAsMap()
	if err != nil {
		// failed to handle job as we require the custom job headers
		failJob(client, job)
		return
	}

	variables, err := job.GetVariablesAsMap()
	if err != nil {
		// failed to handle job as we require the variables
		failJob(client, job)
		return
	}

	var verification_value float64
	if variables["verification_value"] == nil {
		verification_value = 0.0
	} else {
		verification_value = variables["verification_value"].(float64)
	}
	var image_id []uint32 = stringToUint32Array(headers["image_id"])
	var receipt []byte
	if variables["receipt"] == "" {
		receipt = []byte{1, 2, 3}
		log.Println("Receipt was empty")
	} else {
		receipt, err = base64.StdEncoding.DecodeString(variables["receipt"].(string))
		if err != nil {
			fmt.Println("Error decoding:", err)
			return
		}
	}

	// The following writes the value variables[receipt] to a file

	if err := os.WriteFile("receipt.txt", receipt, 0666); err != nil {
		log.Fatal(err)
	}

	//var image_id []float64
	//var receipt []byte

	//var emission_factor = i.(uint32)
	var is_valid_executed = VerifyRequest(verification_value, image_id, receipt)
	//imageid = append(imageid, 1, 2, 3)
	variables["is_valid_executed"] = is_valid_executed
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}
	log.Println("Headers: ", verification_value, headers, is_valid_executed)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("-------------------------------")
	//close(readyClose)
}
