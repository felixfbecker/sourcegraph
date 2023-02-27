package background

import (
	"context"
	"time"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/insights/background/retention"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/insights/store"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/timeutil"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func newRetentionEnqueuer(ctx context.Context, workerBaseStore *basestore.Store, insightStore store.DataSeriesStore) goroutine.BackgroundRoutine {
	return goroutine.NewPeriodicGoroutine(ctx,
		"insights.retention.enqueuer",
		"enqueues series retention jobs",
		2*time.Hour,
		goroutine.HandlerFunc(
			func(ctx context.Context) error {
				seriesArgs := store.GetDataSeriesArgs{ExcludeJustInTime: true, NextSnapshotBefore: timeutil.Now()}
				allSeries, err := insightStore.GetDataSeries(ctx, seriesArgs)
				if err != nil {
					return errors.Wrap(err, "unable to fetch series for retention")
				}
				var multi error
				for _, series := range allSeries {
					_, err = retention.EnqueueJob(ctx, workerBaseStore, &retention.DataRetentionJob{InsightSeriesID: series.ID, SeriesID: series.SeriesID})
					if err != nil {
						multi = errors.Append(multi, errors.Wrapf(err, "seriesID: %d", series.ID))
					}
				}
				return multi
			}))
}
