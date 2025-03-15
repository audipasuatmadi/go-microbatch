package microbatch

import "time"

var (
	defaultBatchTimeoutDuration = 5 * time.Second
	defaultBatchMaxSize         = 5
)
