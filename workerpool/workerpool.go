package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type ProcessFunc func(ctx context.Context) error

type Job struct {
	key     string
	process ProcessFunc
	retried uint
	attempt uint

	retryDelay time.Duration
}

type JobOption func(j *Job)

func NewJob(key string, process ProcessFunc, opts ...JobOption) Job {
	job := Job{
		key:        key,
		process:    process,
		retryDelay: 100 * time.Millisecond,
	}

	for _, opt := range opts {
		opt(&job)
	}

	return job
}

func SetAttempt(n uint) JobOption {
	return func(j *Job) {
		j.attempt = n
	}
}
func SetRetryDelay(delayTime time.Duration) JobOption {
	return func(j *Job) {
		j.retryDelay = delayTime
	}
}

type WorkerPool struct {
	jobQueue  chan Job
	workerNum int
	closed    bool
	mu        sync.Mutex
}

func NewWorkerPool(workerNum int) *WorkerPool {
	jobQueue := make(chan Job, 1000)

	wp := &WorkerPool{
		jobQueue:  jobQueue,
		workerNum: workerNum,
	}

	return wp
}

func exponentialBackoff(n uint, delayTime time.Duration) time.Duration {
	const maxShift = 62

	if n > maxShift {
		n = maxShift - n
	}

	delayTime = delayTime << n
	return delayTime
}

func (wp *WorkerPool) processJob(ctx context.Context, job Job, workerID int) {
	retryJob := func(job Job, err error) {
		if job.retried < job.attempt {
			go func() {
				fmt.Printf("Worker %d encountered an error processing job %s: %v. Retrying...\n", workerID, job.key, err)
				job.retried++
				time.Sleep(exponentialBackoff(job.retried, job.retryDelay)) // exponential backoff delay
				wp.AddJob(job)
			}()
		} else {
			fmt.Printf("Worker %d cannot process job %s: %v.\n", workerID, job.key, err)
		}
	}

	if job.retried > 0 {
		fmt.Printf("Worker %d retrying job %s (attempt %d/%d)...\n", workerID, job.key, job.retried, job.attempt)
	} else {
		fmt.Printf("Worker %d starting job %s...\n", workerID, job.key)
	}
	err := job.process(ctx)
	if err != nil {
		retryJob(job, err)
	} else {
		fmt.Printf("Worker %d finished job!!\n", workerID)

	}

}

func (wp *WorkerPool) startWorker(ctx context.Context, id int) {
	for {
		select {
		case job, ok := <-wp.jobQueue:
			if !ok { // queue is closed
				continue
			}
			wp.processJob(ctx, job, id)
		case <-ctx.Done():
			fmt.Printf("Worker %d received shutdown signal\n", id)
			return
		}
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < int(wp.workerNum); i++ {
		go wp.startWorker(ctx, i)
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	if wp.closed {
		fmt.Println("Cannot add job, the pool was closed")
		return
	}

	retryInterval := time.Millisecond * 100
	for {
		select {
		case wp.jobQueue <- job:
			return
		default:
			// Log or handle full channel case
			fmt.Println("Job queue is full")
			time.Sleep(retryInterval)
			retryInterval *= 2
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
	close(wp.jobQueue)
	fmt.Println("Worker pool is closed successfully")
}
