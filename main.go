package main

import (
	"fmt"

	"github.com/spencerhum/duckmq/db"
)

func main() {
	database := db.Connect()
	defer database.Close()

	fmt.Println("DuckMQ is running")
}
