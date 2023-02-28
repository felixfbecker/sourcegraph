package ranking

import (
	"context"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/ranking/internal/store"
	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/uploads"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/conf/conftypes"
	"github.com/sourcegraph/sourcegraph/internal/lazyregexp"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/schema"
)

type Service struct {
	store           store.Store
	uploadSvc       *uploads.Service
	gitserverClient GitserverClient
	symbolsClient   SymbolsClient
	getConf         conftypes.SiteConfigQuerier
	resultsBucket   *storage.BucketHandle
	operations      *operations
	logger          log.Logger
}

func newService(
	observationCtx *observation.Context,
	store store.Store,
	uploadSvc *uploads.Service,
	gitserverClient GitserverClient,
	symbolsClient SymbolsClient,
	getConf conftypes.SiteConfigQuerier,
	resultsBucket *storage.BucketHandle,
) *Service {
	return &Service{
		store:           store,
		uploadSvc:       uploadSvc,
		gitserverClient: gitserverClient,
		symbolsClient:   symbolsClient,
		getConf:         getConf,
		resultsBucket:   resultsBucket,
		operations:      newOperations(observationCtx),
		logger:          observationCtx.Logger,
	}
}

// GetRepoRank returns a rank vector for the given repository. Repositories are assumed to
// be ordered by each pairwise component of the resulting vector, higher ranks coming earlier.
// We currently rank first by user-defined scores, then by GitHub star count.
func (s *Service) GetRepoRank(ctx context.Context, repoName api.RepoName) (_ []float64, err error) {
	_, _, endObservation := s.operations.getRepoRank.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})

	userRank := repoRankFromConfig(s.getConf.SiteConfig(), string(repoName))

	starRank, err := s.store.GetStarRank(ctx, repoName)
	if err != nil {
		return nil, err
	}

	return []float64{squashRange(userRank), starRank}, nil
}

// copy pasta
// https://github.com/sourcegraph/sourcegraph/blob/942c417363b07c9e0a6377456f1d6a80a94efb99/cmd/frontend/internal/httpapi/search.go#L172
func repoRankFromConfig(siteConfig schema.SiteConfiguration, repoName string) float64 {
	val := 0.0
	if siteConfig.ExperimentalFeatures == nil || siteConfig.ExperimentalFeatures.Ranking == nil {
		return val
	}
	scores := siteConfig.ExperimentalFeatures.Ranking.RepoScores
	if len(scores) == 0 {
		return val
	}
	// try every "directory" in the repo name to assign it a value, so a repoName like
	// "github.com/sourcegraph/zoekt" will have "github.com", "github.com/sourcegraph",
	// and "github.com/sourcegraph/zoekt" tested.
	for i := 0; i < len(repoName); i++ {
		if repoName[i] == '/' {
			val += scores[repoName[:i]]
		}
	}
	val += scores[repoName]
	return val
}

var allPathsPattern = lazyregexp.New(".*")

type RepoPathRanks struct {
	MedianReferenceCount int
	Paths                map[string]PathRank
}

type PathRank struct {
	Count int
}

// GetDocumentRank returns a map from paths within the given repo to their rank vector.
func (s *Service) GetDocumentRanks(ctx context.Context, repoName api.RepoName) (_ RepoPathRanks, err error) {
	_, _, endObservation := s.operations.getDocumentRanks.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})

	documentRanks, ok, err := s.store.GetDocumentRanks(ctx, repoName)
	if err != nil {
		return RepoPathRanks{}, err
	}
	if !ok {
		return RepoPathRanks{}, nil
	}

	pathRanks := make(map[string]PathRank, len(documentRanks))
	for path, rank := range documentRanks {
		pathRanks[path] = PathRank{
			Count: int(rank[1]), // TODO
		}
	}

	return RepoPathRanks{
		MedianReferenceCount: 0, // TODO
		Paths:                pathRanks,
	}, nil
}

func (s *Service) LastUpdatedAt(ctx context.Context, repoIDs []api.RepoID) (map[api.RepoID]time.Time, error) {
	return s.store.LastUpdatedAt(ctx, repoIDs)
}

func (s *Service) UpdatedAfter(ctx context.Context, t time.Time) ([]api.RepoName, error) {
	return s.store.UpdatedAfter(ctx, t)
}

func isPathGenerated(path string) bool {
	return strings.HasSuffix(path, "min.js") || strings.HasSuffix(path, "js.map")
}

func isPathVendored(path string) bool {
	return strings.Contains(path, "vendor/") || strings.Contains(path, "node_modules/")
}

var testPattern = lazyregexp.New("test")

func isPathTest(path string) bool {
	return testPattern.MatchString(path)
}

// Converts a boolean to a [0, 1] rank (where true is ordered before false).
func boolRank(v bool) float64 {
	if v {
		return 1.0
	}

	return 0.0
}

// squashRange maps a value in the range [0, inf) to a value in the range
// [0, 1) monotonically (i.e., (a < b) <-> (squashRange(a) < squashRange(b))).
func squashRange(j float64) float64 {
	return j / (1 + j)
}
