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
	mb, err := microbatch.New[string](context.Background(), microbatch.Config[string]{})
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
	mb.Add(ctx, messages[0])
	fmt.Printf("\n> %s\n", messages[0])
	time.Sleep(2 * time.Second)
	mb.Add(ctx, messages[1])
	fmt.Printf("\n> %s\n", messages[1])

	// Input as many messages as you want
	var userMessage string
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		userMessage, _ = in.ReadString('\n')

		if userMessage == "exit\n" {
			mb.Stop()
			break
		}

		err := mb.Add(context.Background(), strings.TrimSpace(userMessage))
		if err != nil {
			log.Println("error:", err)
		}
	}
}

func simulateSendToLLM(mb *microbatch.Microbatch[string]) {
	for result := range mb.ResultBatch {
		fmt.Printf("\nBatch:\n")
		for _, event := range result {
			fmt.Printf("\n\tEvent: %v\n", event.Payload)
		}
	}
}
