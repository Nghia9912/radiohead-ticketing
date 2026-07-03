package postgres

import (
	"context"
	"database/sql"
	"errors"

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