import React, { FC, useCallback, useMemo, useRef } from 'react'

import { useLocation, useNavigate } from 'react-router-dom'
import { NavbarQueryState } from 'src/stores/navbarSearchQueryState'
import shallow from 'zustand/shallow'

import { SearchBox, Toggles } from '@sourcegraph/branded'
// The experimental search input should be shown on the search home page
// eslint-disable-next-line no-restricted-imports
import { LazyCodeMirrorQueryInput } from '@sourcegraph/branded/src/search-ui/experimental'
import { TraceSpanProvider } from '@sourcegraph/observability-client'
import {
    CaseSensitivityProps,
    SearchPatternTypeProps,
    SubmitSearchParameters,
    canSubmitSearch,
    QueryState,
    SearchModeProps,
    getUserSearchContextNamespaces,
} from '@sourcegraph/shared/src/search'
import { useIsLightTheme } from '@sourcegraph/shared/src/theme'
import { Form } from '@sourcegraph/wildcard'

import { Notices } from '../../../../global/Notices'
import { useLegacyRouteContext } from '../../../../LegacyRouteContext'
import { submitSearch } from '../../../../search/helpers'
import { useLazyCreateSuggestions, useLazyHistoryExtension } from '../../../../search/input/lazy'
import { useRecentSearches } from '../../../../search/input/useRecentSearches'
import { useExperimentalQueryInput } from '../../../../search/useExperimentalSearchInput'
import {
    useExperimentalFeatures,
    useNavbarQueryState,
    setSearchCaseSensitivity,
    setSearchPatternType,
    setSearchMode,
} from '../../../../stores'

import styles from './SearchPageInput.module.scss'

const queryStateSelector = (
    state: NavbarQueryState
): Pick<CaseSensitivityProps, 'caseSensitive'> & SearchPatternTypeProps & Pick<SearchModeProps, 'searchMode'> => ({
    caseSensitive: state.searchCaseSensitivity,
    patternType: state.searchPatternType,
    searchMode: state.searchMode,
})

interface SearchPageInputProps {
    queryState: QueryState
    setQueryState: (newState: QueryState) => void
    hardCodedSearchContextSpec?: string
}

