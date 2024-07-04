package main

import (
	"bufio"
	"context"
	"encoding/base64"
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
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

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

//var readyClose = make(chan struct{})

func verifyCarbonEmissionCalculation(client worker.JobClient, job entities.Job) {
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
	var image_id []uint32
	image_id, err = StringToUint32Slice(headers["image_id"])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var receipt []byte
	if variables["receipt"] == "" {
		receipt = []byte{1, 2, 3}
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
	var is_valid_executed = VerifyCarbonEmissionRequest(verification_value, image_id, receipt)
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

func proveCarbonEmissionCalculation(client worker.JobClient, job entities.Job) {
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
		Operation = "multiply"
	} else {
		Operation = variables["operation"].(string)
	}

	var ImageId []uint32
	if headers["image_id"] == "" {
		ImageId = []uint32{1, 2, 3}
	} else {
		ImageId = []uint32{1, 2, 3}
		//ImageId = variables["image_id"].([]uint32)
	}

	var verification_value, imageid, bites = ProveCarbonEmissionRequest(VariableA, VariableB, Operation, ImageId)
	//imageid = append(imageid, 1, 2, 3)
	variables["verification_value"] = verification_value
	variables["imageid"] = imageid
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

func failJob(client worker.JobClient, job entities.Job) {
	log.Println("Failed to complete job", job.GetKey())

	ctx := context.Background()
	_, err := client.NewFailJobCommand().JobKey(job.GetKey()).Retries(job.Retries - 1).Send(ctx)
	if err != nil {
		panic(err)
	}
}

func main() {
	//gatewayAddr := os.Getenv("ZEEBE_ADDRESS")
	//plainText := false

	//if gatewayAddr == "" {
	//	gatewayAddr = ZeebeAddr
	//	plainText = true
	//}
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	client, err := zbc.NewClient(&config)
	if err != nil {
		panic(err)
	}

	// check connection
	ctx := context.Background()
	response, err := client.NewTopologyCommand().Send(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to " + response.GetBrokers()[0].GetHost())

	//zbClient, err := zbc.NewClient(&zbc.ClientConfig{
	//	GatewayAddress:         gatewayAddr,
	//	UsePlaintextConnection: plainText,
	//})

	// deploy process
	//ctx := context.Background()
	//response, err := zbClient.NewDeployResourceCommand().AddResourceFile("order-process.bpmn").Send(ctx)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(response.String())

	// create a new process instance
	//variables := make(map[string]interface{})
	//variables["orderId"] = "31243"

	//request, err := client.NewCreateInstanceCommand().BPMNProcessId("order-process").LatestVersion().VariablesFromMap(variables)
	//if err != nil {
	//	panic(err)
	//}

	//result, err := request.Send(ctx)
	//if err != nil {
	//	panic(err)
	//}

	//fmt.Println(result.String())

	verificationServiceJobWorker := client.NewJobWorker().JobType("verification-service").Handler(verifyCarbonEmissionCalculation).Open()
	provingServiceJobWorker := client.NewJobWorker().JobType("proving-service").Handler(proveCarbonEmissionCalculation).Open()

	// Shut down worker when pushing enter in console
	buf := bufio.NewReader(os.Stdin)
	buf.ReadLine()

	fmt.Println("Shutting down...")

	//<-readyClose
	provingServiceJobWorker.Close()
	provingServiceJobWorker.AwaitClose()
	verificationServiceJobWorker.Close()
	verificationServiceJobWorker.AwaitClose()
}
