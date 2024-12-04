package db

import (
	"context"
	"fmt"
	"time"

	"github.com/uptrace/bun"
)

type PrequalParametersResponse struct {
	bun.BaseModel `bun:"table:prequal_parameters_response"`

	Id                int       `bun:"id,pk,autoincrement" json:"id"`
	MaxLifeTime       int       `bun:"max_life_time" json:"max_life_time"`
	PoolSize          int       `bun:"pool_size" json:"pool_size"`
	ProbeFactor       float64   `bun:"probe_factor" json:"probe_factor"`
	ProbeRemoveFactor int   `bun:"probe_remove_factor" json:"probe_remove_factor"`
	Mu                int   `bun:"mu" json:"mu"`
	CreatedAt         time.Time `bun:"created_at,default:current_timestamp" json:"created_at"`
	UpdatedAt         time.Time `bun:"updated_at,default:current_timestamp" json:"updated_at"`
}

// Fetch latest created row
func GetPrequalParametersResponse(ctx context.Context) (PrequalParametersResponse, error) {
	var response PrequalParametersResponse
	err := db.NewSelect().
		Model(&response).
		Order("created_at DESC").
		Limit(1).
		Scan(ctx)
	return response, err
}

// Insert new row
func AddPrequalParametersResponse(ctx context.Context, response PrequalParametersResponse) error {
	_, err := db.NewInsert().Model(&response).Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to add prequal parameters response: %v", err)
	}

	message := fmt.Sprintf(
		"Added prequal parameters: MaxLifeTime=%d, PoolSize=%d, ProbeFactor=%.2f, ProbeRemoveFactor=%d, Mu=%d",
		response.MaxLifeTime, response.PoolSize, response.ProbeFactor, response.ProbeRemoveFactor, response.Mu,
	)
	if logErr := LogActivity(ctx, "success", message, nil); logErr != nil {
		return fmt.Errorf("failed to log activity for prequal parameters response: %v", logErr)
	}
	return err
}
