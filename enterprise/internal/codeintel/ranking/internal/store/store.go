package store

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/keegancsmith/sqlf"
	"github.com/lib/pq"
	otlog "github.com/opentracing/opentracing-go/log"
	logger "github.com/sourcegraph/log"

	"github.com/sourcegraph/sourcegraph/enterprise/internal/codeintel/uploads/shared"
	"github.com/sourcegraph/sourcegraph/internal/api"
	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/basestore"
	"github.com/sourcegraph/sourcegraph/internal/database/batch"
	"github.com/sourcegraph/sourcegraph/internal/database/dbutil"
	"github.com/sourcegraph/sourcegraph/internal/observation"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

// Store provides the interface for ranking storage.
type Store interface {
	// Transactions
	Transact(ctx context.Context) (Store, error)
	Done(err error) error

	GetStarRank(ctx context.Context, repoName api.RepoName) (float64, error)
	GetDocumentRanks(ctx context.Context, repoName api.RepoName) (map[string][2]float64, bool, error)
	LastUpdatedAt(ctx context.Context, repoIDs []api.RepoID) (map[api.RepoID]time.Time, error)
	UpdatedAfter(ctx context.Context, t time.Time) ([]api.RepoName, error)

	InsertDefinitionsForRanking(ctx context.Context, rankingGraphKey string, rankingBatchSize int, definitions []shared.RankingDefinitions) (err error)
	InsertReferencesForRanking(ctx context.Context, rankingGraphKey string, rankingBatchSize int, references shared.RankingReferences) (err error)
	InsertPathCountInputs(ctx context.Context, rankingGraphKey string, batchSize int) (numReferenceRecordsProcessed int, numInputsInserted int, err error)
	InsertPathRanks(ctx context.Context, graphKey string, batchSize int) (numPathRanksInserted float64, numInputsProcessed float64, err error)

	VacuumStaleGraphs(ctx context.Context, derivativeGraphKey string) (
		metadataRecordsDeleted int,
		inputRecordsDeleted int,
		err error,
	)

	VacuumStaleDefinitionsAndReferences(ctx context.Context, graphKey string) (
		numStaleDefinitionRecordsDeleted int,
		numStaleReferenceRecordsDeleted int,
		err error,
	)

	GetUploadsForRanking(ctx context.Context, graphKey, objectPrefix string, batchSize int) ([]shared.ExportedUpload, error)

	ProcessStaleExportedUploads(
		ctx context.Context,
		graphKey string,
		batchSize int,
		deleter func(ctx context.Context, objectPrefix string) error,
	) (totalDeleted int, err error)
}

// store manages the ranking store.
type store struct {
	db         *basestore.Store
	logger     logger.Logger
	operations *operations
}

// New returns a new ranking store.
func New(observationCtx *observation.Context, db database.DB) Store {
	return newInternal(observationCtx, db)
}

func newInternal(observationCtx *observation.Context, db database.DB) *store {
	return &store{
		db:         basestore.NewWithHandle(db.Handle()),
		logger:     logger.Scoped("ranking.store", ""),
		operations: newOperations(observationCtx),
	}
}

func (s *store) Transact(ctx context.Context) (Store, error) {
	return s.transact(ctx)
}

func (s *store) transact(ctx context.Context) (*store, error) {
	tx, err := s.db.Transact(ctx)
	if err != nil {
		return nil, err
	}

	return &store{
		logger:     s.logger,
		db:         tx,
		operations: s.operations,
	}, nil
}

func (s *store) Done(err error) error {
	return s.db.Done(err)
}

func (s *store) GetStarRank(ctx context.Context, repoName api.RepoName) (float64, error) {
	rank, _, err := basestore.ScanFirstFloat(s.db.Query(ctx, sqlf.Sprintf(getStarRankQuery, repoName)))
	return rank, err
}

const getStarRankQuery = `
SELECT
	s.rank
FROM (
	SELECT
		name,
		percent_rank() OVER (ORDER BY stars) AS rank
	FROM repo
) s
WHERE s.name = %s
`

