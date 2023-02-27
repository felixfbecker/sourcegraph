import { ApolloError } from '@apollo/client'

import { useQuery } from '@sourcegraph/http-client'

import {
    RepoCodeIntelStatusResult,
    RepoCodeIntelStatusVariables,
    InferredAvailableIndexersFields,
    PreciseIndexFields,
    RepoCodeIntelStatusSummaryFields,
    RepoCodeIntelStatusCommitGraphFields,
} from '../../../../graphql-operations'
import { repoCodeIntelStatusQuery } from '../backend'

export interface UseRepoCodeIntelStatusPayload {
    lastIndexScan?: string
    lastUploadRetentionScan?: string
    availableIndexers: InferredAvailableIndexersFields[]
    recentActivity: PreciseIndexFields[]
}

interface UseRepoCodeIntelStatusResult {
    data?: {
        summary: RepoCodeIntelStatusSummaryFields
        commitGraph: RepoCodeIntelStatusCommitGraphFields
    }
    error?: ApolloError
    loading: boolean
}

export const useRepoCodeIntelStatus = (variables: RepoCodeIntelStatusVariables): UseRepoCodeIntelStatusResult => {
    const { data, error, loading } = useQuery<RepoCodeIntelStatusResult, RepoCodeIntelStatusVariables>(
        repoCodeIntelStatusQuery,
        {
            variables,
            fetchPolicy: 'cache-and-network',
        }
    )

    const repo = data?.repository

    if (!repo) {
        return { loading, error }
    }

    return {
        data: {
            summary: repo.codeIntelSummary,
            commitGraph: repo.codeIntelligenceCommitGraph,
        },
        error,
        loading,
    }
}
