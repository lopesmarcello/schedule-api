package services

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lopesmarcello/schedule-api/internal/repositories/pg"
)

type AvailabilityService struct {
	pool    *pgxpool.Pool
	queries *pg.Queries
}

func NewAvailabilityService(pool *pgxpool.Pool) AvailabilityService {
	return AvailabilityService{
		pool:    pool,
		queries: pg.New(pool),
	}
}

func (as *AvailabilityService) SetAvailability(ctx context.Context, userID, day int, start, end string) error {
	startTime, err := time.Parse("15:04:05", start)
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}
	startTimeInMS := int64(startTime.Hour()*3600000000 + startTime.Minute()*600000000 + startTime.Second()*1000000 + startTime.Nanosecond()/1000)

	startPg := pgtype.Time{
		Microseconds: startTimeInMS,
		Valid:        true,
	}

	endTime, err := time.Parse("15:04:05", start)
	if err != nil {
		return fmt.Errorf("invalid start time: %w", err)
	}
	endTimeInMS := int64(endTime.Hour()*3600000000 + endTime.Minute()*600000000 + endTime.Second()*1000000 + endTime.Nanosecond()/1000)

	endPg := pgtype.Time{
		Microseconds: endTimeInMS,
		Valid:        true,
	}

	userPgID := pgtype.Int4{
		Int32: int32(userID),
		Valid: true,
	}

	params := pg.SetAvailabilityParams{
		UserID:    userPgID,
		DayOfWeek: int32(day),
		StartTime: startPg,
		EndTime:   endPg,
	}

	err = as.queries.SetAvailability(ctx, params)
	if err != nil {
		return fmt.Errorf("error setting availability %w", err)
	}
	return nil
}
