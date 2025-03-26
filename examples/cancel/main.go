package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/audipasuatmadi/go-microbatch"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	mb, err := microbatch.New[string](ctx, microbatch.Config[string]{})
	if err != nil {
		log.Fatal(err)
	}

	messages := []string{
		"Hello,",
		"I",
		"Wanna",
		"Ask",
	}

	go mb.Start()
	go simulateSendToLLM(mb)

	// Simulate how people type in message platforms, where people
	// type multiple messages which semantically related as one message.
	for _, message := range messages {
		mb.Add(ctx, message)
		fmt.Printf("\n> %s\n", message)
	}

	// Trigger cancel long no response
	time.AfterFunc(2*time.Second, func() {
		fmt.Println("Cancelling context...")
		cancel()
	})

	// adding batch event already closed
	mb.Add(ctx, "Bye")
}

func simulateSendToLLM(mb *microbatch.Microbatch[string]) {
	for result := range mb.ResultBatch {
		fmt.Printf("\nBatch:\n")
		for _, event := range result {
			fmt.Printf("\n\tEvent: %v\n", event.Payload)
		}
	}
}
