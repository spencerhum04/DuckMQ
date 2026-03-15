package queue

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Job struct {
	ID      int64
	Type    string
	Payload map[string]any
}

func Enqueue(db *sql.DB, jobType string, payload map[string]any) (*Job, error) {
	// Marshalling basically converts the map into JSON
	payloadBytes, err := json.Marshal(payload)

	if err != nil {
		return nil, fmt.Errorf("Failed to marshal payload: %w", err)
	}

	row := db.QueryRow(`
		INSERT INTO jobs (type, payload, status)
		VALUES ($1, $2, 'pending')
		RETURNING id, type, payload
	`, jobType, payloadBytes)

	var job Job
	var rawPayload []byte

	err = row.Scan(&job.ID, &job.Type, &rawPayload)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert job: %w", err)
	}

	if err := json.Unmarshal(rawPayload, &job.Payload); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal payload: %w", err)
	}

	return &job, nil
}