export const SearchPageInput: FC<SearchPageInputProps> = props => {
    const { queryState, setQueryState, hardCodedSearchContextSpec } = props

    const {
        authenticatedUser,
        globbing,
        isSourcegraphDotCom,
        telemetryService,
        platformContext,
        searchContextsEnabled,
        settingsCascade,
        selectedSearchContextSpec: dynamicSearchContextSpec,
        fetchSearchContexts,
        setSelectedSearchContextSpec,
    } = useLegacyRouteContext()

    const selectedSearchContextSpec = hardCodedSearchContextSpec || dynamicSearchContextSpec

    const location = useLocation()
    const navigate = useNavigate()

    const isLightTheme = useIsLightTheme()
    const { caseSensitive, patternType, searchMode } = useNavbarQueryState(queryStateSelector, shallow)
    const [experimentalQueryInput] = useExperimentalQueryInput()
    const applySuggestionsOnEnter =
        useExperimentalFeatures(features => features.applySearchQuerySuggestionOnEnter) ?? true

    const { recentSearches } = useRecentSearches()
    const recentSearchesRef = useRef(recentSearches)
    recentSearchesRef.current = recentSearches

    const submitSearchOnChange = useCallback(
        (parameters: Partial<SubmitSearchParameters> = {}) => {
            const query = queryState.query

            if (canSubmitSearch(query, selectedSearchContextSpec)) {
                submitSearch({
                    source: 'home',
                    query,
                    historyOrNavigate: navigate,
                    location,
                    patternType,
                    caseSensitive,
                    searchMode,
                    // In the new query input, context is either omitted (-> global)
                    // or explicitly specified.
                    selectedSearchContextSpec: experimentalQueryInput ? undefined : selectedSearchContextSpec,
                    ...parameters,
                })
            }
        },
        [
            queryState.query,
            selectedSearchContextSpec,
            navigate,
            location,
            patternType,
            caseSensitive,
            searchMode,
            experimentalQueryInput,
        ]
    )
    const submitSearchOnChangeRef = useRef(submitSearchOnChange)
    submitSearchOnChangeRef.current = submitSearchOnChange

    const onSubmit = useCallback(
        (event?: React.FormEvent): void => {
            event?.preventDefault()
            submitSearchOnChangeRef.current()
        },
        [submitSearchOnChangeRef]
    )

    const suggestionSource = useLazyCreateSuggestions(
        experimentalQueryInput,
        useMemo(
            () => ({
                platformContext,
                authenticatedUser,
                fetchSearchContexts,
                getUserSearchContextNamespaces,
                isSourcegraphDotCom,
            }),
            [platformContext, authenticatedUser, fetchSearchContexts, isSourcegraphDotCom]
        )
    )

    const experimentalExtensions = useLazyHistoryExtension(
        experimentalQueryInput,
        recentSearchesRef,
        submitSearchOnChangeRef
    )

    // TODO (#48103): Remove/simplify when new search input is released
    const input = experimentalQueryInput ? (
        <LazyCodeMirrorQueryInput
            patternType={patternType}
            interpretComments={false}
            queryState={queryState}
            onChange={setQueryState}
            onSubmit={onSubmit}
            isLightTheme={isLightTheme}
            placeholder="Search for code or files..."
            suggestionSource={suggestionSource}
            extensions={experimentalExtensions}
        >
            <Toggles
                patternType={patternType}
                caseSensitive={caseSensitive}
                setPatternType={setSearchPatternType}
                setCaseSensitivity={setSearchCaseSensitivity}
                searchMode={searchMode}
                setSearchMode={setSearchMode}
                settingsCascade={settingsCascade}
                navbarSearchQuery={queryState.query}
                showCopyQueryButton={false}
                showSmartSearchButton={false}
            />
        </LazyCodeMirrorQueryInput>
    ) : (
        <SearchBox
            platformContext={platformContext}
            globbing={globbing}
            getUserSearchContextNamespaces={getUserSearchContextNamespaces}
            fetchSearchContexts={fetchSearchContexts}
            selectedSearchContextSpec={selectedSearchContextSpec}
            setSelectedSearchContextSpec={setSelectedSearchContextSpec}
            telemetryService={telemetryService}
            authenticatedUser={authenticatedUser}
            isSourcegraphDotCom={isSourcegraphDotCom}
            settingsCascade={settingsCascade}
            searchContextsEnabled={searchContextsEnabled}
            showSearchContext={searchContextsEnabled}
            showSearchContextManagement={true}
            caseSensitive={caseSensitive}
            patternType={patternType}
            setPatternType={setSearchPatternType}
            setCaseSensitivity={setSearchCaseSensitivity}
            searchMode={searchMode}
            setSearchMode={setSearchMode}
            queryState={queryState}
            onChange={setQueryState}
            onSubmit={onSubmit}
            isExternalServicesUserModeAll={window.context.externalServicesUserMode === 'all'}
            structuralSearchDisabled={window.context?.experimentalFeatures?.structuralSearch === 'disabled'}
            applySuggestionsOnEnter={applySuggestionsOnEnter}
            showSearchHistory={true}
            recentSearches={recentSearches}
        />
    )
    return (
        <div className="d-flex flex-row flex-shrink-past-contents">
            <Form className="flex-grow-1 flex-shrink-past-contents" onSubmit={onSubmit}>
                <div data-search-page-input-container={true} className={styles.inputContainer}>
                    <TraceSpanProvider name="SearchBox">
                        <div className="d-flex flex-grow-1 w-100">{input}</div>
                    </TraceSpanProvider>
                </div>
                <Notices className="my-3 text-center" location="home" settingsCascade={settingsCascade} />
            </Form>
        </div>
    )
}
