package store

import (
	"fmt"

	"github.com/sourcegraph/sourcegraph/internal/metrics"
	"github.com/sourcegraph/sourcegraph/internal/observation"
)

type operations struct {
	getStaleSourcedCommits *observation.Operation
	// Ranking
	vacuumStaleDefinitionsAndReferences       *observation.Operation
	vacuumStaleGraphs                         *observation.Operation
	insertDefinitionsAndReferencesForDocument *observation.Operation
	insertDefinitionsForRanking               *observation.Operation
	insertReferencesForRanking                *observation.Operation
	insertPathCountInputs                     *observation.Operation
	insertPathRanks                           *observation.Operation
}

var m = new(metrics.SingletonREDMetrics)

func newOperations(observationCtx *observation.Context) *operations {
	m := m.Get(func() *metrics.REDMetrics {
		return metrics.NewREDMetrics(
			observationCtx.Registerer,
			"codeintel_ranking_store",
			metrics.WithLabels("op"),
			metrics.WithCountHelp("Total number of method invocations."),
		)
	})

	op := func(name string) *observation.Operation {
		return observationCtx.Operation(observation.Op{
			Name:              fmt.Sprintf("codeintel.ranking.store.%s", name),
			MetricLabelValues: []string{name},
			Metrics:           m,
		})
	}

	return &operations{
		getStaleSourcedCommits:                    op("StaleSourcedCommits"),
		vacuumStaleDefinitionsAndReferences:       op("VacuumStaleDefinitionsAndReferences"),
		vacuumStaleGraphs:                         op("VacuumStaleGraphs"),
		insertDefinitionsAndReferencesForDocument: op("InsertDefinitionsAndReferencesForDocument"),
		insertDefinitionsForRanking:               op("InsertDefinitionsForRanking"),
		insertReferencesForRanking:                op("InsertReferencesForRanking"),
		insertPathCountInputs:                     op("InsertPathCountInputs"),
		insertPathRanks:                           op("InsertPathRanks"),
	}
}
