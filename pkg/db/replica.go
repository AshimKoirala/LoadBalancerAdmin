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
	return nil
}

// remove replica by id
func RemoveReplica(ctx context.Context, id int64) error {
	switch {
	case id <= 0:
		return fmt.Errorf("valid replica ID must be provided")
	default:
		_, err := db.NewDelete().Model((*Replica)(nil)).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return fmt.Errorf("error removing replica: %v", err)
		}
		fmt.Println("Replica removed successfully.")
		return nil
	}
}

// Change status of a replica to active, inactive, or disabled
func ChangeStatus(ctx context.Context, id int64, newStatus string) error {
	if newStatus != ACTIVE && newStatus != INACTIVE && newStatus != DISABLED {
		return fmt.Errorf("invalid status: %s. Allowed values are 'active', 'inactive', or 'disabled'", newStatus)
	}

	_, err := db.NewUpdate().
		Model((*Replica)(nil)).
		Set("status = ?", newStatus).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error changing replica status: %v", err)
	}

	fmt.Printf("Replica ID %d status changed to %s successfully.\n", id, newStatus)
	return nil
}

// find a replica by id
func GetReplicaByID(id int64) (Replica, error) {
	var replica Replica
	err := db.NewSelect().Model(&replica).Where("id = ?", id).Scan(ctx)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return replica, err
	}
	return replica, nil
}
