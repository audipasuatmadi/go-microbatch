# Go Microbatch
Microbatch mechanism written in Go. Could be useful in use cases such as:
- Handling frequent real-time data from IoT devices where you have to batch-insert them to your OLAP data store of choice for performance.
- Handling frequent user messages to LLM, where multiple messages of user sent at a short-time is treated as one message for efficiency without destroying semantic meaning.
- Many others.

# Installation
```bash
go get github.com/audipasuatmadi/go-microbatch
```

# The Idea of Microbatching
Microbatching groups multiple incoming data into one batch of data. This means less operation for the data, opting for throughput and sacrificing latency. Thus allows us to have less operations and faster performance on a large amount of data.

<img width="1349" alt="image" src="https://github.com/user-attachments/assets/f0a0f756-98ac-4968-b7d1-2fca8c5dca22" />


Other than for performance benefit, in nowadays, it allows for a unique use case in Large-Language Models. As user tend to send multiple messages from messaging platform such as Whatsapp, we surely don't want each message to be treated as individual message such as the picture below.
<img width="1306" alt="image" src="https://github.com/user-attachments/assets/e7993b34-e9ea-42bc-9ab3-63c02b6bd0e9" />


# Usage
## Example Usage
```go
func main() {
    // Initialize the microbatch
    mb := microbatch.New[string](microbatch.NewParams{
        // When MaxSize is reached, it will flush the batch.
        MaxSize:       int32(10),
        // Or when the interval is reached, it will flush the batch.
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
}

func simulateSendToLLM(mb microbatch.Microbatch[string]) {
    for {
        // This will block until the batch is flushed.
        data := mb.ReadData(context.Background())
        // You got all the data at once, can do whatever you wish with it.
        fmt.Printf("\nsending \"%s\" to LLM\n> ", strings.Join(data, " | "))
    }
}

```

# Contributing
This project is "very" far from perfect. Every contribution is really appreciated.
