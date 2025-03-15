package microbatch_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/audipasuatmadi/go-microbatch"
	"github.com/stretchr/testify/assert"
)

func Test_ReadsDataWhenReachesBufferSize(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	m, err := microbatch.New[int32](context.Background(), microbatch.Config[int32]{Strategy: &microbatch.SizeBasedStrategy[int32]{MaxSize: 10}})
	if err != nil {
		t.Fatal(err)
	}

	dummyData := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	go m.Start()

	// the reader is in another thread... reading
	go func(m *microbatch.Microbatch[int32]) {
		var allData = []int32{}
		result := <-m.ResultBatch
		for _, event := range result {
			allData = append(allData, event.Payload)
		}
		assert.Equal(t, 10, len(allData))
		assert.Equal(t, []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, allData)

		wg.Done()
	}(m)

	for _, temperature := range dummyData {
		m.Add(context.Background(), temperature)
	}
	wg.Wait()
}

func Test_ReadsDataWhenReachesBufferSizeAndFillsTheRemainingData(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	m, err := microbatch.New[int32](context.Background(), microbatch.Config[int32]{})
	if err != nil {
		t.Fatal(err)
	}

	dummyData := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	go m.Start()

	// the reader is in another thread... reading
	go func(m *microbatch.Microbatch[int32]) {
		var allData = []int32{}
		result := <-m.ResultBatch
		for _, event := range result {
			allData = append(allData, event.Payload)
		}
		fmt.Println(allData)
		assert.Equal(t, 5, len(allData))
		assert.Equal(t, []int32{1, 2, 3, 4, 5}, allData)

		allData = []int32{}
		result = <-m.ResultBatch
		for _, event := range result {
			allData = append(allData, event.Payload)
		}
		assert.Equal(t, 5, len(allData))
		assert.Equal(t, []int32{6, 7, 8, 9, 10}, allData)

		wg.Done()
	}(m)

	for _, temperature := range dummyData {
		m.Add(context.Background(), temperature)
	}

	wg.Wait()
}

// func Test_ThreadSafeOnTheMicrobatch(t *testing.T) {
// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	m := microbatch.New[int32](microbatch.NewParams{
// 		MaxSize:       int32(5),
// 		FlushInterval: time.Hour,
// 	})

// 	dummyData := []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

// 	// the reader is in another thread... reading
// 	go func(m microbatch.Microbatch[int32]) {
// 		dataset1 := m.ReadData(context.Background())
// 		assert.Equal(t, 5, len(dataset1))

// 		dataset2 := m.ReadData(context.Background())
// 		assert.Equal(t, 5, len(dataset2))

// 		var mappedNums map[int32]bool = make(map[int32]bool)
// 		for _, v := range dataset1 {
// 			mappedNums[v] = true
// 		}
// 		for _, v := range dataset2 {
// 			mappedNums[v] = true
// 		}

// 		for _, v := range mappedNums {
// 			if v == false {
// 				assert.Fail(t, "the dataset is not correct")
// 			}
// 		}

// 		wg.Done()
// 	}(m)

// 	for _, temperature := range dummyData {
// 		go m.Add(context.Background(), temperature)
// 	}

// 	wg.Wait()
// }

// func Test_ReadsDataWhenReachesTheFlushInterval(t *testing.T) {
// 	var wg sync.WaitGroup
// 	wg.Add(1)

// 	m := microbatch.New[int32](microbatch.NewParams{
// 		MaxSize:       int32(10),
// 		FlushInterval: 1000 * time.Millisecond,
// 	})

// 	dummyData := []int32{1, 2, 3, 4, 5}
// 	dummyData2 := []int32{6, 7, 8, 9, 10}

// 	// the reader is in another thread... reading
// 	go func(m microbatch.Microbatch[int32]) {
// 		allData := m.ReadData(context.Background())
// 		assert.Equal(t, 5, len(allData))
// 		assert.Equal(t, []int32{1, 2, 3, 4, 5}, allData)

// 		// The next read it should got the remaining data
// 		allData = m.ReadData(context.Background())
// 		assert.Equal(t, 5, len(allData))
// 		assert.Equal(t, []int32{6, 7, 8, 9, 10}, allData)

// 		wg.Done()
// 	}(m)

// 	for _, temperature := range dummyData {
// 		m.Add(context.Background(), temperature)
// 	}

// 	time.Sleep(1010 * time.Millisecond)
// 	for _, temperature := range dummyData2 {
// 		m.Add(context.Background(), temperature)
// 	}

// 	wg.Wait()
// }
