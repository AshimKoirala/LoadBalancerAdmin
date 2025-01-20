package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
)

type Replica struct {
	bun.BaseModel `bun:"table:replicas"`

	Id                  int64     `json:"id" bun:"id,pk,autoincrement"`
	Name                string    `json:"name" bun:"name,unique,notnull"`
	URL                 string    `json:"url" bun:"url,unique,notnull"`
	Status              string    `json:"status" bun:"status,notnull"`
	HealthCheckEndpoint string    `json:"health_check_point" bun:"health_check_endpoint,notnull"`
	CreatedAt           time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt           time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
}

const (
	INACTIVE = "inactive"
	ACTIVE   = "active"
	DISABLED = "disabled"
)

func AddReplica(ctx context.Context, name, url, healthCheckEndpoint string) error {
	replica := &Replica{
		Name:                name,
		URL:                 url,
		Status:              INACTIVE,
		HealthCheckEndpoint: healthCheckEndpoint,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	var findReplica Replica
	err := db.NewSelect().
		Model(&findReplica).
		Where("url = ?", url).
		Where("name = ?", name).
		Scan(ctx)

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error occurred: %v", err)
	}

	if err == nil {
		// Update the replica's status to active
		_, updateErr := db.NewUpdate().
			Model(&findReplica).
			Set("status = ?", ACTIVE).
			Set("updated_at = ?", time.Now()).
			Set("health_check_endpoint = ?", replica.HealthCheckEndpoint).
			Where("url = ?", url).
			Exec(ctx)
		if updateErr != nil {
			return fmt.Errorf("error updating replica: %v", updateErr)
		}
		replica = &findReplica
	} else {
		// Insert new replica
		_, err = db.NewInsert().Model(replica).Exec(ctx)
		if err != nil {
			return fmt.Errorf("error adding replica: %v", err)
		}
	}

	// Log activity
	if err = LogActivity(ctx, "success", fmt.Sprintf("Replica '%s' is ready to be active", name), &replica.Id); err != nil {
		return fmt.Errorf("error logging activity: %v", err)
	}

	return nil
}

// remove replica by id
func RemoveReplica(ctx context.Context, id *int64, url *string) error {
	var replica *Replica
	var err error

	if id != nil {
		log.Printf("Attempting to remove replica with ID: %d", *id)
		replica, err = GetReplicaById(ctx, *id)
	} else if url != nil {
		log.Printf("Attempting to remove replica with URL: %s", *url)
		replica, err = GetReplicaByUrl(ctx, *url)
	} else {
		return fmt.Errorf("either id or url must be provided")
	}

	if err != nil {
		log.Printf("Error fetching replica: %v", err)
		return fmt.Errorf("error fetching replica: %v", err)
	}

	// Set status to disabled
	query := db.NewUpdate().
		Model((*Replica)(nil)).
		Set("status = ?", "disabled")

	if id != nil {
		query.Where("id = ?", *id)
	} else if url != nil {
		query.Where("url = ?", *url)
	}

	_, err = query.Exec(ctx)
	if err != nil {
		return fmt.Errorf("error disabling replica: %v", err)
	}
	log.Printf("Successfully disabled replica with ID: %d", replica.Id)

	// Delete the replica from the database
	// _, err = db.NewDelete().Model((*Replica)(nil)).Where("id = ?", id).Exec(ctx)
	// if err != nil {
	//     log.Printf("Error deleting replica with ID: %d, error: %v", id, err)
	//     return fmt.Errorf("error removing replica: %v", err)
	// }

	// log.Printf("Successfully deleted replica with ID: %d", id)

	// Log the activity
	if err := LogActivity(ctx, "warning", fmt.Sprintf("Replica '%s' is disabled", replica.Name), &replica.Id); err != nil {
		log.Printf("Error logging activity for replica '%s': %v", replica.Name, err)
	}

	return nil
}

// Change status of a replica to active, inactive, or disabled
func UpdateStatus(ctx context.Context, id int64, newStatus string) error {
	replica, err := GetReplicaById(ctx, id)
	if err != nil {
		return fmt.Errorf("error fetching replica: %v", err)
	}

	if newStatus != ACTIVE && newStatus != INACTIVE && newStatus != DISABLED {
		return fmt.Errorf("invalid status: %s. Allowed values are 'active', 'inactive', or 'disabled'", newStatus)
	}

	// Check if the status has changed and set the new status
	oldStatus := replica.Status

	_, err = db.NewUpdate().
		Model((*Replica)(nil)).
		Set("status = ?", newStatus).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("error changing replica status: %v", err)
	}

	// Log the status change
	if oldStatus != newStatus {

		// Deactivated status
		if oldStatus == ACTIVE && newStatus == DISABLED {
			if err := LogActivity(ctx, "warning", fmt.Sprintf("Replica '%s' is deactivated", replica.Name), &id); err != nil {
				return fmt.Errorf("error logging activity: %v", err)
			}
		}

		// Activated status
		if oldStatus == DISABLED && newStatus == ACTIVE {
			if err := LogActivity(ctx, "success", fmt.Sprintf("Replica '%s' is queued for activation", replica.Name), &id); err != nil {
				return fmt.Errorf("error logging activity: %v", err)
			}
		}

		if oldStatus == INACTIVE && newStatus == ACTIVE {
			if err := LogActivity(ctx, "success", fmt.Sprintf("Replica '%s' is activated and running", replica.Name), &id); err != nil {
				return fmt.Errorf("error logging activity: %v", err)
			}
		}
	}

	return nil
}

// find a replica by id
func GetReplicaById(ctx context.Context, id int64) (*Replica, error) {
	var replica Replica
	err := db.NewSelect().Model(&replica).Where("id = ?", id).Scan(ctx)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return &replica, nil
}

func GetReplicaByUrl(ctx context.Context, url string) (*Replica, error) {
	replica := new(Replica)
	err := db.NewSelect().
		Model(replica).
		Where("url = ?", url).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching replica by URL: %v", err)
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

func GetReplicaByName(ctx context.Context, name string) (*Replica, error) {
	var replica Replica
	err := db.NewSelect().
		Model(&replica).
		Where("name = ?", name).
		Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch replica by name: %v", err)
	}
	return &replica, nil
}

func UpdateStatusByUrl(url string, newStatus string) error {
	var ctx = context.Background()
	_, err := db.NewUpdate().
		Model((*Replica)(nil)).
		Set("status = ?", newStatus).
		Set("updated_at = ?", time.Now()).
		Where("url = ?", url).
		Returning("id").
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("error updating replica status: %v", err)
	}

	return nil
}