func (s *store) GetDocumentRanks(ctx context.Context, repoName api.RepoName) (map[string][2]float64, bool, error) {
	pathRanksWithPrecision := map[string][2]float64{}
	scanner := func(s dbutil.Scanner) (bool, error) {
		var (
			precision  float64
			serialized string
		)
		if err := s.Scan(&precision, &serialized); err != nil {
			return false, err
		}

		pathRanks := map[string]float64{}
		if err := json.Unmarshal([]byte(serialized), &pathRanks); err != nil {
			return false, err
		}

		for path, newRank := range pathRanks {
			if oldRank, ok := pathRanksWithPrecision[path]; ok && oldRank[0] <= precision {
				continue
			}

			pathRanksWithPrecision[path] = [2]float64{precision, newRank}
		}

		return true, nil
	}

	if err := basestore.NewCallbackScanner(scanner)(s.db.Query(ctx, sqlf.Sprintf(getDocumentRanksQuery, repoName))); err != nil {
		return nil, false, err
	}
	return pathRanksWithPrecision, true, nil
}

const getDocumentRanksQuery = `
SELECT
	precision,
	payload
FROM codeintel_path_ranks pr
JOIN repo r ON r.id = pr.repository_id
WHERE
	r.name = %s AND
	r.deleted_at IS NULL AND
	r.blocked IS NULL
`

func (s *store) LastUpdatedAt(ctx context.Context, repoIDs []api.RepoID) (map[api.RepoID]time.Time, error) {
	pairs, err := scanLastUpdatedAtPairs(s.db.Query(ctx, sqlf.Sprintf(lastUpdatedAtQuery, pq.Array(repoIDs))))
	if err != nil {
		return nil, err
	}

	return pairs, nil
}

const lastUpdatedAtQuery = `
SELECT
	repository_id,
	updated_at
FROM codeintel_path_ranks
WHERE repository_id = ANY(%s)
`

var scanLastUpdatedAtPairs = basestore.NewMapScanner(func(s dbutil.Scanner) (repoID api.RepoID, t time.Time, _ error) {
	err := s.Scan(&repoID, &t)
	return repoID, t, err
})

func (s *store) UpdatedAfter(ctx context.Context, t time.Time) ([]api.RepoName, error) {
	names, err := basestore.ScanStrings(s.db.Query(ctx, sqlf.Sprintf(updatedAfterQuery, t)))
	if err != nil {
		return nil, err
	}

	repoNames := make([]api.RepoName, 0, len(names))
	for _, name := range names {
		repoNames = append(repoNames, api.RepoName(name))
	}

	return repoNames, nil
}

const updatedAfterQuery = `
SELECT r.name
FROM codeintel_path_ranks pr
JOIN repo r ON r.id = pr.repository_id
WHERE pr.updated_at >= %s
ORDER BY r.name
`

func (s *store) SetDocumentRanks(ctx context.Context, repoName api.RepoName, precision float64, ranks map[string]float64) error {
	serialized, err := json.Marshal(ranks)
	if err != nil {
		return err
	}

	return s.db.Exec(ctx, sqlf.Sprintf(setDocumentRanksQuery, repoName, precision, serialized))
}

const setDocumentRanksQuery = `
	INSERT INTO codeintel_path_ranks AS pr (repository_id, precision, payload)
	VALUES (
		(SELECT id FROM repo WHERE name = %s),
		%s,
		%s
	)
	ON CONFLICT (repository_id, precision) DO
	UPDATE
		SET payload = EXCLUDED.payload
	`

func (s *store) InsertDefinitionsForRanking(
	ctx context.Context,
	rankingGraphKey string,
	rankingBatchNumber int,
	definitions []shared.RankingDefinitions,
) (err error) {
	ctx, _, endObservation := s.operations.insertDefinitionsForRanking.With(
		ctx,
		&err,
		observation.Args{},
	)
	defer endObservation(1, observation.Args{})

	tx, err := s.db.Transact(ctx)
	if err != nil {
		return err
	}
	defer func() { err = tx.Done(err) }()

	inserter := func(inserter *batch.Inserter) error {
		batchDefinitions := make([]shared.RankingDefinitions, 0, rankingBatchNumber)
		for _, def := range definitions {
			batchDefinitions = append(batchDefinitions, def)

			if len(batchDefinitions) == rankingBatchNumber {
				if err := insertDefinitions(ctx, inserter, rankingGraphKey, batchDefinitions); err != nil {
					return err
				}
				batchDefinitions = make([]shared.RankingDefinitions, 0, rankingBatchNumber)
			}
		}

		if len(batchDefinitions) > 0 {
			if err := insertDefinitions(ctx, inserter, rankingGraphKey, batchDefinitions); err != nil {
				return err
			}
		}

		return nil
	}

	if err := batch.WithInserter(
		ctx,
		tx.Handle(),
		"codeintel_ranking_definitions",
		batch.MaxNumPostgresParameters,
		[]string{
			"upload_id",
			"symbol_name",
			"repository",
			"document_path",
			"graph_key",
		},
		inserter,
	); err != nil {
		return err
	}

	return nil
}

