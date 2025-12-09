package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lopesmarcello/schedule-api/internal/repositories/pg"
)

type AppointmentsService struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
}

func NewAppointmentsService(pool *pgxpool.Pool) AppointmentsService {
	return AppointmentsService{
		pool:    pool,
		queries: pg.New(pool),
	}
}

type Slot struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

var (
	bufferMinutes   = 0
	durationMinutes = 60
	timeLayout      = "15:04:05"
)

func (as *AppointmentsService) GetAvailableSpots(ctx context.Context, userID int32, dateStr string) ([]Slot, string, error) {
	// Parse date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, "", fmt.Errorf("invalid date: %w", err)
	}
	if date.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, "date is in the past", fmt.Errorf("date is in the past. %v", date)
	}

	// Fetch availability
	availParams := pg.GetAvailabilityForDayParams{
		UserID:    pgtype.Int4{Int32: userID, Valid: true},
		DayOfWeek: int32(date.Weekday()),
	}

	availabilities, err := as.queries.GetAvailabilityForDay(ctx, availParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, "no availability found for day", nil
		}
		return nil, "error fetching availability", fmt.Errorf("error: %w", err)
	}

	if len(availabilities) == 0 {
		return nil, "no availability entries for day", nil
	}
	avail := availabilities[0]

	// Convert pgtype.Time to time.Time (microseconds since midnight)
	startTime := time.Time{}.Add(time.Duration(avail.StartTime.Microseconds) * time.Microsecond)
	endTime := time.Time{}.Add(time.Duration(avail.EndTime.Microseconds) * time.Microsecond)
	apptParams := pg.GetAppointmentsForDateParams{ // Assuming this exists
		UserID:          pgtype.Int4{Int32: userID, Valid: true},
		AppointmentDate: pgtype.Date{Time: date, Valid: true},
	}
	appts, err := as.queries.GetAppointmentsForDate(ctx, apptParams)
	if err != nil {
		return nil, "failed to fetch appointments", fmt.Errorf("failed to fetch appointments: %w", err)
	}

	// Generate possible slots
	slotDuration := time.Duration(durationMinutes) * time.Minute
	bufferDuration := time.Duration(bufferMinutes) * time.Minute
	current := startTime
	var possibleSlots []Slot
	log.Printf("startTime: %v, endTime: %v", startTime, endTime)
	for {
		slotEnd := current.Add(slotDuration)
		if slotEnd.After(endTime) {
			break // Stop if slot would overrun end_time
		}
		possibleSlots = append(possibleSlots, Slot{
			Start: current.Format("15:04:05"),
			End:   slotEnd.Format("15:04:05"),
		})
		current = slotEnd.Add(bufferDuration) // Advance by buffer for next slot
	}

	// Filter free slots
	var freeSlots []Slot
	for _, slot := range possibleSlots {
		slotStart, _ := time.Parse("15:04:05", slot.Start) // Ignore err for simplicity
		slotEnd, _ := time.Parse("15:04:05", slot.End)
		isFree := true
		for _, appt := range appts {
			apptStart := time.Time{}.Add(time.Duration(appt.StartTime.Microseconds) * time.Microsecond)
			apptEnd := time.Time{}.Add(time.Duration(appt.EndTime.Microseconds) * time.Microsecond)
			if slotEnd.After(apptStart) && slotStart.Before(apptEnd) {
				isFree = false
				break
			}
		}
		if isFree {
			freeSlots = append(freeSlots, slot)
		}
	}
	log.Printf("free slots: %v", freeSlots)

	return freeSlots, "", nil
}