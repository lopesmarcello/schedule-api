package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lopesmarcello/schedule-api/internal/api"
	"github.com/lopesmarcello/schedule-api/internal/services"
)

func init() {
	godotenv.Load()
}

func main() {
	fmt.Println("hello from api!")

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable search_path=public",
		os.Getenv("SCHED_DATABASE_USER"),
		os.Getenv("SCHED_DATABASE_PASSWORD"),
		os.Getenv("SCHED_DATABASE_HOST"),
		os.Getenv("SCHED_DATABASE_PORT"),
		os.Getenv("SCHED_DATABASE_NAME"),
	),
	)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("Unable to ping database %v", err))
	}

	userService := services.NewUserService(pool)
	availabilityService := services.NewAvailabilityService(pool)
	appointmentsService := services.NewAppointmentsService(pool)

	api := api.NewAPI(*userService, *availabilityService, *appointmentsService)

	// Start server on 8080
	// 0.0.0.0:8080 || localhost:8080
	if err := api.Router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
