package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Job interface {
	Run() error
}

type Worker struct {
	ID   int
	Jobs chan Job
}

func NewWorker() *Worker {
	jobs := make(chan Job, 10)
	worker := &Worker{
		Jobs: jobs,
	}

	return worker
}

func (w *Worker) Start() {
	for job := range w.Jobs {
		fmt.Printf("Worker %d processing...\n", w.ID)
		err := job.Run()
		fmt.Printf("Worker %d completed!!\n", w.ID)
		if err != nil {
			fmt.Println("Job error:", err)
		}
	}
}

type WorkerPool struct {
	Jobs    chan Job
	Workers []*Worker
	closed  bool
	mu      sync.Mutex
}

func NewWorkerPool(workerCount int) *WorkerPool {
	jobs := make(chan Job, 10)
	wp := &WorkerPool{
		Jobs:    jobs,
		Workers: make([]*Worker, workerCount),
	}

	for i := 0; i < workerCount; i++ {
		worker := &Worker{
			ID:   i + 1,
			Jobs: jobs,
		}
		wp.Workers[i] = worker
		go worker.Start()
	}

	return wp
}

func (wp *WorkerPool) Enqueue(job Job) {
	if wp.closed {
		fmt.Println("Cannot enqueue, the pool was closed")
		return
	}
	for {
		select {
		case wp.Jobs <- job:
			return
		default:
			// Log or handle full channel case
			// fmt.Println("Job queue is full")
			// Full job channel
			time.Sleep(time.Second)
		}
	}
}

func (wp *WorkerPool) Close() {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	if wp.closed {
		return
	}

	wp.closed = true
	close(wp.Jobs)
	fmt.Println("Worker pool is closed successfully")
}

type TestJob struct{}

func NewTestJob() TestJob {
	return TestJob{}
}

func (tj TestJob) Run() error {
	time.Sleep(10 * time.Second)
	// time.Sleep(time.Duration(rand.Int64N(int64(5)) * int64(time.Second)))
	fmt.Println("Test job run background!!!")
	return nil
}

func main() {
	wp := NewWorkerPool(3)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testJob := NewTestJob()
		go wp.Enqueue(testJob)
		w.Write([]byte("Test successfully"))
	})

	http.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		wp.Close()
	})

	addr := ":8081"
	fmt.Println("Server is listening at: ", addr)
	http.ListenAndServe(addr, nil)
}
