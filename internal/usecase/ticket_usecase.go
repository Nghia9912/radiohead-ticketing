package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/NghiaHoang/radiohead-ticketing/internal/domain"
	"github.com/NghiaHoang/radiohead-ticketing/internal/repository/redis"
	"github.com/NghiaHoang/radiohead-ticketing/pkg/messagebroker"
)

type TicketUseCaseImpl struct {
	pgRepo        domain.TicketRepository
	redisRepo     *redis.TicketRedisRepo
	kafkaProducer *messagebroker.EventProducer
}

func NewTicketUseCase(pg domain.TicketRepository, rd *redis.TicketRedisRepo, kp *messagebroker.EventProducer) domain.TicketUsecase {
	return &TicketUseCaseImpl{pgRepo: pg, redisRepo: rd, kafkaProducer: kp}
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
// ConfirmPurchase handles the final checkout logic
func (u *TicketUseCaseImpl) ConfirmPurchase(ctx context.Context, userID, ticketID, orderID, email string) error {
	// 1. Fetch current status and version from DB
	// TODO: Implement GetTicket in pgRepo to fetch the real currentVersion
	var currentVersion int32 = 1 // Placeholder for now

	// 2. Execute Optimistic Locking to transition status to SOLD
	err := u.pgRepo.UpdateStatus(ctx, ticketID, domain.TicketLocked, domain.TicketSold, currentVersion)
	if err != nil {
		return err // Could be ErrOptimisticLockConflict
	}

	// 3. Publish Event to Kafka asynchronously
	event := messagebroker.TicketSoldEvent{
		OrderID:  orderID,
		UserID:   userID,
		TicketID: ticketID,
		Email:    email,
	}
	
	// Publish the event to Kafka. Do not block the main flow.
	err = u.kafkaProducer.PublishTicketSold("ticket_sold_events", event)
	if err != nil {
		// In production, we should log this error properly
	}

	return nil
}