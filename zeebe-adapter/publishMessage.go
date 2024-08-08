package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func publishMessageHandler(client zbc.Client, messageName string) func(jobClient worker.JobClient, job entities.Job) {
	return func(jobClient worker.JobClient, job entities.Job) {
		jobKey := job.GetKey()
		variables, err := job.GetVariablesAsMap()
		if err != nil {
			// failed to handle job as we require the variables
			failJob(client, job)
			return
		}

		log.Println("Processing job", jobKey, "of type", job.Type)

		// Publish a message (using the main client)
		err = publishMessage(client, job, messageName, variables)
		if err != nil {
			log.Printf("Failed to publish message: %v", err)
		} else {
			fmt.Println("Message published successfully")
		}

		// Complete the job
		_, err = jobClient.NewCompleteJobCommand().JobKey(job.Key).Send(context.Background())
		if err != nil {
			log.Printf("Failed to complete job: %v", err)
		} else {
			fmt.Printf("Job %d completed successfully\n", job.Key)
		}
	}
}

// publishMessage sends a message to Zeebe
func publishMessage(zeebeClient zbc.Client, job entities.Job, messageName string, variables map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create the message request
	messageRequest, err := zeebeClient.NewPublishMessageCommand().
		MessageName(messageName).
		CorrelationKey(string(job.GetKey())). // Use job ID or other relevant correlation key
		VariablesFromMap(variables)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	messageRequest.Send(ctx)

	//	fmt.Printf("Published message with key: %d\n", messageRequest.GetMessageKey())
	return nil
}
