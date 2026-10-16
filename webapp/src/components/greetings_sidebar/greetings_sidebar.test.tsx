// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {act, render, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {getGreetingCount, resetGreetingCount} from 'client';
import {IntlProvider} from 'react-intl';
import {Provider} from 'react-redux';
import {legacy_createStore as createStore} from 'redux';

import GreetingsSidebar from './greetings_sidebar';

jest.mock('client', () => ({getGreetingCount: jest.fn(), resetGreetingCount: jest.fn()}));
const mockGetGreetingCount = jest.mocked(getGreetingCount);
const mockResetGreetingCount = jest.mocked(resetGreetingCount);

function makeState(channelRoles: string) {
    return {
        entities: {
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
                roles: {'channel-id': new Set(channelRoles.split(' '))},
            },
            roles: {
                roles: {
                    system_user: {permissions: []},
                    channel_user: {permissions: ['read_channel']},
                    channel_admin: {permissions: ['manage_channel_roles']},
                },
            },
        },
    };
}

type State = ReturnType<typeof makeState>;
type SwitchChannel = {type: 'SWITCH_CHANNEL'; channelId: string};

function renderSidebar(channelRoles = 'channel_user') {
    const initialState = makeState(channelRoles);
    const store = createStore((state: State = initialState, action: SwitchChannel | {type: string}) => {
        if (action.type !== 'SWITCH_CHANNEL') {
            return state;
        }
        const {channelId} = action as SwitchChannel;
        return {...state, entities: {...state.entities, channels: {...state.entities.channels, currentChannelId: channelId}}};
    });
    render(
        <Provider store={store}>
            <IntlProvider locale='en'>
                <GreetingsSidebar/>
            </IntlProvider>
        </Provider>,
    );
    return store;
}

beforeEach(() => {
    jest.resetAllMocks();
});

test('shows the current channel\'s greeting count', async () => {
    mockGetGreetingCount.mockResolvedValue(2);

    renderSidebar();

    expect(await screen.findByText('2 greetings in Town Square')).toBeVisible();
    expect(mockGetGreetingCount).toHaveBeenCalledWith('channel-id');
    expect(screen.queryByRole('button', {name: 'Reset'})).not.toBeInTheDocument();
});

test('shows the server error and retries', async () => {
    mockGetGreetingCount.mockRejectedValueOnce(new Error('You don\'t have access to this channel'));
    mockGetGreetingCount.mockResolvedValueOnce(0);

    renderSidebar();

    expect(await screen.findByText('You don\'t have access to this channel')).toBeVisible();
    await userEvent.click(screen.getByRole('button', {name: 'Try again'}));
    expect(await screen.findByText('No greetings in Town Square yet')).toBeVisible();
});

test('lets channel admins reset the count', async () => {
    mockGetGreetingCount.mockResolvedValueOnce(3);
    mockGetGreetingCount.mockResolvedValueOnce(0);
    mockResetGreetingCount.mockResolvedValue();

    renderSidebar('channel_user channel_admin');

    expect(await screen.findByText('3 greetings in Town Square')).toBeVisible();
    await userEvent.click(screen.getByRole('button', {name: 'Reset'}));
    expect(await screen.findByText('No greetings in Town Square yet')).toBeVisible();
    expect(mockResetGreetingCount).toHaveBeenCalledWith('channel-id');
});

test('ignores a reset that finishes after switching channels', async () => {
    let rejectReset: (error: Error) => void = () => {};
    mockResetGreetingCount.mockReturnValue(new Promise((_, reject) => {
        rejectReset = reject;
    }));
    mockGetGreetingCount.mockResolvedValueOnce(3);
    mockGetGreetingCount.mockResolvedValueOnce(5);

    const store = renderSidebar('channel_user channel_admin');
    expect(await screen.findByText('3 greetings in Town Square')).toBeVisible();
    await userEvent.click(screen.getByRole('button', {name: 'Reset'}));

    act(() => {
        store.dispatch({type: 'SWITCH_CHANNEL', channelId: 'other-channel-id'});
    });
    expect(await screen.findByText('5 greetings in Off-Topic')).toBeVisible();

    await act(async () => rejectReset(new Error('Failed to reset greetings')));

    expect(screen.getByText('5 greetings in Off-Topic')).toBeVisible();
    expect(screen.queryByText('Failed to reset greetings')).not.toBeInTheDocument();
    expect(mockGetGreetingCount).toHaveBeenCalledTimes(2);
});