func insertDefinitions(
	ctx context.Context,
	inserter *batch.Inserter,
	rankingGraphKey string,
	definitions []shared.RankingDefinitions,
) error {
	for _, def := range definitions {
		if err := inserter.Insert(
			ctx,
			def.UploadID,
			def.SymbolName,
			def.Repository,
			def.DocumentPath,
			rankingGraphKey,
		); err != nil {
			return err
		}
	}
	return nil
}

func (s *store) InsertReferencesForRanking(
	ctx context.Context,
	rankingGraphKey string,
	rankingBatchNumber int,
	references shared.RankingReferences,
) (err error) {
	ctx, _, endObservation := s.operations.insertReferencesForRanking.With(
		ctx,
		&err,
		observation.Args{},
	)
	defer endObservation(1, observation.Args{})

	tx, err := s.db.Transact(ctx)
	if err != nil {
		return err
	}
	defer func() { err = tx.Done(err) }()

	inserter := func(inserter *batch.Inserter) error {
		batchSymbolNames := make([]string, 0, rankingBatchNumber)
		for _, ref := range references.SymbolNames {
			batchSymbolNames = append(batchSymbolNames, ref)

			if len(batchSymbolNames) == rankingBatchNumber {
				if err := inserter.Insert(ctx, references.UploadID, pq.Array(batchSymbolNames), rankingGraphKey); err != nil {
					return err
				}
				batchSymbolNames = make([]string, 0, rankingBatchNumber)
			}
		}

		if len(batchSymbolNames) > 0 {
			if err := inserter.Insert(ctx, references.UploadID, pq.Array(batchSymbolNames), rankingGraphKey); err != nil {
				return err
			}
		}

		return nil
	}

	if err := batch.WithInserter(
		ctx,
		tx.Handle(),
		"codeintel_ranking_references",
		batch.MaxNumPostgresParameters,
		[]string{"upload_id", "symbol_names", "graph_key"},
		inserter,
	); err != nil {
		return err
	}

	return nil
}

func (s *store) InsertPathCountInputs(
	ctx context.Context,
	derivativeGraphKey string,
	batchSize int,
) (
	numReferenceRecordsProcessed int,
	numInputsInserted int,
	err error,
) {
	ctx, _, endObservation := s.operations.insertPathCountInputs.With(ctx, &err, observation.Args{})
	defer endObservation(1, observation.Args{})

	parts := strings.Split(derivativeGraphKey, "-")
	if len(parts) < 2 {
		return 0, 0, errors.Newf("unexpected graph key format %q", derivativeGraphKey)
	}
	// Remove last segment, which indicates the current time bucket
	parentGraphKey := strings.Join(parts[:len(parts)-1], "-")

	rows, err := s.db.Query(ctx, sqlf.Sprintf(
		insertPathCountInputsQuery,
		parentGraphKey,
		derivativeGraphKey,
		batchSize,
		derivativeGraphKey,
		parentGraphKey,
		derivativeGraphKey,
	))
	if err != nil {
		return 0, 0, err
	}
	defer func() { err = basestore.CloseRows(rows, err) }()

	for rows.Next() {
		if err := rows.Scan(
			&numReferenceRecordsProcessed,
			&numInputsInserted,
		); err != nil {
			return 0, 0, err
		}
	}

	return numReferenceRecordsProcessed, numInputsInserted, nil
}

