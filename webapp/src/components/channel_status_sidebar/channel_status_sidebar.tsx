// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {statusReceived} from 'channel_status/state';
import {canEditChannelStatus, findStatusOption, STATUS_OPTIONS} from 'channel_status/statuses';
import {useChannelStatus} from 'channel_status/use_channel_status';
import {clearChannelStatus, setChannelStatus} from 'client';
import {useId, useState} from 'react';
import {FormattedMessage, useIntl} from 'react-intl';
import {useDispatch, useSelector} from 'react-redux';

import {Button} from '@mattermost/compass-ui/components/button';
import {SectionNotice} from '@mattermost/compass-ui/components/section-notice';
import {Select} from '@mattermost/compass-ui/components/select';
import {Tag} from '@mattermost/compass-ui/components/tag';
import type {Channel} from '@mattermost/types/channels';
import type {GlobalState} from '@mattermost/types/store';

import {getCurrentChannel} from 'mattermost-redux/selectors/entities/channels';

import './channel_status_sidebar.css';

export default function ChannelStatusSidebar() {
    const channel = useSelector(getCurrentChannel);
    if (!channel) {
        return null;
    }

    // A new key per channel resets the form, so a save that finishes after the user
    // switches channels can't show its result in the wrong one.
    return (
        <ChannelStatusPanel
            key={channel.id}
            channel={channel}
        />
    );
}

function ChannelStatusPanel({channel}: {channel: Channel}) {
    const intl = useIntl();
    const dispatch = useDispatch();
    const headingId = useId();
    const {status, error, retry} = useChannelStatus(channel.id);
    const canEdit = useSelector((state: GlobalState) => canEditChannelStatus(state, channel));
    const [selected, setSelected] = useState<string>();
    const [saving, setSaving] = useState(false);
    const [saveError, setSaveError] = useState<string>();

    if (error) {
        return (
            <div className='ChannelStatusSidebar'>
                <SectionNotice
                    type='danger'
                    title={intl.formatMessage({id: 'channel_status.load_error', defaultMessage: 'Couldn’t load the channel status'})}
                    description={error}
                    primaryButtonLabel={intl.formatMessage({id: 'channel_status.retry', defaultMessage: 'Try again'})}
                    onPrimaryAction={retry}
                />
            </div>
        );
    }

    if (status === undefined) {
        return (
            <div className='ChannelStatusSidebar'>
                <p className='ChannelStatusSidebar__text'>
                    <FormattedMessage
                        id='channel_status.loading'
                        defaultMessage='Loading…'
                    />
                </p>
            </div>
        );
    }

    const current = findStatusOption(status);
    const save = (next: string) => {
        setSaving(true);
        setSaveError(undefined);
        (next ? setChannelStatus(channel.id, next) : clearChannelStatus(channel.id)).then(
            () => {
                dispatch(statusReceived(channel.id, next));
                setSelected(undefined);
                setSaving(false);
            },
            (err: Error) => {
                setSaveError(err.message);
                setSaving(false);
            },
        );
    };

    return (
        <div className='ChannelStatusSidebar'>
            <p className='ChannelStatusSidebar__text'>
                <FormattedMessage
                    id='channel_status.intro'
                    defaultMessage='Show everyone in {channel} how things are going. You can also use {command}.'
                    values={{
                        channel: channel.display_name,
                        command: <code key='command'>{'/channelstatus'}</code>,
                    }}
                />
            </p>

            <section
                className='ChannelStatusSidebar__current'
                aria-labelledby={headingId}
                aria-live='polite'
            >
                <h3
                    id={headingId}
                    className='ChannelStatusSidebar__heading'
                >
                    <FormattedMessage
                        id='channel_status.current'
                        defaultMessage='Current status'
                    />
                </h3>
                {current ? (
                    <Tag
                        label={intl.formatMessage(current.message)}
                        type={current.tagType}
                        size='small'
                    />
                ) : (
                    <span className='ChannelStatusSidebar__text'>
                        <FormattedMessage
                            id='channel_status.none'
                            defaultMessage='No status set'
                        />
                    </span>
                )}
            </section>

            {canEdit ? (
                <div className='ChannelStatusSidebar__form'>
                    <Select
                        label={intl.formatMessage({id: 'channel_status.select_label', defaultMessage: 'New status'})}
                        placeholder={intl.formatMessage({id: 'channel_status.select_placeholder', defaultMessage: 'Choose a status'})}
                        options={STATUS_OPTIONS.map((option) => ({
                            value: option.value,
                            label: intl.formatMessage(option.message),
                            leadingVisual: <span className={`ChannelStatusSidebar__dot ChannelStatusSidebar__dot--${option.value}`}/>,
                        }))}
                        value={selected ?? status}
                        onChange={setSelected}
                        disabled={saving}
                    />
                    <div className='ChannelStatusSidebar__actions'>
                        <Button
                            emphasis='primary'
                            size='small'
                            loading={saving}
                            disabled={!selected || selected === status}
                            onClick={() => selected && save(selected)}
                        >
                            <FormattedMessage
                                id='channel_status.save'
                                defaultMessage='Save'
                            />
                        </Button>
                        {current && (
                            <Button
                                emphasis='tertiary'
                                size='small'
                                disabled={saving}
                                onClick={() => save('')}
                            >
                                <FormattedMessage
                                    id='channel_status.clear'
                                    defaultMessage='Clear status'
                                />
                            </Button>
                        )}
                    </div>
                    {saveError && (
                        <SectionNotice
                            type='danger'
                            title={intl.formatMessage({id: 'channel_status.save_error', defaultMessage: 'Couldn’t update the status'})}
                            description={saveError}
                        />
                    )}
                </div>
            ) : (
                <p className='ChannelStatusSidebar__text'>
                    <FormattedMessage
                        id='channel_status.read_only'
                        defaultMessage='You don’t have permission to change this channel’s status.'
                    />
                </p>
            )}
        </div>
    );
}
