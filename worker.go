package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
			"github.com/aws/aws-sdk-go-v2/aws"
)

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest, awsConfig aws.Config) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool),
		AwsConfig:	 awsConfig	}
	return worker
}

type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
	AwsConfig 	aws.Config
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w *Worker ) Start() {
	go func() {


		svc := dynamodb.New(w.AwsConfig)

		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work

			select {
			case work := <-w.Work:
				// Receive a work request.
				fmt.Printf("worker%d: Received insert request %s seconds\n", w.ID, work.record)
				req:= svc.PutItemRequest(work.record)

				resp, err := req.Send()
				if err == nil { // resp is now filled
					fmt.Println(resp)
				}
				fmt.Println("Inserted Record %s ", resp)
			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

