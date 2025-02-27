# Go Microbatch
Microbatch mechanism written in Go. Could be useful in use cases such as:
- Handling frequent real-time data from IoT devices where you have to batch-insert them to your OLAP data store of choice for performance.
- Handling frequent user messages to LLM, where multiple messages of user sent at a short-time is treated as one message for efficiency without destroying semantic meaning.
- Many others.

# Installation
```bash
go get github.com/audipasuatmadi/go-microbatch
```

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