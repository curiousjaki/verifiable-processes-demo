package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	proto "github.com/ciphersmaug/verifiable-process/zeebe-go-client/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProveRequest(VariableA float64, VariableB float64, Operation string, ImageId []uint32) (float64, []uint32, []byte) {

	var request = proto.ProveRequest{
		VariableA: VariableA,
		VariableB: VariableB,
		Operation: Operation,
		ImageId:   ImageId,
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	r, err := c.ProveExecution(ctx, &request)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	conn.Close()
	return r.ResponseValue, r.ImageId, r.Receipt
}

func proveComputationHandler(client worker.JobClient, job entities.Job) {
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

	var VariableA float64
	if variables["variable_a"] == nil {
		VariableA = 0.0
	} else {
		VariableA = variables["variable_a"].(float64)
	}
	var VariableB float64
	if variables["variable_b"] == nil {
		VariableB = 0.0
	} else {
		VariableB = variables["variable_b"].(float64)
	}

	var Operation string
	if variables["operation"] == nil {
		Operation = "mul"
	} else {
		Operation = variables["operation"].(string)
	}

	var ImageId []uint32
	if headers["image_id"] == "" {
		ImageId = []uint32{0, 1, 2, 3, 4, 5, 6, 7}
		log.Println("Image Id was empty and is now: ", ImageId)
	} else {
		// convert string to []uint32
		ImageId = stringToUint32Array(headers["image_id"])
	}

	var verification_value, imageid, bites = ProveRequest(VariableA, VariableB, Operation, ImageId)
	variables["verification_value"] = verification_value
	//variables["image_id"] = imageid
	variables["receipt"] = base64.StdEncoding.EncodeToString(bites)
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}
	log.Println("Verification Value: ", verification_value, " Image Id: ", imageid)

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("-------------------------------")
	//close(readyClose)
}
