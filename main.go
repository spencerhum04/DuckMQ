package main

import (
	"fmt"
	"log"

	"github.com/spencerhum/duckmq/db"
	"github.com/spencerhum/duckmq/queue"
)

func main() {
	database := db.Connect()
	defer database.Close()

	worker.Start(database)
}
