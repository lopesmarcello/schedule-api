package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	dbName := os.Getenv("SCHED_DATABASE_NAME")
	dbUser := os.Getenv("SCHED_DATABASE_USER")

	fmt.Println("Trying to connect:", dbName, dbUser)

	cmd := exec.Command(
		"tern", "migrate",
		"--migrations", "./internal/repositories/pg/migrations",
		"--config", "./internal/repositories/pg/migrations/tern.conf",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("failed to migrate")
		fmt.Println("Output:\n", string(output))
		fmt.Println("Error:\n", err)
		panic(err)
	}

	fmt.Println("Executed successfully", string(output))
}
