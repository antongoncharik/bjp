package jobprocessor

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

type Job struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type JobHandler func(job *Job) error

type JobProcessor struct {
	redisPool *redis.Pool
	jobQueue  string
	handlers  map[string]JobHandler
}

func New(redisAddr, jobQueue string) *JobProcessor {
	return &JobProcessor{
		redisPool: &redis.Pool{
			MaxIdle:   10,
			MaxActive: 10,
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", redisAddr)
			},
		},
		jobQueue: jobQueue,
		handlers: make(map[string]JobHandler),
	}
}

func (jp *JobProcessor) RegisterHandler(jobType string, handler JobHandler) {
	jp.handlers[jobType] = handler
}

func (jp *JobProcessor) EnqueueJob(job Job) error {
	conn := jp.redisPool.Get()
	defer conn.Close()

	jobData, err := json.Marshal(job)
	if err != nil {
		return err
	}

	_, err = conn.Do("RPUSH", jp.jobQueue, jobData)
	if err != nil {
		return err
	}

	return nil
}

func (jp *JobProcessor) Start() {
	fmt.Println("Job Processor started")
	for {
		job, err := jp.FetchJob()
		if err != nil {
			log.Printf("Failed to fetch job: %v\n", err)
			continue
		}

		if handler, exists := jp.handlers[job.Type]; exists {
			if err := handler(job); err != nil {
				log.Printf("Failed to process job: %v\n", err)
			}
		} else {
			log.Printf("No handler registered for job type: %s\n", job.Type)
		}
	}
}

func (jp *JobProcessor) FetchJob() (*Job, error) {
	conn := jp.redisPool.Get()
	defer conn.Close()

	jobData, err := redis.Bytes(conn.Do("BRPOP", jp.jobQueue, 0))
	if err != nil {
		return nil, err
	}

	var job Job
	err = json.Unmarshal(jobData, &job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
