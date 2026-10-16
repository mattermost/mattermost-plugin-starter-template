// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {getGreetingCount, resetGreetingCount} from 'client';
import {useEffect, useRef, useState} from 'react';
import {FormattedMessage, useIntl} from 'react-intl';
import {useSelector} from 'react-redux';

import RefreshIcon from '@mattermost/compass-icons/components/refresh';
import {Button} from '@mattermost/compass-ui/components/button';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {SectionNotice} from '@mattermost/compass-ui/components/section-notice';
import type {GlobalState} from '@mattermost/types/store';

import {Permissions} from 'mattermost-redux/constants';
import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';
import {haveIChannelPermission} from 'mattermost-redux/selectors/entities/roles';

import './greetings_sidebar.css';

type LoadState =
    | {status: 'loading'}
    | {status: 'loaded'; count: number}
    | {status: 'error'; message: string};

export default function GreetingsSidebar() {
    const intl = useIntl();
    const channel = useSelector(getCurrentChannel);
    const channelId = channel?.id;
    const [state, setState] = useState<LoadState>({status: 'loading'});
    const [reloadCount, setReloadCount] = useState(0);
    const reload = () => setReloadCount((count) => count + 1);

    // Only hides the button; the server makes the real permission check.
    const canReset = useSelector((state: GlobalState) => haveIChannelPermission(state, channel?.team_id, channelId, Permissions.MANAGE_CHANNEL_ROLES));

    const currentChannelId = useRef(channelId);
    useEffect(() => {
        currentChannelId.current = channelId;
    }, [channelId]);

    const reset = () => {
        if (!channelId) {
            return;
        }

        // Ignore the result if the user has switched channels while the reset was in flight.
        const resetChannelId = channelId;
        const stillCurrent = () => currentChannelId.current === resetChannelId;
        setState({status: 'loading'});
        resetGreetingCount(resetChannelId).then(
            () => stillCurrent() && reload(),
            (error: Error) => stillCurrent() && setState({status: 'error', message: error.message}),
        );
    };

    useEffect(() => {
        if (!channelId) {
            return undefined;
        }

        // Ignore responses for a channel the user has already switched away from.
        let stale = false;
        setState({status: 'loading'});
        getGreetingCount(channelId).then(
            (count) => !stale && setState({status: 'loaded', count}),
            (error: Error) => !stale && setState({status: 'error', message: error.message}),
        );
        return () => {
            stale = true;
        };
    }, [channelId, reloadCount]);

    if (state.status === 'error') {
        return (
            <div className='GreetingsSidebar'>
                <SectionNotice
                    type='danger'
                    title={intl.formatMessage({id: 'greetings.error', defaultMessage: 'Couldn’t load greetings'})}
                    description={state.message}
                    primaryButtonLabel={intl.formatMessage({id: 'greetings.retry', defaultMessage: 'Try again'})}
                    onPrimaryAction={reload}
                />
            </div>
        );
    }

    return (
        <div className='GreetingsSidebar'>
            <p className='GreetingsSidebar__text'>
                <FormattedMessage
                    id='greetings.intro'
                    defaultMessage='Run {command} in this channel, then refresh to see the count change.'
                    values={{command: <code key='command'>{'/hello @someone'}</code>}}
                />
            </p>
            <p
                className='GreetingsSidebar__count'
                aria-live='polite'
            >
                {state.status === 'loading' ? (
                    <FormattedMessage
                        id='greetings.loading'
                        defaultMessage='Loading…'
                    />
                ) : (
                    <FormattedMessage
                        id='greetings.count'
                        defaultMessage='{count, plural, =0 {No greetings in {channel} yet} one {# greeting in {channel}} other {# greetings in {channel}}}'
                        values={{count: state.count, channel: channel?.display_name}}
                    />
                )}
            </p>
            <div className='GreetingsSidebar__actions'>
                <Button
                    emphasis='tertiary'
                    size='small'
                    leadingIcon={<Icon glyph={<RefreshIcon/>}/>}
                    loading={state.status === 'loading'}
                    onClick={reload}
                >
                    <FormattedMessage
                        id='greetings.refresh'
                        defaultMessage='Refresh'
                    />
                </Button>
                {canReset && (
                    <Button
                        emphasis='tertiary'
                        size='small'
                        destructive={true}
                        disabled={state.status === 'loading'}
                        onClick={reset}
                    >
                        <FormattedMessage
                            id='greetings.reset'
                            defaultMessage='Reset'
                        />
                    </Button>
                )}
            </div>
        </div>
    );
}
