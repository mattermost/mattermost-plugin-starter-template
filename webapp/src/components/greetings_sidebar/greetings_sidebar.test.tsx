// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {render, screen} from '@testing-library/react';
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
                channels: {'channel-id': {id: 'channel-id', team_id: 'team-id', display_name: 'Town Square', type: 'O'}},
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

function renderSidebar(channelRoles = 'channel_user') {
    const state = makeState(channelRoles);
    return render(
        <Provider store={createStore(() => state)}>
            <IntlProvider locale='en'>
                <GreetingsSidebar/>
            </IntlProvider>
        </Provider>,
    );
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
