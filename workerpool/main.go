package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)


type SendEmailJob interface {
	Run(ctx context.Context) error
}

type sendEmailJob struct {
	email    string
	isFailed bool
}

func NewSendEmailJob(email string, isFailed bool) SendEmailJob {
	return &sendEmailJob{
		email:    email,
		isFailed: isFailed,
	}
}

func (j sendEmailJob) Run(ctx context.Context) error {
	// time.Sleep(time.Duration(rand.Int64N(int64(5)) * int64(time.Second)))
	fmt.Printf("Sending to %s...\n", j.email)
	if j.isFailed {
		return fmt.Errorf("error send mail")
	}
	time.Sleep(10 * time.Second)
	fmt.Printf("Sending to %s successfully!\n", j.email)
	return nil
}

func main() {
	wp := NewWorkerPool(3)
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	ctx := context.Background()
	wp.Start(ctx)
	defer wp.Close()

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		job := NewSendEmailJob("test@gmail.com", false)
		go wp.AddJob(NewJob("SendEmail", job.Run))

		w.Write([]byte("Test successfully"))
	})

	http.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) {
		job := NewSendEmailJob("test@gmail.com", true)
		go wp.AddJob(NewJob("SendEmail", job.Run, SetAttempt(3)))

		w.Write([]byte("Test successfully"))
	})

	addr := ":8081"
	fmt.Println("Server is listening at: ", addr)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println(err)
	}
}
