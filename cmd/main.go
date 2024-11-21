package main

import (
	"fmt"
	"log"

	"github.com/antongoncharik/bjp/example"
	"github.com/antongoncharik/bjp/jobprocessor"
)

func main() {
	redisAddr := ":6379"
	jobQueue := "jobs"
	processor := jobprocessor.New(redisAddr, jobQueue)

	processor.RegisterHandler("send_email", example.SendEmailHandler)
	processor.RegisterHandler("generate_report", example.GenerateReportHandler)

	go processor.Start()

	for i := 1; i <= 10; i++ {
		emailJob := jobprocessor.Job{Type: "send_email", Data: fmt.Sprintf("user%d@example.com", i)}
		reportJob := jobprocessor.Job{Type: "generate_report", Data: fmt.Sprintf("report%d", i)}

		if err := processor.EnqueueJob(emailJob); err != nil {
			log.Fatalf("Failed to enqueue job: %v", err)
		}

		if err := processor.EnqueueJob(reportJob); err != nil {
			log.Fatalf("Failed to enqueue job: %v", err)
		}
	}

	select {}
}
