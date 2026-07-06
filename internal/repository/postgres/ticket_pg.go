package postgres

import (
	"context"
	"database/sql"

	"github.com/NghiaHoang/radiohead-ticketing/internal/domain"
)

type TicketPGRepo struct {
	db *sql.DB
}

func NewTicketPGRepo(db *sql.DB) domain.TicketRepository {
	return &TicketPGRepo{db: db}
}

func (r *TicketPGRepo) UpdateStatus(ctx context.Context, ticketID string, oldStatus, newStatus domain.TicketStatus, version int32) error {
	query := `
		UPDATE tickets 
		SET status = $1, version = version + 1 
		WHERE id = $2 AND status = $3 AND version = $4
	`

	result, err := r.db.ExecContext(ctx, query, newStatus, ticketID, oldStatus, version)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If rowsAffected == 0, it means someone else has modified the status or version of this ticket.
	// This is the core mechanism of Optimistic Locking.
	if rowsAffected == 0 {
		return domain.ErrOptimisticLockConflict
	}

	return nil
}

func (r *TicketPGRepo) GetTicket(ctx context.Context, ticketID string) (*domain.Ticket, error) {
	query := `SELECT id, event_id, ticket_type_id, seat_identifier, status, version FROM tickets WHERE id = $1`
	var t domain.Ticket
	err := r.db.QueryRowContext(ctx, query, ticketID).Scan(
		&t.ID, &t.EventID, &t.TicketTypeID, &t.SeatIdentifier, &t.Status, &t.Version,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TicketPGRepo) GetAvailableTickets(ctx context.Context, eventID string, limit int) ([]*domain.Ticket, error) {
	query := `SELECT id, event_id, ticket_type_id, seat_identifier, status, version 
	          FROM tickets WHERE event_id = $1 AND status = $2 LIMIT $3`

	rows, err := r.db.QueryContext(ctx, query, eventID, domain.TicketAvailable, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []*domain.Ticket
	for rows.Next() {
		var t domain.Ticket
		if err := rows.Scan(&t.ID, &t.EventID, &t.TicketTypeID, &t.SeatIdentifier, &t.Status, &t.Version); err != nil {
			return nil, err
		}
		tickets = append(tickets, &t)
	}
	return tickets, nil
}
