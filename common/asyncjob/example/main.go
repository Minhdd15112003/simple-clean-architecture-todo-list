package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"social-todo-list/common/asyncjob"
)

func exameJob() {
	job := asyncjob.NewJob(func(ctx context.Context) error {
		fmt.Println("JOB")
		return nil
	}, asyncjob.WithName("Job"))

	if err := job.Execute(context.Background()); err != nil {
		log.Println(err)

		for {
			err := job.Retry(context.Background())

			if err == nil || job.State() == asyncjob.StateRetryFailed {
				break
			}
		}
	}
}

func exameGroupJob() {
	job1 := asyncjob.NewJob(func(ctx context.Context) error {
		return errors.New("ccccccccccc")
	}, asyncjob.WithName("Job 1"))

	job2 := asyncjob.NewJob(func(ctx context.Context) error {
		return nil
	}, asyncjob.WithName("Job 2"))

	err := asyncjob.NewGroup(false, job1, job2).Run(context.Background())

	if err != nil {
		log.Println("err: ", err)
	}
}

func main() {
	exameGroupJob()
}