const insertPathCountInputsQuery = `
WITH
refs AS (
	SELECT
		rr.id,
		rr.symbol_names
	FROM codeintel_ranking_references rr
	WHERE
		rr.graph_key = %s AND
		NOT EXISTS (
			SELECT 1
			FROM codeintel_ranking_references_processed rrp
			WHERE
				rrp.graph_key = %s AND
				rrp.codeintel_ranking_reference_id = rr.id
		)
	ORDER BY rr.id
	LIMIT %s
),
locked_refs AS (
	INSERT INTO codeintel_ranking_references_processed (graph_key, codeintel_ranking_reference_id)
	SELECT %s, r.id FROM refs r
	ON CONFLICT DO NOTHING
	RETURNING codeintel_ranking_reference_id
),
referenced_symbols AS (
	SELECT unnest(symbol_names) AS symbol_name
	FROM refs
	WHERE id IN (SELECT id FROM locked_refs)
),
referenced_definitions AS (
	SELECT
		rd.repository,
		rd.document_path,
		rd.graph_key
	FROM codeintel_ranking_definitions rd
	WHERE
		rd.graph_key = %s AND
		rd.symbol_name IN (SELECT symbol_name FROM referenced_symbols)
),
ins AS (
	INSERT INTO codeintel_ranking_path_counts_inputs (repository, document_path, count, graph_key)
	SELECT
		rd.repository,
		rd.document_path,
		COUNT(*),
		%s
	FROM referenced_definitions rd
	GROUP BY rd.repository, rd.document_path, rd.graph_key
	RETURNING 1
)
SELECT
	(SELECT COUNT(*) FROM locked_refs),
	(SELECT COUNT(*) FROM ins)
`

func (s *store) InsertPathRanks(
	ctx context.Context,
	derivativeGraphKey string,
	batchSize int,
) (numPathRanksInserted float64, numInputsProcessed float64, err error) {
	ctx, _, endObservation := s.operations.insertPathRanks.With(
		ctx,
		&err,
		observation.Args{LogFields: []otlog.Field{
			otlog.String("derivativeGraphKey", derivativeGraphKey),
		}},
	)
	defer endObservation(1, observation.Args{})

	// Validation for testing (test graph keys do not have hyphens)
	if parts := strings.Split(derivativeGraphKey, "-"); len(parts) < 2 {
		return 0, 0, errors.Newf("unexpected graph key format %q", derivativeGraphKey)
	}

	rows, err := s.db.Query(ctx, sqlf.Sprintf(insertPathRanksQuery, derivativeGraphKey, batchSize))
	if err != nil {
		return 0, 0, err
	}
	defer func() { err = basestore.CloseRows(rows, err) }()

	if !rows.Next() {
		return 0, 0, errors.New("no rows from count")
	}

	if err = rows.Scan(&numPathRanksInserted, &numInputsProcessed); err != nil {
		return 0, 0, err
	}

	return numPathRanksInserted, numInputsProcessed, nil
}

const insertPathRanksQuery = `
WITH
input_ranks AS (
	SELECT
		pci.id,
		r.id AS repository_id,
		pci.document_path AS path,
		pci.count
	FROM codeintel_ranking_path_counts_inputs pci
	JOIN repo r ON lower(r.name) = (lower(pci.repository::text) COLLATE "C")
	WHERE
		pci.graph_key = %s AND
		NOT pci.processed AND
		r.deleted_at IS NULL AND
		r.blocked IS NULL
	ORDER BY pci.graph_key, pci.repository, pci.id
	LIMIT %s
	FOR UPDATE SKIP LOCKED
),
processed AS (
	UPDATE codeintel_ranking_path_counts_inputs
	SET processed = true
	WHERE id IN (SELECT ir.id FROM input_ranks ir)
	RETURNING 1
),
inserted AS (
	INSERT INTO codeintel_path_ranks AS pr (repository_id, precision, payload)
	SELECT
		temp.repository_id,
		1,
		sg_jsonb_concat_agg(temp.row)
	FROM (
		SELECT
			cr.repository_id,
			jsonb_build_object(cr.path, SUM(count)) AS row
		FROM input_ranks cr
		GROUP BY cr.repository_id, cr.path
	) temp
	GROUP BY temp.repository_id
	ON CONFLICT (repository_id, precision) DO UPDATE SET
		graph_key = EXCLUDED.graph_key,
		payload = CASE
			WHEN pr.graph_key != EXCLUDED.graph_key
				THEN EXCLUDED.payload
			ELSE
				(
					SELECT sg_jsonb_concat_agg(row) FROM (
						SELECT jsonb_build_object(key, SUM(value::int)) AS row
						FROM
							(
								SELECT * FROM jsonb_each(pr.payload)
								UNION
								SELECT * FROM jsonb_each(EXCLUDED.payload)
							) AS both_payloads
						GROUP BY key
					) AS combined_json
				)
			END
	RETURNING 1
)
SELECT
	(SELECT COUNT(*) FROM processed) AS num_processed,
	(SELECT COUNT(*) FROM inserted) AS num_inserted
`

