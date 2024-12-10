package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type ActivityLog struct {
	bun.BaseModel `bun:"table:activity_logs"`

	Id        int64     `bun:"id,pk,autoincrement"`
	Type      string    `bun:"type,notnull"`
	Message   string    `bun:"message,notnull"`
	ReplicaId *int64    `bun:"replica_id"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp"`

	Replica *Replica `bun:"rel:belongs-to,join:replica_id=id"`
}

// adds new activity log entry.
func AddActivityLog(ctx context.Context, log ActivityLog) error {
	_, err := db.NewInsert().Model(&log).Exec(ctx)
	return err
}

// retrieves all activity logs in descending order
func FetchActivityLogs(ctx context.Context) ([]ActivityLog, error) {
	var logs []ActivityLog
	err := db.NewSelect().
		Model(&logs).
		Relation("Replica"). // Fetch associated replica details
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching activity logs: %v", err)
	}
	return logs, nil
}

// reusable function to log activity
func LogActivity(ctx context.Context, activityType, message string, replicaId *int64) error {
	log := ActivityLog{
		Type:      activityType,
		Message:   message,
		ReplicaId: replicaId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return AddActivityLog(ctx, log)
}
