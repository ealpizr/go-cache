package main

import (
	"sync"
	"time"

	"github.com/ealpizr/go-cache/src/types"
)

func ExpensiveFunction(value interface{}) types.FunctionResult {
	time.Sleep(5 * time.Second)
	return types.FunctionResult{value, nil}
}

func main() {
	nWorkers := 3
	var wg sync.WaitGroup
	jobs := []*types.Job{
		types.NewJob(0),
		types.NewJob(1),
		types.NewJob(2),
		types.NewJob(3),
		types.NewJob(0),
		types.NewJob(1),
		types.NewJob(0),
		types.NewJob(2),
		types.NewJob(3),
		types.NewJob(1),
	}
	j := make(chan types.Job, nWorkers)
	fc := types.NewFunctionCache(ExpensiveFunction)

	wg.Add(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go func(i int, j chan types.Job) {
			w := types.NewWorker(i, j, &wg, fc)
			w.Work()
		}(i, j)
	}

	for _, job := range jobs {
		j <- *job
	}
	close(j)

	wg.Wait()
}
