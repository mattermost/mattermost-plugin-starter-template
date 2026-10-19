// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {render} from '@testing-library/react';
import type {ReactElement} from 'react';
import {IntlProvider} from 'react-intl';
import {Provider} from 'react-redux';
import {legacy_createStore as createStore, type UnknownAction} from 'redux';

import {pluginStateKey, reducer, type ChannelStatusState} from './state';

const SWITCH_CHANNEL = 'test/switch_channel';

export function switchChannel(channelId: string) {
    return {type: SWITCH_CHANNEL, channelId};
}

function makeEntities(canEdit: boolean) {
    return {
        general: {config: {}, license: {}},
        preferences: {myPreferences: {}},
        users: {currentUserId: 'user-id', profiles: {'user-id': {id: 'user-id', roles: 'system_user'}}},
        teams: {myMembers: {}},
        channels: {
            currentChannelId: 'channel-id',
            channels: {
                'channel-id': {id: 'channel-id', team_id: 'team-id', display_name: 'Town Square', type: 'O'},
                'other-channel-id': {id: 'other-channel-id', team_id: 'team-id', display_name: 'Off-Topic', type: 'O'},
            },
            roles: {
                'channel-id': new Set(['channel_user']),
                'other-channel-id': new Set(['channel_user']),
            },
        },
        roles: {
            roles: {
                system_user: {permissions: []},
                channel_user: {permissions: canEdit ? ['read_channel', 'manage_public_channel_properties'] : ['read_channel']},
            },
        },
    };
}

type State = {entities: ReturnType<typeof makeEntities>} & Record<string, unknown>;

// Renders with a real Redux store that includes the plugin's reducer, like the host does.
export function renderWithStore(ui: ReactElement, {canEdit = true} = {}) {
    const initialState: State = {entities: makeEntities(canEdit)};
    const store = createStore((state: State = initialState, action: UnknownAction) => {
        let next = state;
        if (action.type === SWITCH_CHANNEL) {
            const channels = {...state.entities.channels, currentChannelId: action.channelId as string};
            next = {...state, entities: {...state.entities, channels}};
        }
        return {...next, [pluginStateKey]: reducer(next[pluginStateKey] as ChannelStatusState | undefined, action)};
    });

    render(
        <Provider store={store}>
            <IntlProvider locale='en'>
                {ui}
            </IntlProvider>
        </Provider>,
    );
    return store;
}
