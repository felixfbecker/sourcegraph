package background

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/sourcegraph/sourcegraph/internal/observation"
)

type rankMappingOperations struct {
	numReferenceRecordsProcessed prometheus.Counter
	numInputsInserted            prometheus.Counter
}

func newRankMappingOperations(observationCtx *observation.Context) *rankMappingOperations {
	counter := func(name, help string) prometheus.Counter {
		counter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: name,
			Help: help,
		})

		observationCtx.Registerer.MustRegister(counter)
		return counter
	}

	numReferenceRecordsProcessed := counter(
		"src_codeintel_ranking_reference_records_processed_total",
		"The number of reference rows processed.",
	)
	numInputsInserted := counter(
		"src_codeintel_ranking_inputs_inserted_total",
		"The number of input rows inserted.",
	)

	return &rankMappingOperations{
		numReferenceRecordsProcessed: numReferenceRecordsProcessed,
		numInputsInserted:            numInputsInserted,
	}
}

type rankReducingOperations struct {
	numPathCountsInputsRowsProcessed prometheus.Counter
	numPathRanksInserted             prometheus.Counter
}

func newRankReducingOperations(observationCtx *observation.Context) *rankReducingOperations {
	counter := func(name, help string) prometheus.Counter {
		counter := prometheus.NewCounter(prometheus.CounterOpts{
			Name: name,
			Help: help,
		})

		observationCtx.Registerer.MustRegister(counter)
		return counter
	}

	numPathCountInputsRowsProcessed := counter(
		"src_codeintel_ranking_path_count_inputs_rows_processed_total",
		"The number of input rows processed.",
	)
	numPathRanksInserted := counter(
		"src_codeintel_ranking_path_ranks_inserted_total",
		"The number of path ranks inserted.",
	)

	return &rankReducingOperations{
		numPathCountsInputsRowsProcessed: numPathCountInputsRowsProcessed,
		numPathRanksInserted:             numPathRanksInserted,
	}
}
