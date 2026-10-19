// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {findStatusOption} from 'channel_status/statuses';
import {useChannelStatus} from 'channel_status/use_channel_status';
import {useIntl} from 'react-intl';

import {Tag} from '@mattermost/compass-ui/components/tag';
import type {Channel} from '@mattermost/types/channels';

import './channel_status_tag.css';

type Props = {
    channel: Channel;
    onOpen: () => void;
};

// Rendered in the channel header; shows nothing until the channel has a status.
export default function ChannelStatusTag({channel, onOpen}: Props) {
    const intl = useIntl();
    const {status} = useChannelStatus(channel.id);
    const option = findStatusOption(status);
    if (!option) {
        return null;
    }

    const label = intl.formatMessage(option.message);
    return (
        <button
            type='button'
            className='ChannelStatusTag'
            aria-label={intl.formatMessage({id: 'channel_status.tag_label', defaultMessage: 'Channel status: {status}'}, {status: label})}
            onClick={onOpen}
        >
            <Tag
                label={label}
                type={option.tagType}
                size='small'
            />
        </button>
    );
}
