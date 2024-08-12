package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

var (
	addr = flag.String("addr", "localhost:50051", "the risczero address to connect to")
)

func getVariablesAndFlags() (proving_service *bool, verification_service *bool, message_service *bool, zeebe_addr *string) {
	proving_service = flag.Bool("prove", false, "run the proving service")
	verification_service = flag.Bool("verify", false, "run the verification service")
	message_service = flag.Bool("message", false, "run the message service")
	zeebe_addr = flag.String("zeebe", "localhost:26500", "the address of the zeebe cluster")
	flag.Parse()
	println(*proving_service, *verification_service, *message_service, *zeebe_addr)
	return proving_service, verification_service, message_service, zeebe_addr
}
func stringToUint32Array(s string) []uint32 {
	//log.Println(s)
	var uint32Array []uint32
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	for _, si := range strings.Split(s, ",") {
		si = strings.ReplaceAll(si, ",", "")
		si = strings.TrimSpace(si)
		ui, err := strconv.ParseUint(si, 10, 32)
		if err != nil {
			log.Println(err)
			panic(err)
		}
		uint32Array = append(uint32Array, uint32(ui))
	}
	return uint32Array
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
	run_proving_service, run_verification_service, run_message_service, zeebe_addr := getVariablesAndFlags()
	fmt.Println(*zeebe_addr)
	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: *zeebe_addr}
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

	var provingServiceJobWorker, verificationServiceJobWorker, messageServiceJobWorker worker.JobWorker
	if *run_proving_service {
		provingServiceJobWorker = client.NewJobWorker().
			JobType("proving-service").
			Handler(proveComputationHandler).
			MaxJobsActive(1).
			Open()
	}
	if *run_verification_service {
		verificationServiceJobWorker = client.NewJobWorker().
			JobType("verification-service").
			Handler(verifyComputationHandler).
			MaxJobsActive(1).
			Open()
	}
	if *run_message_service {
		messageServiceJobWorker = client.NewJobWorker().
			JobType("message-service").
			Handler(publishMessageHandler(client, "Message-Verifiable-Computation")).
			MaxJobsActive(1).
			Open()
	}

	// Shut down worker when pushing enter in console
	//buf := bufio.NewReader(os.Stdin)
	//buf.ReadLine()

	//fmt.Println("Shutting down...")

	//<-readyClose
	if *run_proving_service {
		//provingServiceJobWorker.Close()
		provingServiceJobWorker.AwaitClose()
	}
	if *run_verification_service {
		//verificationServiceJobWorker.Close()
		verificationServiceJobWorker.AwaitClose()
	}
	if *run_message_service {
		//messageServiceJobWorker.Close()
		messageServiceJobWorker.AwaitClose()
	}
}
