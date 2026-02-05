package asyncjob

import (
	"context"
	"social-todo-list/middleware"
	"sync"
)

type group struct {
	jobs         []Job
	isConcurrent bool
	wg           *sync.WaitGroup // wg để đếm số job đã hoàn thành trong nhóm
}

func NewGroup(isConcurrent bool, jobs ...Job) *group {
	g := &group{
		isConcurrent: isConcurrent,
		jobs:         jobs,
		wg:           new(sync.WaitGroup), // khởi tạo WaitGroup
	}

	return g
}

func (g *group) Run(ctx context.Context) error {
	g.wg.Add(len(g.jobs)) // nhận vào số job cần hoàn thành

	errChan := make(chan error, len(g.jobs)) // kênh để nhận lỗi từ các job

	for i := range g.jobs {
		if g.isConcurrent {
			go func(aj Job) {
				defer middleware.RecoverGoroutine()

				errChan <- g.runJob(ctx, aj) // gửi lỗi (nếu có) vào kênh
				g.wg.Done()
			}(g.jobs[i])
			continue

		}

		job := g.jobs[i]

		err := g.runJob(ctx, job)
		errChan <- err // gửi vào channel TRƯỚC
		if err != nil {
			return err // return SAU
		}

		g.wg.Done() // đánh dấu một job đã hoàn thành
	}
	g.wg.Wait() // chờ tất cả các job hoàn thành

	var err error
	for i := 1; i <= len(g.jobs); i++ {
		if v := <-errChan; v != nil { // nếu luống có lỗi thì rút ra và chả về
			return v
		}
	}
	return err
}

func (g *group) runJob(ctx context.Context, j Job) error {
	if err := j.Execute(ctx); err != nil {
		for {
			// log.Println(err)
			if j.State() == StateRetryFailed {
				return err
			}
			if j.Retry(ctx) == nil {
				return nil
			}
		}
	}
	return nil
}
