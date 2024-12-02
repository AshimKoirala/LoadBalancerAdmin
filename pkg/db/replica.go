package db

import (
	"context"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Replica struct {
	bun.BaseModel `bun:"table:replicas"`

	ID                  int64     `json:"id" bun:"id,pk,autoincrement"`
	Name                string    `json:"name" bun:"name,unique,notnull"`
	URL                 string    `json:"url" bun:"url,unique,notnull"`
	Status              string    `bun:"status,notnull"`
	HealthCheckEndpoint string    `bun:"healthcheck_endpoint,notnull"`
	CreatedAt           time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt           time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
}

const (
	INACTIVE = "inactive"
	ACTIVE   = "active"
	DISABLED = "disabled"
)

// add replica
func AddReplica(ctx context.Context, name, url, healthCheckEndpoint string) error {
	replica := &Replica{
		Name:                name,
		URL:                 url,
		Status:              INACTIVE,
		HealthCheckEndpoint: healthCheckEndpoint,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	_, err := db.NewInsert().Model(replica).Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding replica: %v", err)
	}
	log := ActivityLog{
		Type:      "success",
		Message:   fmt.Sprintf("Replica '%s' added successfully", name),
		ReplicaID: &replica.ID,
	}
	if logErr := AddActivityLog(ctx, log); logErr != nil {
		return fmt.Errorf("error logging activity: %v", logErr)
	}

	return nil
}

// remove replica by id
func RemoveReplica(ctx context.Context, id int64) error {
    log.Printf("Attempting to remove replica with ID: %d", id)

    replica, err := GetReplicaByID(ctx, id)
    if err != nil {
        log.Printf("Error fetching replica with ID: %d, error: %v", id, err)
        return fmt.Errorf("error fetching replica: %v", err)
    }

    // Delete the replica from the database
    _, err = db.NewDelete().Model((*Replica)(nil)).Where("id = ?", id).Exec(ctx)
    if err != nil {
        log.Printf("Error deleting replica with ID: %d, error: %v", id, err)
        return fmt.Errorf("error removing replica: %v", err)
    }

    log.Printf("Successfully deleted replica with ID: %d", id)

    // Log the activity even if it fails
    logEntry := ActivityLog{
        Type:      "warning",
        Message:   fmt.Sprintf("Replica '%s' removed", replica.Name),
        ReplicaID: &id,
    }

    if logErr := AddActivityLog(ctx, logEntry); logErr != nil {
        log.Printf("Error logging activity for replica '%s': %v", replica.Name, logErr)
    }

    return nil
}


// Change status of a replica to active, inactive, or disabled
func ChangeStatus(ctx context.Context, id int64, newStatus string) error {

	replica, err := GetReplicaByID(ctx ,id)
	if err != nil {
		return fmt.Errorf("error fetching replica: %v", err)
	}

	if newStatus != ACTIVE && newStatus != INACTIVE && newStatus != DISABLED {
		return fmt.Errorf("invalid status: %s. Allowed values are 'active', 'inactive', or 'disabled'", newStatus)
	}

	_, err = db.NewUpdate().
		Model((*Replica)(nil)).
		Set("status = ?", newStatus).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error changing replica status: %v", err)
	}

	log := ActivityLog{
		Type:      "success",
		Message:   fmt.Sprintf("Replica '%s' status changed to '%s'", replica.Name, newStatus),
		ReplicaID: &id,
	}
	if logErr := AddActivityLog(ctx, log); logErr != nil {
		return fmt.Errorf("error logging activity: %v", logErr)
	}

	return nil
}

// find a replica by id
func GetReplicaByID(ctx context.Context, id int64) (Replica, error) {
	var replica Replica
	err := db.NewSelect().Model(&replica).Where("id = ?", id).Scan(ctx)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return replica, err
	}
	return replica, nil
}
func GetReplicas(ctx context.Context) ([]Replica, error) {
	var replicas []Replica
	err := db.NewSelect().Model(&replicas).Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching replicas: %v", err)
	}
	return replicas, nil
}