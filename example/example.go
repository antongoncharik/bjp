package example

import (
	"fmt"
	"time"

	"github.com/antongoncharik/bjp/jobprocessor"
)

func SendEmailHandler(job *jobprocessor.Job) error {
	fmt.Printf("Sending email to: %s\n", job.Data)
	time.Sleep(2 * time.Second)
	fmt.Println("Email sent!")
	return nil
}

func GenerateReportHandler(job *jobprocessor.Job) error {
	fmt.Printf("Generating report for: %s\n", job.Data)
	time.Sleep(3 * time.Second)
	fmt.Println("Report generated!")
	return nil
}
