import React from 'react';

import {MentionBadge} from '@mattermost/compass-ui/components/mention-badge';
import {ReactionPill} from '@mattermost/compass-ui/components/reaction-pill';
import {ShortcutTag, ShortcutTagGroup} from '@mattermost/compass-ui/components/shortcut-tag';
import {StatusBadge} from '@mattermost/compass-ui/components/status-badge';
import {Tag} from '@mattermost/compass-ui/components/tag';
import {UnreadBadge} from '@mattermost/compass-ui/components/unread-badge';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function TagPreview() {
    return (
        <Tag
            label='BETA'
            type='info'
        />
    );
}

function TagDetail() {
    return (
        <SpecimenRow>
            <Tag label='Default'/>
            <Tag
                label='Info'
                type='info'
            />
            <Tag
                label='Success'
                type='success'
            />
            <Tag
                label='Warning'
                type='warning'
            />
            <Tag
                label='Danger'
                type='danger'
            />
        </SpecimenRow>
    );
}

function ShortcutTagPreview() {
    return <ShortcutTag label='K'/>;
}

function ShortcutTagDetail() {
    return (
        <SpecimenRow>
            <ShortcutTagGroup labels={['Ctrl', 'K']}/>
        </SpecimenRow>
    );
}

function MentionBadgePreview() {
    return (
        <MentionBadge
            count={3}
            location='channel'
            size='medium'
        />
    );
}

function MentionBadgeDetail() {
    return (
        <SpecimenRow>
            <MentionBadge
                count={1}
                location='channel'
            />
            <MentionBadge
                count={12}
                location='sidebar'
            />
            <MentionBadge
                count={120}
                location='menu-item'
            />
        </SpecimenRow>
    );
}

function StatusBadgePreview() {
    return (
        <StatusBadge
            status='online'
            size='medium'
        />
    );
}

function StatusBadgeDetail() {
    return (
        <SpecimenRow>
            <StatusBadge
                status='online'
                size='medium'
            />
            <StatusBadge
                status='away'
                size='medium'
            />
            <StatusBadge
                status='do-not-disturb'
                size='medium'
            />
            <StatusBadge
                status='offline'
                size='medium'
            />
        </SpecimenRow>
    );
}

function UnreadBadgePreview() {
    return <UnreadBadge/>;
}

function UnreadBadgeDetail() {
    return (
        <SpecimenRow>
            <UnreadBadge size='6'/>
            <UnreadBadge size='8'/>
        </SpecimenRow>
    );
}

function ReactionPillPreview() {
    return (
        <ReactionPill
            emoji='👍'
            label='Ada'
        />
    );
}

function ReactionPillDetail() {
    return (
        <SpecimenRow>
            <ReactionPill
                type='reaction'
                emoji='🎉'
                label='Grace'
            />
            <ReactionPill
                type='hand-raise'
                label='Alan'
            />
        </SpecimenRow>
    );
}

export const statusEntries: CatalogEntry[] = [
    {
        id: 'tag',
        name: 'Tag',
        category: 'status',
        description: 'Compact pill for role, tier, or status metadata.',
        Preview: TagPreview,
        Detail: TagDetail,
    },
    {
        id: 'shortcut-tag',
        name: 'Shortcut Tag',
        category: 'status',
        description: 'Keyboard key chip for menus and tooltips.',
        Preview: ShortcutTagPreview,
        Detail: ShortcutTagDetail,
    },
    {
        id: 'mention-badge',
        name: 'Mention Badge',
        category: 'status',
        description: 'Numeric mention count for sidebars and lists.',
        Preview: MentionBadgePreview,
        Detail: MentionBadgeDetail,
    },
    {
        id: 'status-badge',
        name: 'Status Badge',
        category: 'status',
        description: 'Presence: online, away, DND, or offline.',
        Preview: StatusBadgePreview,
        Detail: StatusBadgeDetail,
    },
    {
        id: 'unread-badge',
        name: 'Unread Badge',
        category: 'status',
        description: 'Small unread dot when a count would be noisy.',
        Preview: UnreadBadgePreview,
        Detail: UnreadBadgeDetail,
    },
    {
        id: 'reaction-pill',
        name: 'Reaction Pill',
        category: 'status',
        description: 'Transient call reaction or hand-raise signal.',
        Preview: ReactionPillPreview,
        Detail: ReactionPillDetail,
    },
];
