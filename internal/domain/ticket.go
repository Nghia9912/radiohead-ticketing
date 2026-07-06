package domain

import (
	"context"
	"errors"
)

var (
	ErrOptimisticLockConflict = errors.New("optimistic lock conflict: ticket has been updated by another transaction")
	ErrTicketUnavailable      = errors.New("ticket is currently unavailable or locked by another user")
)

// TicketStatus manages the lifecycle of a ticket
type TicketStatus string

const (
	TicketAvailable TicketStatus = "AVAILABLE"
	TicketLocked    TicketStatus = "LOCKED" // Held in Redis
	TicketSold      TicketStatus = "SOLD"   // Successfully paid and stored in Postgres
)

// Ticket represents the atomic ticket entity
type Ticket struct {
	ID             string       `json:"id"`
	EventID        string       `json:"event_id"`
	TicketTypeID   string       `json:"ticket_type_id"`
	SeatIdentifier string       `json:"seat_identifier"`
	Status         TicketStatus `json:"status"`
	Version        int32        `json:"version"` // Critical: Used for Optimistic Locking
}

// TicketRepository defines the contract for communicating with the database layer
type TicketRepository interface {
	GetTicket(ctx context.Context, ticketID string) (*Ticket, error)
	GetAvailableTickets(ctx context.Context, eventID string, limit int) ([]*Ticket, error)
	// UpdateStatus requires the current version to prevent race conditions
	UpdateStatus(ctx context.Context, ticketID string, oldStatus, newStatus TicketStatus, currentVersion int32) error
}

// TicketUsecase defines the usecase interface
type TicketUsecase interface {
	HoldTicket(ctx context.Context, userID string, ticketID string) error
	ConfirmPurchase(ctx context.Context, userID, ticketID, orderID, email string) error
}