package ranking

import (
	"context"

	"cloud.google.com/go/storage"
	"github.com/sourcegraph/log"
	"google.golang.org/api/option"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/ranking/internal/background"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/ranking/internal/lsifstore"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/ranking/internal/store"
	codeintelshared "github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/shared"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/uploads"
	"github.com/sourcegraph/sourcegraph/internal/conf"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/env"
	"github.com/sourcegraph/sourcegraph/internal/goroutine"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/internal/symbols"
)

func NewService(
	observationCtx *observation.Context,
	db database.DB,
	codeIntelDB codeintelshared.CodeIntelDB,
	uploadSvc *uploads.Service,
	gitserverClient GitserverClient,
) *Service {
	if resultsGraphKey == "" {
		// The codenav default
		resultsGraphKey = "dev"
	}

	resultsBucket := func() *storage.BucketHandle {
		if resultsBucketCredentialsFile == "" {
			return nil
		}

		var opts []option.ClientOption
		if resultsBucketCredentialsFile != "" {
			opts = append(opts, option.WithCredentialsFile(resultsBucketCredentialsFile))
		}

		client, err := storage.NewClient(context.Background(), opts...)
		if err != nil {
			log.Scoped("ranking", "").Error("failed to create storage client", log.Error(err))
			return nil
		}

		return client.Bucket(resultsBucketName)
	}()

	return newService(
		scopedContext("service", observationCtx),
		store.New(scopedContext("store", observationCtx), db),
		lsifstore.New(scopedContext("lsifstore", observationCtx), codeIntelDB),
		uploadSvc,
		gitserverClient,
		symbols.DefaultClient,
		conf.DefaultClient(),
		resultsBucket,
	)
}

var (
	// TODO - move these into background config
	resultsBucketName             = env.Get("CODEINTEL_RANKING_RESULTS_BUCKET", "lsif-pagerank-experiments", "The GCS bucket.")
	resultsGraphKey               = env.Get("CODEINTEL_RANKING_RESULTS_GRAPH_KEY", "dev", "An identifier of the graph export. Change to start a new import from the configured bucket.")
	resultsObjectKeyPrefix        = env.Get("CODEINTEL_RANKING_RESULTS_OBJECT_KEY_PREFIX", "ranks/", "The object key prefix that holds results of the last PageRank batch job.")
	resultsBucketCredentialsFile  = env.Get("CODEINTEL_RANKING_RESULTS_GOOGLE_APPLICATION_CREDENTIALS_FILE", "", "The path to a service account key file with access to GCS.")
	exportObjectKeyPrefix         = env.Get("CODEINTEL_RANKING_DEVELOPMENT_EXPORT_OBJECT_KEY_PREFIX", "", "The object key prefix that should be used for development exports.")
	developmentExportRepositories = env.Get("CODEINTEL_RANKING_DEVELOPMENT_EXPORT_REPOSITORIES", "github.com/sourcegraph/sourcegraph,github.com/sourcegraph/lsif-go", "Comma-separated list of repositories whose ranks should be exported for development.")

	rankingMapReduceBatchSize   = env.MustGetInt("CODEINTEL_UPLOADS_MAP_REDUCE_RANKING_BATCH_SIZE", 10000, "How many references, definitions, and path counts to map and reduce at once.")
	rankingGraphKey             = env.Get("CODEINTEL_UPLOADS_RANKING_GRAPH_KEY", "dev", "Backdoor value used to restart the ranking export procedure.")
	rankingGraphBatchSize       = env.MustGetInt("CODEINTEL_UPLOADS_RANKING_GRAPH_BATCH_SIZE", 16, "How many uploads to process at once.")
	rankingGraphDeleteBatchSize = env.MustGetInt("CODEINTEL_UPLOADS_RANKING_GRAPH_DELETE_BATCH_SIZE", 32, "How many stale uploads to delete at once.")

	// Backdoor tuning for dotcom
	mergeBatchSize = env.MustGetInt("CODEINTEL_RANKING_MERGE_BATCH_SIZE", 5000, "")
)

func scopedContext(component string, observationCtx *observation.Context) *observation.Context {
	return observation.ScopedContext("codeintel", "ranking", component, observationCtx)
}

func NewGraphExporters(observationCtx *observation.Context, rankingService *Service) []goroutine.BackgroundRoutine {
	return []goroutine.BackgroundRoutine{
		background.NewRankingGraphExporter(
			observationCtx,
			rankingService,
			ConfigExportInst.NumRankingRoutines,
			ConfigExportInst.RankingInterval,
			ConfigExportInst.RankingBatchSize,
			ConfigExportInst.RankingJobsEnabled,
		),
		background.NewRankingGraphMapper(
			observationCtx,
			rankingService,
			ConfigExportInst.NumRankingRoutines,
			ConfigExportInst.RankingInterval,
			ConfigExportInst.RankingJobsEnabled,
		),
		background.NewRankingGraphReducer(
			observationCtx,
			rankingService,
			ConfigExportInst.NumRankingRoutines,
			ConfigExportInst.RankingInterval,
			ConfigExportInst.RankingJobsEnabled,
		),
	}
}
