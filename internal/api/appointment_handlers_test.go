package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/lopesmarcello/schedule-api/internal/repositories/pg"
	"github.com/lopesmarcello/schedule-api/internal/services"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(os.Stderr)
}

func loadEnv() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
func newTestDb() *pgxpool.Pool {
	loadEnv()
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable search_path=public",
		os.Getenv("SCHED_TEST_DATABASE_USER"),
		os.Getenv("SCHED_TEST_DATABASE_PASSWORD"),
		os.Getenv("SCHED_TEST_DATABASE_HOST"),
		os.Getenv("SCHED_TEST_DATABASE_PORT"),
		os.Getenv("SCHED_TEST_DATABASE_NAME"),
	),
	)
	if err != nil {
		panic(err)
	}

	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Sprintf("Unable to ping database %v", err))
	}
	return pool
}

func setupTestAPI() *gin.Engine {
	r := gin.Default()
	dbpool := newTestDb()

	api := &API{
		Router:              r,
		UserService:         services.NewUserService(dbpool),
		AvailabilityService: services.NewAvailabilityService(dbpool),
		AppointmentsService: services.NewAppointmentsService(dbpool),
	}

	api.BindRoutes()

	return r
}
func createTestUser(dbpool *pgxpool.Pool) {
	queries := pg.New(dbpool)
	ctx := context.Background()

	// Clean up previous test data
	if err := queries.DeleteAllAppointments(ctx); err != nil {
		log.Fatalf("failed to delete all appointments: %v", err)
	}

	if err := queries.DeleteAllAvailabilities(ctx); err != nil {
		log.Fatalf("failed to delete all availabilities: %v", err)
	}
	if err := queries.DeleteAllUsers(ctx); err != nil {
		log.Fatalf("failed to delete all users: %v", err)
	}

	user, err := queries.CreateUser(ctx, pg.CreateUserParams{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "password",
		Slug:         "test-user",
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	// January 1, 2026 is a Thursday, so day of week is 4
	_, err = queries.CreateAvailability(ctx, pg.CreateAvailabilityParams{
		UserID:    pgtype.Int4{Int32: user.ID, Valid: true},
		DayOfWeek: 4,
		StartTime: pgtype.Time{Microseconds: 9 * 60 * 60 * 1000000, Valid: true},
		EndTime:   pgtype.Time{Microseconds: 17 * 60 * 60 * 1000000, Valid: true},
	})

	if err != nil {
		log.Fatalf("failed to create availability: %v", err)
	}
}

func TestGetSlots(t *testing.T) {
	r := setupTestAPI()
	dbpool := newTestDb()
	createTestUser(dbpool)

	req, _ := http.NewRequest("GET", "/api/v1/test-user/slots/2026-01-01", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{
		"success": true,
		"data": [
			"09:00",
			"10:00",
			"11:00",
			"12:00",
			"13:00",
			"14:00",
			"15:00",
			"16:00"
		]
	}`, w.Body.String())
}

func TestGetSlotsNoAvailability(t *testing.T) {
	r := setupTestAPI()
	dbpool := newTestDb()
	queries := pg.New(dbpool)
	ctx := context.Background()

	// Clean up previous test data
	if err := queries.DeleteAllAppointments(ctx); err != nil {
		log.Fatalf("failed to delete all appointments: %v", err)
	}
	if err := queries.DeleteAllAvailabilities(ctx); err != nil {
		log.Fatalf("failed to delete all availabilities: %v", err)
	}
	if err := queries.DeleteAllUsers(ctx); err != nil {
		log.Fatalf("failed to delete all users: %v", err)
	}

	_, err := queries.CreateUser(ctx, pg.CreateUserParams{
		Name:         "User No Availability",
		Email:        "noavail@example.com",
		PasswordHash: "password",
		Slug:         "user-no-avail",
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	req, _ := http.NewRequest("GET", "/api/v1/user-no-avail/slots/2026-01-01", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"success": true, "data": null}`, w.Body.String())
}