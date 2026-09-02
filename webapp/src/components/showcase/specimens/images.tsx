import React from 'react';

import PaletteOutlineIcon from '@mattermost/compass-icons/components/palette-outline';
import {Emoji} from '@mattermost/compass-ui/components/emoji';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {TeamAvatar} from '@mattermost/compass-ui/components/team-avatar';
import {UserAvatar} from '@mattermost/compass-ui/components/user-avatar';
import {UserAvatarGroup} from '@mattermost/compass-ui/components/user-avatar-group';

import {PLACEHOLDER_AVATAR, SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function IconPreview() {
    return (
        <Icon
            glyph={<PaletteOutlineIcon/>}
            size='24'
        />
    );
}

function IconDetail() {
    return (
        <SpecimenRow label='Size'>
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='16'
            />
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='20'
            />
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='24'
            />
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='32'
            />
        </SpecimenRow>
    );
}

function EmojiPreview() {
    return (
        <Emoji
            emoji='🎉'
            size='24'
        />
    );
}

function EmojiDetail() {
    return (
        <SpecimenRow>
            <Emoji
                emoji='👍'
                size='16'
            />
            <Emoji
                emoji='🎉'
                size='24'
            />
            <Emoji
                emoji='🚀'
                size='32'
            />
        </SpecimenRow>
    );
}

function UserAvatarPreview() {
    return (
        <UserAvatar
            alt='Ada Lovelace'
            name='Ada Lovelace'
            size='32'
            status={true}
        />
    );
}

function UserAvatarDetail() {
    return (
        <SpecimenRow>
            <UserAvatar
                alt='Ada Lovelace'
                name='Ada Lovelace'
                size='24'
            />
            <UserAvatar
                alt='Ada Lovelace'
                name='Ada Lovelace'
                size='32'
                status={true}
            />
            <UserAvatar
                alt='Ada Lovelace'
                name='Ada Lovelace'
                size='48'
                src={PLACEHOLDER_AVATAR}
            />
        </SpecimenRow>
    );
}

function TeamAvatarPreview() {
    return (
        <TeamAvatar
            initials='UX'
            size='32'
        />
    );
}

function TeamAvatarDetail() {
    return (
        <SpecimenRow>
            <TeamAvatar
                initials='UX'
                size='24'
            />
            <TeamAvatar
                initials='EN'
                size='32'
                badge={3}
            />
            <TeamAvatar
                initials='DS'
                size='48'
                state='active'
            />
        </SpecimenRow>
    );
}

function UserAvatarGroupPreview() {
    return (
        <UserAvatarGroup
            avatars={[
                {key: 'a', name: 'Ada Lovelace'},
                {key: 'b', name: 'Grace Hopper'},
                {key: 'c', name: 'Alan Turing'},
            ]}
        />
    );
}

function UserAvatarGroupDetail() {
    return (
        <SpecimenRow>
            <UserAvatarGroup
                size='24'
                max={3}
                avatars={[
                    {key: 'a', name: 'Ada Lovelace'},
                    {key: 'b', name: 'Grace Hopper'},
                    {key: 'c', name: 'Alan Turing'},
                    {key: 'd', name: 'Katherine Johnson'},
                ]}
            />
        </SpecimenRow>
    );
}

export const imageEntries: CatalogEntry[] = [
    {
        id: 'icon',
        name: 'Icon',
        category: 'images',
        description: 'Standard wrapper for compass-icons glyphs.',
        Preview: IconPreview,
        Detail: IconDetail,
    },
    {
        id: 'emoji',
        name: 'Emoji',
        category: 'images',
        description: 'Sized Unicode emoji aligned to the icon scale.',
        Preview: EmojiPreview,
        Detail: EmojiDetail,
    },
    {
        id: 'user-avatar',
        name: 'User Avatar',
        category: 'images',
        description: 'Circular person avatar with initials fallback and status.',
        Preview: UserAvatarPreview,
        Detail: UserAvatarDetail,
    },
    {
        id: 'team-avatar',
        name: 'Team Avatar',
        category: 'images',
        description: 'Rounded-square team mark, distinct from user avatars.',
        Preview: TeamAvatarPreview,
        Detail: TeamAvatarDetail,
    },
    {
        id: 'user-avatar-group',
        name: 'User Avatar Group',
        category: 'images',
        description: 'Stacked participants with overflow count.',
        Preview: UserAvatarGroupPreview,
        Detail: UserAvatarGroupDetail,
    },
];
