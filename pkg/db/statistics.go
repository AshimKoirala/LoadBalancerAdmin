package db

import (
	"context"
	"log"
	"time"

	"github.com/uptrace/bun"
)

type Statistics struct {
	bun.BaseModel `bun:"table:statistics"`

	Id                 int64     `json:"id" bun:"id,pk,autoincrement"`
	URL                string    `json:"url" bun:"url,unique,notnull"`
	SuccessfulRequests int64     `json:"successful_requests" bun:"successful_requests,default:0"`
	FailedRequests     int64     `json:"failed_requests" bun:"failed_requests,default:0"`
	CreatedAt          time.Time `json:"created_at" bun:"created_at,default:current_timestamp"`
	UpdatedAt          time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp"`
}

type StatisticsData struct {
	URL                string
	SuccessfulRequests int64
	FailedRequests     int64
}

func BatchAddStatistics(statistics *[]StatisticsData) error {
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Fatalf("Failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	for _, statistic := range *statistics {
		stat := Statistics{
			URL:                statistic.URL,
			SuccessfulRequests: statistic.SuccessfulRequests,
			FailedRequests:     statistic.FailedRequests,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		_, err := db.NewInsert().Model(&stat).On("CONFLICT (url) DO UPDATE").
			Set("successful_requests = statistics.successful_requests + ?", statistic.SuccessfulRequests).
			Set("failed_requests = statistics.failed_requests + ?", statistic.FailedRequests).
			Set("updated_at = NOW()").
			Exec(ctx)

		if err != nil {
			log.Print("Error inserting statistics:", err)
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
		return err
	}

	log.Printf("Statistics updated/inserted successfully")
	return nil
}
