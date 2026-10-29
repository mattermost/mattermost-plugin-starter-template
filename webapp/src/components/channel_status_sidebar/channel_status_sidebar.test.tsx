// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {act, screen, within} from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import {statusReceived} from 'channel_status/state';
import {renderWithStore, switchChannel} from 'channel_status/test_utils';
import {clearChannelStatus, fetchChannelStatus, setChannelStatus} from 'client';

import ChannelStatusSidebar from './channel_status_sidebar';

jest.mock('client', () => ({fetchChannelStatus: jest.fn(), setChannelStatus: jest.fn(), clearChannelStatus: jest.fn()}));
const mockFetchChannelStatus = jest.mocked(fetchChannelStatus);
const mockSetChannelStatus = jest.mocked(setChannelStatus);
const mockClearChannelStatus = jest.mocked(clearChannelStatus);

const currentStatus = async () => within(await screen.findByRole('region', {name: 'Current status'}));

beforeEach(() => {
    jest.resetAllMocks();
});

test('lets members with permission change the status', async () => {
    mockFetchChannelStatus.mockResolvedValue('at-risk');
    mockSetChannelStatus.mockResolvedValue();
    renderWithStore(<ChannelStatusSidebar/>);

    expect(await (await currentStatus()).findByText('At risk')).toBeVisible();
    expect(mockFetchChannelStatus).toHaveBeenCalledWith('channel-id');

    await userEvent.click(screen.getByRole('combobox', {name: 'New status'}));
    await userEvent.click(screen.getByRole('option', {name: 'Off track'}));
    await userEvent.click(screen.getByRole('button', {name: 'Save'}));

    expect(mockSetChannelStatus).toHaveBeenCalledWith('channel-id', 'off-track');
    expect(await (await currentStatus()).findByText('Off track')).toBeVisible();
});

test('clears the status', async () => {
    mockFetchChannelStatus.mockResolvedValue('on-track');
    mockClearChannelStatus.mockResolvedValue();
    renderWithStore(<ChannelStatusSidebar/>);

    await userEvent.click(await screen.findByRole('button', {name: 'Clear status'}));

    expect(mockClearChannelStatus).toHaveBeenCalledWith('channel-id');
    expect(await (await currentStatus()).findByText('No status set')).toBeVisible();
});

test('is read-only for members without permission', async () => {
    mockFetchChannelStatus.mockResolvedValue('');
    renderWithStore(<ChannelStatusSidebar/>, {canEdit: false});

    expect(await screen.findByText('You don’t have permission to change this channel’s status.')).toBeVisible();
    expect(screen.queryByRole('button', {name: 'Save'})).not.toBeInTheDocument();
});

test('shows load errors and retries', async () => {
    mockFetchChannelStatus.mockRejectedValueOnce(new Error('Failed to load the channel status'));
    mockFetchChannelStatus.mockResolvedValueOnce('on-track');
    renderWithStore(<ChannelStatusSidebar/>);

    expect(await screen.findByText('Failed to load the channel status')).toBeVisible();
    await userEvent.click(screen.getByRole('button', {name: 'Try again'}));

    expect(await (await currentStatus()).findByText('On track')).toBeVisible();
});

test('drops the load error when a status update arrives', async () => {
    mockFetchChannelStatus.mockRejectedValueOnce(new Error('Failed to load the channel status'));
    const store = renderWithStore(<ChannelStatusSidebar/>);
    expect(await screen.findByText('Failed to load the channel status')).toBeVisible();

    act(() => {
        store.dispatch(statusReceived('channel-id', 'off-track'));
    });

    expect(await (await currentStatus()).findByText('Off track')).toBeVisible();
    expect(screen.queryByText('Failed to load the channel status')).not.toBeInTheDocument();
});

test('keeps a status update that arrives before the load finishes', async () => {
    let resolveFetch: (status: string) => void = () => {};
    mockFetchChannelStatus.mockReturnValueOnce(new Promise((resolve) => {
        resolveFetch = resolve;
    }));
    const store = renderWithStore(<ChannelStatusSidebar/>);

    act(() => {
        store.dispatch(statusReceived('channel-id', 'off-track'));
    });
    await act(async () => resolveFetch('at-risk'));

    expect((await currentStatus()).getByText('Off track')).toBeVisible();
});

test('announces save errors', async () => {
    mockFetchChannelStatus.mockResolvedValueOnce('');
    mockSetChannelStatus.mockRejectedValueOnce(new Error('Failed to update the channel status'));
    renderWithStore(<ChannelStatusSidebar/>);

    await userEvent.click(await screen.findByRole('combobox', {name: 'New status'}));
    await userEvent.click(screen.getByRole('option', {name: 'On track'}));
    await userEvent.click(screen.getByRole('button', {name: 'Save'}));

    expect(await within(screen.getByRole('alert')).findByText('Failed to update the channel status')).toBeVisible();
});

test('ignores a save that fails after switching channels', async () => {
    let rejectSave: (error: Error) => void = () => {};
    mockSetChannelStatus.mockReturnValue(new Promise((_, reject) => {
        rejectSave = reject;
    }));
    mockFetchChannelStatus.mockResolvedValueOnce('');
    mockFetchChannelStatus.mockResolvedValueOnce('at-risk');
    const store = renderWithStore(<ChannelStatusSidebar/>);

    await userEvent.click(await screen.findByRole('combobox', {name: 'New status'}));
    await userEvent.click(screen.getByRole('option', {name: 'On track'}));
    await userEvent.click(screen.getByRole('button', {name: 'Save'}));

    act(() => {
        store.dispatch(switchChannel('other-channel-id'));
    });
    expect(await (await currentStatus()).findByText('At risk')).toBeVisible();

    await act(async () => rejectSave(new Error('Failed to update the channel status')));

    expect(screen.queryByText('Failed to update the channel status')).not.toBeInTheDocument();
    expect((await currentStatus()).getByText('At risk')).toBeVisible();
});
