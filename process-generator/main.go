package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

const ZeebeAddr = "0.0.0.0:26500"

func main() {

	var nFlag = flag.Int("n", 10, "number of process instances to generate")
	var tFlag = flag.Int("t", 10, "seconds to wait between process instance generation")
	flag.Parse()

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
	ctx := context.Background()
	response, err := zbClient.NewDeployResourceCommand().AddResourceFile("prove-verify-service.bpmn").Send(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.String())
	go generateLoop(nFlag, tFlag, client, ctx)
	// create a new process instance
	// Shut down worker when pushing enter in console
	buf := bufio.NewReader(os.Stdin)
	buf.ReadLine()

	fmt.Println("Shutting down...")
}

func generateLoop(nFlag *int, tFlag *int, client zbc.Client, ctx context.Context) {
	for i := 0; i < *nFlag; i++ {
		variables := make(map[string]interface{})
		variables["variable_a"] = rand.Float64() * (rand.Float64() * 100)
		variables["variable_b"] = rand.Float64() * (rand.Float64() * 100)
		variables["operation"] = "multiply"
		go generateProcessInstance(client, ctx, "verification-process", variables)
		time.Sleep(time.Duration(*tFlag) * time.Second)
	}
}

func generateProcessInstance(client zbc.Client, ctx context.Context, processId string, variables map[string]interface{}) {
	request, err := client.NewCreateInstanceCommand().BPMNProcessId(processId).LatestVersion().VariablesFromMap(variables)
	if err != nil {
		panic(err)
	}
	result, err := request.Send(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.String())
}
