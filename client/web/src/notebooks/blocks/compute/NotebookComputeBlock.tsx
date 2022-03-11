import { noop } from 'lodash'
import React from 'react'
import ElmComponent from 'react-elm-components'

import { PlatformContext } from '@sourcegraph/shared/src/platform/context'
import { ThemeProps } from '@sourcegraph/shared/src/theme'

import { BlockProps, ComputeBlock } from '../..'
import { useCommonBlockMenuActions } from '../menu/useCommonBlockMenuActions'
import { NotebookBlock } from '../NotebookBlock'

import { Elm } from './component/src/Main.elm'
import styles from './NotebookComputeBlock.module.scss'

interface ComputeBlockProps extends BlockProps<ComputeBlock>, ThemeProps {
    platformContext: Pick<PlatformContext, 'sourcegraphURL'>
}

interface ElmEvent {
    data: string
    eventType?: string
    id?: string
}

function setupPorts(ports: {
    receiveEvent: { send: (event: ElmEvent) => void }
    openStream: { subscribe: (callback: (args: string[]) => void) => void }
}): void {
    const sources: { [key: string]: EventSource } = {}

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    function sendEventToElm(event: any): void {
        // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access
        const elmEvent = { data: event.data, eventType: event.type || null, id: event.id || null }
        ports.receiveEvent.send(elmEvent)
    }

    function newEventSource(address: string): EventSource {
        sources[address] = new EventSource(address)
        return sources[address]
    }

    function deleteAllEventSources(): void {
        for (const [key] of Object.entries(sources)) {
            deleteEventSource(key)
        }
    }

    function deleteEventSource(address: string): void {
        sources[address].close()
        delete sources[address]
    }

    ports.openStream.subscribe((args: string[]) => {
        deleteAllEventSources() // Close any open streams if we receive a request to open a new stream before seeing 'done'.
        console.log(`stream: ${args[0]}`)
        const address = args[0]

        const eventSource = newEventSource(address)
        eventSource.addEventListener('error', () => {
            console.log('EventSource failed')
        })
        eventSource.addEventListener('results', sendEventToElm)
        eventSource.addEventListener('alert', sendEventToElm)
        eventSource.addEventListener('error', sendEventToElm)
        eventSource.addEventListener('done', () => {
            deleteEventSource(address)
            // Note: 'done:true' is sent in progress too. But we want a 'done' for the entire stream in case we don't see it.
            sendEventToElm({ type: 'done', data: '' })
        })
    })
}

export const NotebookComputeBlock: React.FunctionComponent<ComputeBlockProps> = ({
    id,
    input,
    output,
    isSelected,
    isLightTheme,
    platformContext,
    isReadOnly,
    onRunBlock,
    onSelectBlock,
    ...props
}) => {
    const isInputFocused = false
    const commonMenuActions = useCommonBlockMenuActions({
        isInputFocused,
        isReadOnly,
        ...props,
    })

    return (
        <NotebookBlock
            className={styles.input}
            id={id}
            isReadOnly={isReadOnly}
            isInputFocused={isInputFocused}
            aria-label="Notebook compute block"
            onEnterBlock={noop}
            isSelected={isSelected}
            onRunBlock={noop}
            onSelectBlock={onSelectBlock}
            actions={isSelected ? commonMenuActions : []}
            {...props}
        >
            <div className="elm">
                {/* eslint-disable-next-line @typescript-eslint/no-unsafe-assignment, @typescript-eslint/no-unsafe-member-access */}
                <ElmComponent src={Elm.Main} ports={setupPorts} flags={null} />
            </div>
        </NotebookBlock>
    )
}
