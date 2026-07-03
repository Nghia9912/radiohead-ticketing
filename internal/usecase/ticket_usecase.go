package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/NghiaHoang/radiohead-ticketing/internal/domain"
	"github.com/NghiaHoang/radiohead-ticketing/internal/repository/redis"
)

type TicketUseCaseImpl struct {
	pgRepo    domain.TicketRepository
	redisRepo *redis.TicketRedisRepo
}

func NewTicketUseCase(pg domain.TicketRepository, rd *redis.TicketRedisRepo) domain.TicketUsecase {
	return &TicketUseCaseImpl{pgRepo: pg, redisRepo: rd}
}

// HoldTicket processes the user request to hold/reserve a ticket
func (u *TicketUseCaseImpl) HoldTicket(ctx context.Context, userID string, ticketID string) error {
	// 1. Check if the ticket is AVAILABLE in cache/DB (Skip if the API only returns currently available tickets)
	// For optimization, we assume the client passed a valid ticketID and proceed to lock it on Redis.

	// 2. Lock the ticket in Redis with a 10-minute expiration.
	// This blocks all other requests to the same ticketID at the in-memory layer.
	err := u.redisRepo.LockTicket(ctx, ticketID, userID, 10*time.Minute)
	if err != nil {
		if errors.Is(err, redis.ErrTicketAlreadyLocked) {
			// Fail Fast: Return error immediately without querying the database
			return domain.ErrTicketUnavailable
		}
		return err
	}

	// Success: The ticket is reserved for the userID for 10 minutes.
	// The frontend can now start a countdown timer for checkout.
	return nil
}