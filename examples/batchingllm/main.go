package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/audipasuatmadi/go-microbatch"
)

func main() {
	mb, err := microbatch.New[string](microbatch.Context{
		Ctx: context.Background(),
	}, microbatch.Config[string]{
		Strategy: &microbatch.SizeBasedStrategy[string]{
			MaxSize: 3,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	messages := []string{
		"Hello",
		"How are you?",
	}

	go mb.Start()
	go simulateSendToLLM(mb)

	// Simulate how people type in message platforms, where people
	// type multiple messages which semantically related as one message.
	ctx := context.Background()
	mb.Add(ctx, microbatch.Event[string]{
		Payload: messages[0],
		AddedAt: time.Now(),
	})
	fmt.Printf("\n> %s\n", messages[0])
	time.Sleep(2 * time.Second)
	mb.Add(ctx, microbatch.Event[string]{
		Payload: messages[1],
		AddedAt: time.Now(),
	})
	fmt.Printf("\n> %s\n", messages[1])

	// Input as many messages as you want
	var userMessage string
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		userMessage, _ = in.ReadString('\n')

		mb.Add(context.Background(), microbatch.Event[string]{
			Payload: strings.TrimSpace(userMessage),
			AddedAt: time.Now(),
		})
	}
}

func simulateSendToLLM(mb *microbatch.Microbatch[string]) {
	for result := range mb.ResultStream {
		fmt.Printf("\nBatch:\n")
		for _, event := range result {
			fmt.Printf("\n\tEvent: %v\n", event.Event.Payload)
		}
	}
}
