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
		Status:              ACTIVE,
		HealthCheckEndpoint: healthCheckEndpoint,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	_, err := db.NewInsert().Model(replica).Returning("id").Exec(ctx)
	if err != nil {
		return fmt.Errorf("error adding replica: %v", err)
	}

	if err := LogActivity(ctx, "success", fmt.Sprintf("Replica '%s' added successfully", name), &replica.ID); err != nil {
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
        replica, err = GetReplicaByID(ctx, *id)
    } else if url != nil {
        log.Printf("Attempting to remove replica with URL: %s", *url)
        replica, err = GetReplicaByURL(ctx, *url)
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
	log.Printf("Successfully disabled replica with ID: %d", replica.ID)

    // Delete the replica from the database
    // _, err = db.NewDelete().Model((*Replica)(nil)).Where("id = ?", id).Exec(ctx)
    // if err != nil {
    //     log.Printf("Error deleting replica with ID: %d, error: %v", id, err)
    //     return fmt.Errorf("error removing replica: %v", err)
    // }

    // log.Printf("Successfully deleted replica with ID: %d", id)

    // Log the activity
    if err := LogActivity(ctx, "warning", fmt.Sprintf("Replica '%s' removed", replica.Name), &replica.ID);
	 err != nil {
    log.Printf("Error logging activity for replica '%s': %v", replica.Name, err)
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

	if err := LogActivity(ctx, "success", fmt.Sprintf("Replica '%s' status changed to '%s'", replica.Name, newStatus), &id); err != nil {
	return fmt.Errorf("error logging activity: %v", err)
     }

	return nil
}

// find a replica by id
func GetReplicaByID(ctx context.Context, id int64) (*Replica, error) {
	var replica Replica
	err := db.NewSelect().Model(&replica).Where("id = ?", id).Scan(ctx)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return &replica, nil
}

func GetReplicaByURL(ctx context.Context, url string) (*Replica, error) {
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