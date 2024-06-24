package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

const ZeebeAddr = "0.0.0.0:26500"

//var readyClose = make(chan struct{})

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

	jobWorker := client.NewJobWorker().JobType("payment-service").Handler(handleJob).Open()

	// Shut down worker when pushing enter in console
	buf := bufio.NewReader(os.Stdin)
	buf.ReadLine()

	fmt.Println("Shutting down...")

	//<-readyClose
	jobWorker.Close()
	jobWorker.AwaitClose()
}

func handleJob(client worker.JobClient, job entities.Job) {
	jobKey := job.GetKey()

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

	variables["totalPrice"] = 46.50
	request, err := client.NewCompleteJobCommand().JobKey(jobKey).VariablesFromMap(variables)
	if err != nil {
		// failed to set the updated variables
		failJob(client, job)
		return
	}

	log.Println("Complete job", jobKey, "of type", job.Type)
	log.Println("Processing order:", variables["orderId"])
	log.Println("Collect money using payment method:", headers["method"])

	ctx := context.Background()
	_, err = request.Send(ctx)
	if err != nil {
		panic(err)
	}

	log.Println("Successfully completed job")
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
