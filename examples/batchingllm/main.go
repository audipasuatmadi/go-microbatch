package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/audipasuatmadi/go-microbatch"
)

func main() {
	mb := microbatch.New[string](microbatch.NewParams{
		MaxSize:       int32(10),
		FlushInterval: 3 * time.Second,
	})

	messages := []string{
		"Hello",
		"How are you?",
	}

	go simulateSendToLLM(mb)

	// Simulate how people type in message platforms, where people
	// type multiple messages which semantically related as one message.
	mb.Add(context.Background(), messages[0])
	fmt.Printf("\n> %s\n", messages[0])
	time.Sleep(2 * time.Second)
	mb.Add(context.Background(), messages[1])
	fmt.Printf("\n> %s\n", messages[1])

	// Input as many messages as you want
	var userMessage string
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		userMessage, _ = in.ReadString('\n')

		mb.Add(context.Background(), strings.TrimSpace(userMessage))
	}
}

func simulateSendToLLM(mb microbatch.Microbatch[string]) {
	for {
		data := mb.ReadData(context.Background())
		fmt.Printf("\nsending \"%s\" to LLM\n> ", strings.Join(data, " | "))
	}
}
