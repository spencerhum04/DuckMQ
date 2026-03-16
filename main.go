package main

import (
	"github.com/spencerhum/duckmq/db"
	"github.com/spencerhum/duckmq/worker"
)

func main() {
	database := db.Connect()
	defer database.Close()

	job, err := queue.Enqueue(database, "send email", map[string]any{
		"to":      "test@example.com",
		"subject": "test message from DuckMQ",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Enqueued job id=%d type=%s:", job.ID, job.Type)
}
