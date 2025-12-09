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

	os.Setenv("SCHED_DATABASE_PORT", os.Getenv("SCHED_TEST_DATABASE_PORT"))
	os.Setenv("SCHED_DATABASE_NAME", os.Getenv("SCHED_TEST_DATABASE_NAME"))
	os.Setenv("SCHED_DATABASE_PASSWORD", os.Getenv("SCHED_TEST_DATABASE_PASSWORD"))
	os.Setenv("SCHED_DATABASE_USER", os.Getenv("SCHED_TEST_DATABASE_USER"))
	os.Setenv("SCHED_DATABASE_HOST", os.Getenv("SCHED_TEST_DATABASE_HOST"))

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
