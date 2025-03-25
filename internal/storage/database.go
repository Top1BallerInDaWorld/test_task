package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"
)

type ClicksStorage struct {
	db *pgxpool.Pool
}

func NewClickStorage(ctx context.Context, connectionString string) (*ClicksStorage, error) {
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "new database pool", err)
	}
	return &ClicksStorage{db: pool}, nil
}

func (cs *ClicksStorage) BulkInsertClicks(ctx context.Context, clicks map[int]int) error {
	query := `INSERT INTO clicks (banner_id, count) VALUES (@bannerID, @count)`

	batch := &pgx.Batch{}
	for k, v := range clicks {
		args := pgx.NamedArgs{
			"bannerID": k,
			"count":    v,
		}
		batch.Queue(query, args)
	}

	results := cs.db.SendBatch(ctx, batch)
	defer results.Close()

	for k, v := range clicks {
		_, err := results.Exec()
		if err != nil {
			slog.Error("error inserting clicks", "error", err,
				"bannerID", k, "lost clicks", v)
			continue
		}
	}

	return results.Close()
}

func (cs *ClicksStorage) GetBannerClickStats(ctx context.Context, bannerID int, tsFrom, tsTo time.Time) (int, error) {
	var clicksSummary int64

	err := cs.db.QueryRow(
		ctx,
		`SELECT COALESCE(SUM(count), 0)::bigint 
         FROM clicks
         WHERE banner_id = $1 
         AND timestamp BETWEEN $2 AND $3`,
		bannerID,
		tsFrom,
		tsTo,
	).Scan(&clicksSummary)
	if err != nil {
		return 0, fmt.Errorf("failed to get clicks stats: %w", err)
	}

	return int(clicksSummary), nil
}
