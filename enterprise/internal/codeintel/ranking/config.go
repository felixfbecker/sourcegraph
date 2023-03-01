package ranking

import (
	"time"

	"github.com/sourcegraph/sourcegraph/internal/env"
)

type rankingConfig struct {
	env.BaseConfig

	Interval time.Duration
}

var RankingConfigInst = &rankingConfig{}

func (c *rankingConfig) Load() {
	c.Interval = c.GetInterval("CODEINTEL_RANKING_RECALCULATION_INTERVAL", "72h", "The maximum age of document reference count values used for ranking before being considered stale.")
}

type exportConfig struct {
	env.BaseConfig

	RankingInterval    time.Duration
	NumRankingRoutines int
	RankingBatchSize   int
	RankingJobsEnabled bool
}

var ConfigExportInst = &exportConfig{}

func (c *exportConfig) Load() {
	c.RankingInterval = c.GetInterval("CODEINTEL_UPLOADS_RANKING_INTERVAL", "1s", "How frequently to serialize a batch of the code intel graph for ranking.")
	c.NumRankingRoutines = c.GetInt("CODEINTEL_UPLOADS_RANKING_NUM_ROUTINES", "4", "The number of concurrent ranking graph serializer routines to run per worker instance.")
	c.RankingBatchSize = c.GetInt("CODEINTEL_UPLOADS_RANKING_BATCH_SIZE", "10000", "The number of definitions and references to populate the ranking graph per batch.")
	c.RankingJobsEnabled = c.GetBool("CODEINTEL_UPLOADS_RANKING_JOB_ENABLED", "false", "Whether or not to run the ranking job.")
}
