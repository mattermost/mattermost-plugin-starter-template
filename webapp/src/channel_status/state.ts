// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import manifest from 'manifest';
import type {UnknownAction} from 'redux';

import type {GlobalState} from '@mattermost/types/store';

const STATUS_RECEIVED = `${manifest.id}/status_received`;
const STATUSES_RESET = `${manifest.id}/statuses_reset`;

type StatusReceivedAction = {type: typeof STATUS_RECEIVED; channelId: string; status: string};

export type ChannelStatusState = {

    /** Statuses by channel ID. '' means the channel has none; a missing key means it hasn't been loaded. */
    statuses: Record<string, string>;
};

export function statusReceived(channelId: string, status: string): StatusReceivedAction {
    return {type: STATUS_RECEIVED, channelId, status};
}

export function statusesReset(): UnknownAction {
    return {type: STATUSES_RESET};
}

export function reducer(state: ChannelStatusState = {statuses: {}}, action: UnknownAction): ChannelStatusState {
    switch (action.type) {
    case STATUS_RECEIVED: {
        const {channelId, status} = action as StatusReceivedAction;
        return {statuses: {...state.statuses, [channelId]: status}};
    }
    case STATUSES_RESET:
        return {statuses: {}};
    default:
        return state;
    }
}

export const pluginStateKey = `plugins-${manifest.id}`;

export function getChannelStatus(state: GlobalState, channelId: string): string | undefined {
    const pluginState = (state as unknown as Record<string, ChannelStatusState | undefined>)[pluginStateKey];
    return pluginState?.statuses[channelId];
}
