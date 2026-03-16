package main

import (
	"github.com/spencerhum/duckmq/db"
	"github.com/spencerhum/duckmq/worker"
)

func main() {
	database := db.Connect()
	defer database.Close()

	worker.Start(database)
}