func (s *store) VacuumStaleDefinitionsAndReferences(ctx context.Context, graphKey string) (
	numStaleDefinitionRecordsDeleted int,
	numStaleReferenceRecordsDeleted int,
	err error,
) {
	ctx, _, endObservation := s.operations.vacuumStaleDefinitionsAndReferences.With(ctx, &err, observation.Args{LogFields: []otlog.Field{}})
	defer endObservation(1, observation.Args{})

	rows, err := s.db.Query(ctx, sqlf.Sprintf(vacuumStaleDefinitionsAndReferencesQuery, graphKey, graphKey))
	if err != nil {
		return 0, 0, err
	}
	defer func() { err = basestore.CloseRows(rows, err) }()

	for rows.Next() {
		if err := rows.Scan(
			&numStaleDefinitionRecordsDeleted,
			&numStaleReferenceRecordsDeleted,
		); err != nil {
			return 0, 0, err
		}
	}

	return numStaleDefinitionRecordsDeleted, numStaleReferenceRecordsDeleted, nil
}

const vacuumStaleDefinitionsAndReferencesQuery = `
WITH
locked_definitions AS (
	SELECT id
	FROM codeintel_ranking_definitions
	WHERE
		upload_id NOT IN (SELECT uvt.upload_id FROM lsif_uploads_visible_at_tip uvt WHERE uvt.is_default_branch) AND
		graph_key = %s
	ORDER BY id
	FOR UPDATE
),
locked_references AS (
	SELECT id
	FROM codeintel_ranking_references
	WHERE
		upload_id NOT IN (SELECT uvt.upload_id FROM lsif_uploads_visible_at_tip uvt WHERE uvt.is_default_branch) AND
		graph_key = %s
	ORDER BY id
	FOR UPDATE
),
deleted_definitions AS (
	DELETE FROM codeintel_ranking_definitions
	WHERE id IN (SELECT id FROM locked_definitions)
	RETURNING 1
),
deleted_references AS (
	DELETE FROM codeintel_ranking_references
	WHERE id IN (SELECT id FROM locked_references)
	RETURNING 1
)
SELECT
	(SELECT COUNT(*) FROM deleted_definitions),
	(SELECT COUNT(*) FROM deleted_references)
`

func (s *store) VacuumStaleGraphs(ctx context.Context, derivativeGraphKey string) (
	metadataRecordsDeleted int,
	inputRecordsDeleted int,
	err error,
) {
	ctx, _, endObservation := s.operations.vacuumStaleGraphs.With(ctx, &err, observation.Args{LogFields: []otlog.Field{}})
	defer endObservation(1, observation.Args{})

	rows, err := s.db.Query(ctx, sqlf.Sprintf(vacuumStaleGraphsQuery, derivativeGraphKey, derivativeGraphKey))
	if err != nil {
		return 0, 0, err
	}
	defer func() { err = basestore.CloseRows(rows, err) }()

	for rows.Next() {
		if err := rows.Scan(
			&metadataRecordsDeleted,
			&inputRecordsDeleted,
		); err != nil {
			return 0, 0, err
		}
	}

	return metadataRecordsDeleted, inputRecordsDeleted, nil
}

const vacuumStaleGraphsQuery = `
WITH
locked_references_processed AS (
	SELECT id
	FROM codeintel_ranking_references_processed
	WHERE graph_key != %s
	ORDER BY id
	FOR UPDATE
),
locked_path_counts_inputs AS (
	SELECT id
	FROM codeintel_ranking_path_counts_inputs
	WHERE graph_key != %s
	ORDER BY id
	FOR UPDATE
),
deleted_references_processed AS (
	DELETE FROM codeintel_ranking_references_processed
	WHERE id IN (SELECT id FROM locked_references_processed)
	RETURNING 1
),
deleted_path_counts_inputs AS (
	DELETE FROM codeintel_ranking_path_counts_inputs
	WHERE id IN (SELECT id FROM locked_path_counts_inputs)
	RETURNING 1
)
SELECT
	(SELECT COUNT(*) FROM deleted_references_processed),
	(SELECT COUNT(*) FROM deleted_path_counts_inputs)
`
