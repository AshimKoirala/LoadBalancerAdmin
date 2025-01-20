package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/uptrace/bun"
)

type PrequalParametersResponse struct {
	bun.BaseModel `bun:"table:prequal_parameters_response"`

	Id                int       `bun:"id,pk,autoincrement" json:"id"`
	MaxLifeTime       int       `bun:"max_life_time" json:"max_life_time"`
	PoolSize          int       `bun:"pool_size" json:"pool_size"`
	ProbeFactor       float64   `bun:"probe_factor" json:"probe_factor"`
	ProbeRemoveFactor int       `bun:"probe_remove_factor" json:"probe_remove_factor"`
	Mu                int       `bun:"mu" json:"mu"`
	CreatedAt         time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
	Status            string    `bun:"status,default:inactive" json:"status"`
}

// Fetch latest created row
func GetPrequalParametersResponse(ctx context.Context) (PrequalParametersResponse, error) {
	var response PrequalParametersResponse
	err := db.NewSelect().
		Model(&response).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		log.Printf("Error fetching latest prequal parameters: %v", err)
		return response, err
	}
	// Check if the latest entry is active, if it's not fetch the last active entry
	if response.Status != "active" {
		var lastActive PrequalParametersResponse
		err := db.NewSelect().
			Model(&lastActive).
			Where("status = ?", "active").
			Order("created_at DESC").
			Limit(1).
			Scan(ctx)
		if err != nil {
			return response, fmt.Errorf("error fetching last active entry: %v", err)
		}
		// If the latest is not active, return the last active entry
		response = lastActive
	}

	return response, nil
}

type AddPrequalParametersType struct {
	MaxLifeTime       int     `json:"max_life_time"`
	PoolSize          int     `json:"pool_size"`
	ProbeFactor       float64 `json:"probe_factor"`
	ProbeRemoveFactor int     `json:"probe_remove_factor"`
	Mu                int     `json:"mu"`
	Status            string  `json:"status"`
}

// Insert new row
func AddPrequalParametersResponse(ctx context.Context, response AddPrequalParametersType) (*PrequalParametersResponse, error) {
	payload := &PrequalParametersResponse{
		MaxLifeTime:       response.MaxLifeTime,
		PoolSize:          response.PoolSize,
		ProbeFactor:       response.ProbeFactor,
		ProbeRemoveFactor: response.ProbeRemoveFactor,
		Mu:                response.Mu,
		Status:            "active",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Insert the new record
	_, err := db.NewInsert().Model(payload).Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to add prequal parameters response: %v", err)
	}

	if logErr := LogActivity(ctx, "success", "Prequal Parameters Updated", nil); logErr != nil {
		return nil, fmt.Errorf("failed to log activity for prequal parameters response: %v", logErr)
	}

	return payload, nil
}
