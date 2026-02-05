package asyncjob

import (
	"context"
	"log"
	"time"
)

//task degien parten

type Job interface {
	Execute(ctx context.Context) error
	Retry(ctx context.Context) error
	State() JobState
	SetRetryDurations(time []time.Duration)
}

const (
	defaultMaxTimeout = time.Second * 10
)

var (
	defaultRetryTime = []time.Duration{time.Second, time.Second * 2, time.Second * 4}
)

type JobHandler func(context.Context) error

type JobState int

const (
	StateInit JobState = iota
	StateRuning
	StateFailed
	StateTimeout
	StateCompleted
	StateRetryFailed
)

func (js JobState) String() string {
	return [6]string{
		"Init",
		"Runing",
		"Failed",
		"Timeout",
		"Completed",
	}[js]
}

type jobConfig struct {
	Name       string
	MaxTimeout time.Duration
	Retries    []time.Duration
}

type job struct {
	config     jobConfig
	handler    JobHandler
	state      JobState
	retryIndex int
	stopChan   chan bool
}

func NewJob(handler JobHandler, option ...OptionHdl) *job {
	j := job{
		config: jobConfig{
			MaxTimeout: defaultMaxTimeout,
			Retries:    defaultRetryTime,
		},
		handler:    handler,
		retryIndex: -1,
		state:      StateInit,
		stopChan:   make(chan bool),
	}

	for i := range option {
		option[i](&j.config)
	}
	return &j
}

func (j *job) Execute(ctx context.Context) error {
	log.Printf("Execute %s\n", j.config.Name)
	j.state = StateRuning

	var err error
	err = j.handler(ctx)

	if err != nil {
		j.state = StateFailed
		return err
	}
	j.state = StateCompleted

	return nil
}

func (j *job) Retry(ctx context.Context) error {

	if j.retryIndex == len(j.config.Retries)-1 {
		return nil
	}

	j.retryIndex += 1
	time.Sleep(j.config.Retries[j.retryIndex])

	err := j.Execute(ctx)
	if err == nil {
		j.state = StateCompleted
		return nil
	}

	if j.retryIndex == len(j.config.Retries)-1 {
		j.state = StateFailed
		return err
	}

	j.state = StateFailed
	return err
}

func (j *job) State() JobState { return j.state }
func (j *job) RetryIndex() int { return j.retryIndex }
func (j *job) SetRetryDurations(times []time.Duration) {
	if len(times) == 0 {
		return
	}

	j.config.Retries = times
}

type OptionHdl func(*jobConfig)

func WithName(name string) OptionHdl {
	return func(jc *jobConfig) {
		jc.Name = name
	}
}
func WithRetriesDuration(items []time.Duration) OptionHdl {
	return func(jc *jobConfig) {
		jc.Retries = items
	}
}
