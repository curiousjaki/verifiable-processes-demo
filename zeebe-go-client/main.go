package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

const ZeebeAddr = "0.0.0.0:26500"

var (
	addr = flag.String("addr", "localhost:50051", "the risczero address to connect to")
)

func getVariablesAndFlags() (proving_service *bool, verification_service *bool) {
	proving_service = flag.Bool("run-proving-service", false, "run the proving service")
	verification_service = flag.Bool("run-verification-service", false, "run the verification service")
	flag.Parse()
	println(*proving_service, *verification_service)
	return proving_service, verification_service
}

// StringToUint32Slice converts a string representation of a slice of uint32s into an actual []uint32.
func StringToUint32Slice(s string) ([]uint32, error) {
	// Remove the square brackets from the string
	s = strings.Trim(s, "[]")

	// Split the string into individual number strings
	numberStrings := strings.Fields(s)

	// Create a slice to hold the uint32 values
	uint32Slice := make([]uint32, len(numberStrings))

	// Convert each number string to a uint32 and store it in the slice
	for i, numStr := range numberStrings {
		num, err := strconv.ParseUint(numStr, 10, 32)
		if err != nil {
			return nil, err
		}
		uint32Slice[i] = uint32(num)
	}

	return uint32Slice, nil
}

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}

func main() {
	run_proving_service, run_verification_service := getVariablesAndFlags()
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	client, err := zbc.NewClient(&config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	response, err := client.NewTopologyCommand().Send(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to " + response.GetBrokers()[0].GetHost())

	var provingServiceJobWorker, verificationServiceJobWorker worker.JobWorker
	if *run_proving_service {
		provingServiceJobWorker = client.NewJobWorker().
			JobType("proving-service").
			Handler(proveCarbonEmissionCalculation).
			MaxJobsActive(1).
			Open()
	}
	if *run_verification_service {
		verificationServiceJobWorker = client.NewJobWorker().
			JobType("verification-service").
			Handler(verifyCarbonEmissionCalculation).
			MaxJobsActive(1).
			Open()
	}

	// Shut down worker when pushing enter in console
	buf := bufio.NewReader(os.Stdin)
	buf.ReadLine()

	fmt.Println("Shutting down...")

	//<-readyClose
	if *run_proving_service {
		provingServiceJobWorker.Close()
		provingServiceJobWorker.AwaitClose()
	}
	if *run_verification_service {
		verificationServiceJobWorker.Close()
		verificationServiceJobWorker.AwaitClose()
	}
}
