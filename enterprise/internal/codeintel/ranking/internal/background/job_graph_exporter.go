package background

import (
	"context"
	"time"

	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

func NewRankingGraphExporter(
	observationCtx *observation.Context,
	rankingService RankingService,
	numRankingRoutines int,
	interval time.Duration,
	batchSize int,
	rankingJobEnabled bool,
) goroutine.BackgroundRoutine {
	return goroutine.NewPeriodicGoroutine(
		context.Background(),
		"rank.graph-exporter", "exports SCIP data to ranking definitions and reference tables",
		interval,
		goroutine.HandlerFunc(func(ctx context.Context) error {
			if err := rankingService.ExportRankingGraph(ctx, numRankingRoutines, batchSize, rankingJobEnabled); err != nil {
				return err
			}

			if err := rankingService.VacuumRankingGraph(ctx); err != nil {
				return err
			}

			return nil
		}),
	)
}

func NewRankingGraphMapper(
	observationCtx *observation.Context,
	rankingService RankingService,
	numRankingRoutines int,
	interval time.Duration,
	rankingJobEnabled bool,
) goroutine.BackgroundRoutine {
	operations := newRankMappingOperations(observationCtx)
	return goroutine.NewPeriodicGoroutine(
		context.Background(),
		"rank.graph-mapper", "maps definitions and references data to path_counts_inputs table in store",
		interval,
		goroutine.HandlerFunc(func(ctx context.Context) error {
			numReferenceRecordsProcessed, numInputsInserted, err := rankingService.MapRankingGraph(ctx, numRankingRoutines, rankingJobEnabled)
			if err != nil {
				return err
			}

			operations.numReferenceRecordsProcessed.Add(float64(numReferenceRecordsProcessed))
			operations.numInputsInserted.Add(float64(numInputsInserted))
			return nil
		}),
	)
}

func NewRankingGraphReducer(
	observationCtx *observation.Context,
	rankingService RankingService,
	numRankingRoutines int,
	interval time.Duration,
	rankingJobEnabled bool,
) goroutine.BackgroundRoutine {
	operations := newRankReducingOperations(observationCtx)
	return goroutine.NewPeriodicGoroutine(
		context.Background(),
		"rank.graph-reducer", "reduces path_counts_inputs into a count of paths per repository and stores it in path_ranks table in store.",
		interval,
		goroutine.HandlerFunc(func(ctx context.Context) error {
			numPathRanksInserted, numPathCountsInputsProcessed, err := rankingService.ReduceRankingGraph(ctx, numRankingRoutines, rankingJobEnabled)
			if err != nil {
				return err
			}

			operations.numPathCountsInputsRowsProcessed.Add(numPathCountsInputsProcessed)
			operations.numPathRanksInserted.Add(numPathRanksInserted)
			return nil
		}),
	)
}
