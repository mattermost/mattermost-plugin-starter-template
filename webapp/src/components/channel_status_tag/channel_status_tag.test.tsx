// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {act, screen} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {statusReceived} from 'channel_status/state';
import {renderWithStore} from 'channel_status/test_utils';
import {fetchChannelStatus} from 'client';

import type {Channel} from '@mattermost/types/channels';

import ChannelStatusTag from './channel_status_tag';

jest.mock('client', () => ({fetchChannelStatus: jest.fn()}));
const mockFetchChannelStatus = jest.mocked(fetchChannelStatus);

const channel = {id: 'channel-id', team_id: 'team-id', type: 'O'} as Channel;

test('appears when the channel gets a status and opens the sidebar', async () => {
    mockFetchChannelStatus.mockResolvedValue('');
    const onOpen = jest.fn();
    const store = renderWithStore(
        <ChannelStatusTag
            channel={channel}
            onOpen={onOpen}
        />,
    );

    await act(async () => {});
    expect(screen.queryByRole('button')).not.toBeInTheDocument();

    // What the WebSocket handler in index.tsx dispatches when someone changes the status.
    act(() => {
        store.dispatch(statusReceived('channel-id', 'off-track'));
    });
    await userEvent.click(screen.getByRole('button', {name: 'Channel status: Off track'}));

    expect(onOpen).toHaveBeenCalled();
    expect(mockFetchChannelStatus).toHaveBeenCalledTimes(1);
});
