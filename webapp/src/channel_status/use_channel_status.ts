// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {fetchChannelStatus} from 'client';
import {useEffect, useState} from 'react';
import {useDispatch, useSelector, useStore} from 'react-redux';

import type {GlobalState} from '@mattermost/types/store';

import {getChannelStatus, statusReceived} from './state';

// The header tag and the sidebar both ask for the same channel; share one request.
const pending = new Map<string, Promise<string>>();

function fetchOnce(channelId: string) {
    let request = pending.get(channelId);
    if (!request) {
        request = fetchChannelStatus(channelId).finally(() => pending.delete(channelId));
        pending.set(channelId, request);
    }
    return request;
}

// Loads a channel's status into the store the first time it's needed. The WebSocket handler
// in index.tsx keeps it current after that.
export function useChannelStatus(channelId: string) {
    const dispatch = useDispatch();
    const store = useStore<GlobalState>();
    const status = useSelector((state: GlobalState) => getChannelStatus(state, channelId));
    const loaded = status !== undefined;
    const [error, setError] = useState<{channelId: string; message: string}>();
    const [attempt, setAttempt] = useState(0);

    useEffect(() => {
        if (loaded) {
            return undefined;
        }

        let stale = false;
        setError(undefined);
        fetchOnce(channelId).then(
            (value) => {
                // A WebSocket event that arrived while this request was in flight is newer.
                if (getChannelStatus(store.getState(), channelId) === undefined) {
                    dispatch(statusReceived(channelId, value));
                }
            },
            (err: Error) => !stale && setError({channelId, message: err.message}),
        );
        return () => {
            stale = true;
        };
    }, [channelId, loaded, attempt, dispatch, store]);

    return {
        status,
        error: !loaded && error?.channelId === channelId ? error.message : undefined,
        retry: () => setAttempt((count) => count + 1),
    };
}
