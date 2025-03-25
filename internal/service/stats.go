package service

import (
	"context"
	"fmt"
	"time"
)

type StatGetter interface {
	GetBannerClickStats(context.Context, int, time.Time, time.Time) (int, error)
}

type ClickStats struct {
	statGetter StatGetter
}

func NewClickStats(statGetter StatGetter) *ClickStats {
	return &ClickStats{statGetter: statGetter}
}

func (sg *ClickStats) GetClickStats(ctx context.Context, bannerID int, timeFrom, timeTo time.Time) (int, error) {
	result, err := sg.statGetter.GetBannerClickStats(ctx, bannerID, timeFrom, timeTo)
	if err != nil {
		return 0, fmt.Errorf("failed to get click stats: %w", err)
	}
	return result, nil
}
