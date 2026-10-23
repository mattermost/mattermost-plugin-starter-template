// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {defineMessages, type MessageDescriptor} from 'react-intl';

import type {TagType} from '@mattermost/compass-ui/components/tag';
import type {Channel} from '@mattermost/types/channels';
import type {GlobalState} from '@mattermost/types/store';

import {Permissions} from 'mattermost-redux/constants';
import {haveIChannelPermission} from 'mattermost-redux/selectors/entities/roles';

export type StatusOption = {
    value: string;
    tagType: TagType;
    message: MessageDescriptor;
};

const messages = defineMessages({
    onTrack: {id: 'channel_status.on_track', defaultMessage: 'On track'},
    atRisk: {id: 'channel_status.at_risk', defaultMessage: 'At risk'},
    offTrack: {id: 'channel_status.off_track', defaultMessage: 'Off track'},
});

// Must match the values the server accepts in server/channelstatus.
export const STATUS_OPTIONS: StatusOption[] = [
    {value: 'on-track', tagType: 'success', message: messages.onTrack},
    {value: 'at-risk', tagType: 'warning', message: messages.atRisk},
    {value: 'off-track', tagType: 'danger', message: messages.offTrack},
];

export function findStatusOption(value: string | undefined) {
    return STATUS_OPTIONS.find((option) => option.value === value);
}

// Mirrors the server's check so the UI can hide controls; the server still enforces it.
export function canEditChannelStatus(state: GlobalState, channel: Channel) {
    switch (channel.type) {
    case 'O':
        return haveIChannelPermission(state, channel.team_id, channel.id, Permissions.MANAGE_PUBLIC_CHANNEL_PROPERTIES);
    case 'P':
        return haveIChannelPermission(state, channel.team_id, channel.id, Permissions.MANAGE_PRIVATE_CHANNEL_PROPERTIES);
    default:
        return true;
    }
}
