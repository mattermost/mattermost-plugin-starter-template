import React from 'react';

import {Message} from '@mattermost/compass-ui/components/message';
import {MessageActions} from '@mattermost/compass-ui/components/message-actions';
import {MessageHeader} from '@mattermost/compass-ui/components/message-header';
import {MessageInput} from '@mattermost/compass-ui/components/message-input';
import {MessageSeparator} from '@mattermost/compass-ui/components/message-separator';
import {PinnedSavedIndicators} from '@mattermost/compass-ui/components/pinned-saved-indicators';
import {ThreadFooter} from '@mattermost/compass-ui/components/thread-footer';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function MessagePreview() {
    return (
        <Message
            avatarAlt='Ada Lovelace'
            username='ada.lovelace'
            timestamp='2:14 PM'
            showMessageActions={false}
        >
            {'Hello from Compass UI.'}
        </Message>
    );
}

function MessageDetail() {
    return (
        <SpecimenRow>
            <Message
                avatarAlt='Ada Lovelace'
                username='ada.lovelace'
                timestamp='2:14 PM'
                isBot={true}
                botLabel='BOT'
                showMessageActions={false}
            >
                {'A bot message with a header and body.'}
            </Message>
        </SpecimenRow>
    );
}

function MessageHeaderPreview() {
    return (
        <MessageHeader
            username='ada.lovelace'
            timestamp='2:14 PM'
        />
    );
}

function MessageHeaderDetail() {
    return (
        <SpecimenRow>
            <MessageHeader
                username='compass-bot'
                timestamp='Yesterday'
                isBot={true}
                botLabel='BOT'
            />
        </SpecimenRow>
    );
}

function MessageInputPreview() {
    return (
        <MessageInput
            width='narrow'
            placeholder='Reply…'
        />
    );
}

function MessageInputDetail() {
    return (
        <SpecimenRow>
            <MessageInput
                width='narrow'
                placeholder='Write a message'
                showAttachments={true}
            />
        </SpecimenRow>
    );
}

function MessageActionsPreview() {
    return (
        <MessageActions
            type='rhs'
            visible={true}
        />
    );
}

function MessageActionsDetail() {
    return (
        <SpecimenRow>
            <MessageActions
                type='center-channel'
                visible={true}
            />
        </SpecimenRow>
    );
}

function MessageSeparatorPreview() {
    return (
        <MessageSeparator
            type='date'
            label='Today'
        />
    );
}

function MessageSeparatorDetail() {
    return (
        <>
            <SpecimenRow>
                <MessageSeparator
                    type='date'
                    label='Yesterday'
                />
            </SpecimenRow>
            <SpecimenRow>
                <MessageSeparator type='new-messages'/>
            </SpecimenRow>
        </>
    );
}

function ThreadFooterPreview() {
    return (
        <ThreadFooter
            replyCount={4}
            avatars={[
                {key: 'a', name: 'Ada Lovelace'},
                {key: 'b', name: 'Grace Hopper'},
            ]}
        />
    );
}

function ThreadFooterDetail() {
    return (
        <SpecimenRow>
            <ThreadFooter
                replyCount={8}
                badge='mention'
                mentionCount={2}
                following={true}
                lastReplyTime='5 min ago'
                avatars={[
                    {key: 'a', name: 'Ada Lovelace'},
                    {key: 'b', name: 'Grace Hopper'},
                    {key: 'c', name: 'Alan Turing'},
                ]}
            />
        </SpecimenRow>
    );
}

function PinnedSavedPreview() {
    return <PinnedSavedIndicators/>;
}

function PinnedSavedDetail() {
    return (
        <SpecimenRow>
            <PinnedSavedIndicators/>
        </SpecimenRow>
    );
}

export const messagingEntries: CatalogEntry[] = [
    {
        id: 'message',
        name: 'Message',
        category: 'messaging',
        description: 'Post with avatar, header, and body.',
        Preview: MessagePreview,
        Detail: MessageDetail,
    },
    {
        id: 'message-header',
        name: 'Message Header',
        category: 'messaging',
        description: 'Username, optional bot tag, and timestamp.',
        Preview: MessageHeaderPreview,
        Detail: MessageHeaderDetail,
    },
    {
        id: 'message-input',
        name: 'Message Input',
        category: 'messaging',
        description: 'Composer; narrow width matches RHS layouts.',
        Preview: MessageInputPreview,
        Detail: MessageInputDetail,
    },
    {
        id: 'message-actions',
        name: 'Message Actions',
        category: 'messaging',
        description: 'Hover toolbar for react, save, reply, and more.',
        Preview: MessageActionsPreview,
        Detail: MessageActionsDetail,
    },
    {
        id: 'message-separator',
        name: 'Message Separator',
        category: 'messaging',
        description: 'Date or new-messages divider in the stream.',
        Preview: MessageSeparatorPreview,
        Detail: MessageSeparatorDetail,
    },
    {
        id: 'thread-footer',
        name: 'Thread Footer',
        category: 'messaging',
        description: 'Reply count, participants, and follow state.',
        Preview: ThreadFooterPreview,
        Detail: ThreadFooterDetail,
    },
    {
        id: 'pinned-saved-indicators',
        name: 'Pinned Saved Indicators',
        category: 'messaging',
        description: 'Pinned and saved row above a message.',
        Preview: PinnedSavedPreview,
        Detail: PinnedSavedDetail,
    },
];
