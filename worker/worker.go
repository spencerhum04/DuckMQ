package worker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Job struct {
	ID      int64
	Type    string
	Payload map[string]any
}

func Start(db *sql.DB) {
	fmt.Println("worker started, polling for jobs...")

	for {
		job, err := dequeue(db)
		if err != nil {
			log.Println("error dequeuing:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if job == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		fmt.Printf("picked up job id=%d type=%s\n", job.ID, job.Type)
		processJob(db, job)
	}
}

func dequeue(db *sql.DB) (*Job, error) {
	row := db.QueryRow(`
		UPDATE jobs
		SET status = 'running', updated_at = NOW()
		WHERE id = (
			SELECT id FROM jobs
			WHERE status = 'pending'
			AND run_at <= NOW()
			ORDER BY run_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, type, payload
	`)

	var job Job
	var rawPayload []byte

	err := row.Scan(&job.ID, &job.Type, &rawPayload)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("dequeue failed: %w", err)
	}

	if err := json.Unmarshal(rawPayload, &job.Payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &job, nil
}

func processJob(db *sql.DB, job *Job) {
	fmt.Printf("processing job id=%d payload=%v\n", job.ID, job.Payload)
	time.Sleep(500 * time.Millisecond)

	_, err := db.Exec(`
		UPDATE jobs
		SET status = 'done', updated_at = NOW()
		WHERE id = $1
	`, job.ID)

	if err != nil {
		log.Printf("failed to mark job %d done: %v\n", job.ID, err)
		return
	}

	fmt.Printf("job id=%d done\n", job.ID)
}
