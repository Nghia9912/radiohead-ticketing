package domain

import (
	"context"
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
	GetAvailableTickets(ctx context.Context, eventID string, limit int) ([]*Ticket, error)
	// UpdateStatus requires the current version to prevent race conditions
	UpdateStatus(ctx context.Context, ticketID string, oldStatus, newStatus TicketStatus, currentVersion int32) error
}