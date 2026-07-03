package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrTicketAlreadyLocked = errors.New("ticket is currently locked by another user")

type TicketRedisRepo struct {
	client *redis.Client
}

func NewTicketRedisRepo(client *redis.Client) *TicketRedisRepo {
	return &TicketRedisRepo{client: client}
}

// reserveScript ensures Atomicity. 
// KEYS[1]: ticket_id, ARGV[1]: user_id, ARGV[2]: TTL (seconds)
const reserveScript = `
local current_holder = redis.call("GET", KEYS[1])
if current_holder == false then
    redis.call("SET", KEYS[1], ARGV[1], "EX", ARGV[2])
    return 1
elseif current_holder == ARGV[1] then
    redis.call("EXPIRE", KEYS[1], ARGV[2])
    return 1
else
    return 0
end
`

// LockTicket executes the Lua script to acquire a temporary hold/lock
func (r *TicketRedisRepo) LockTicket(ctx context.Context, ticketID string, userID string, duration time.Duration) error {
	key := "ticket:lock:" + ticketID
	ttl := int(duration.Seconds())

	// Run the script on the Redis server
	result, err := r.client.Eval(ctx, reserveScript, []string{key}, userID, ttl).Result()
	if err != nil {
		return err
	}

	if result.(int64) == 0 {
		return ErrTicketAlreadyLocked
	}

	return nil
}