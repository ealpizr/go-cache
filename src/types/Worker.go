package types

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id int
	j  chan Job
	fc *FunctionCache
	wg *sync.WaitGroup
}

func NewWorker(ID int, j chan Job, wg *sync.WaitGroup, fc *FunctionCache) *Worker {
	return &Worker{ID, j, fc, wg}
}

func (w *Worker) Work() {
	defer w.wg.Done()
	fmt.Printf("[worker %d] waiting for jobs\n", w.id)
	for job := range w.j {
		start := time.Now()
		fr, exists := w.fc.Get(job.value)
		if !exists {
			fmt.Printf("[worker %d] key for %d missing\n", w.id, job.value)
			if w.fc.InProgress[job] {
				fmt.Printf("[worker %d] job for %d is already being calculated. waiting\n", w.id, job.value)
				c := make(chan FunctionResult)
				w.fc.Waiting[job] = append(w.fc.Waiting[job], c)
				fr = <-c
				fmt.Printf("[worker %d] function result for job %d received\n", w.id, job.value)
			} else {
				fmt.Printf("[worker %d] calculating for job %d\n", w.id, job.value)
				w.fc.InProgress[job] = true
				fr = w.fc.f(job)
				w.fc.Set(job.value, fr)
				for _, c := range w.fc.Waiting[job] {
					c <- fr
				}
				w.fc.InProgress[job] = false
			}
		} else {
			fmt.Printf("[worker %d] key for %d found\n", w.id, job.value)
		}
		fmt.Printf("[worker %d] finished job for %d. Took %s\n", w.id, fr.Value, time.Since(start))
	}

}
