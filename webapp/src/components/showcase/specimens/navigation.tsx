import React from 'react';

import PaletteOutlineIcon from '@mattermost/compass-icons/components/palette-outline';
import {AppBarItem} from '@mattermost/compass-ui/components/app-bar-item';
import {ChannelSidebarItem} from '@mattermost/compass-ui/components/channel-sidebar-item';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {MenuItem} from '@mattermost/compass-ui/components/menu-item';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function ChannelSidebarItemPreview() {
    return (
        <ChannelSidebarItem
            name='design'
            status='mention'
            mentionCount={3}
        />
    );
}

function ChannelSidebarItemDetail() {
    return (
        <>
            <SpecimenRow>
                <ChannelSidebarItem
                    name='town-square'
                    active={true}
                />
            </SpecimenRow>
            <SpecimenRow>
                <ChannelSidebarItem
                    name='off-topic'
                    status='unread'
                />
            </SpecimenRow>
            <SpecimenRow>
                <ChannelSidebarItem
                    name='Ada Lovelace'
                    leadingVisual='direct-message'
                    showAvatarStatus={true}
                />
            </SpecimenRow>
        </>
    );
}

function MenuItemPreview() {
    return <MenuItem label='Edit channel'/>;
}

function MenuItemDetail() {
    return (
        <SpecimenRow>
            <div style={{width: '100%'}}>
                <MenuItem label='Copy link'/>
                <MenuItem
                    label='Mute channel'
                    mentionCount={2}
                />
                <MenuItem
                    label='Leave channel'
                    destructive={true}
                />
            </div>
        </SpecimenRow>
    );
}

function AppBarItemPreview() {
    return (
        <AppBarItem
            icon={(
                <Icon
                    glyph={<PaletteOutlineIcon/>}
                    size='20'
                />
            )}
            label='Compass UI'
            mentionBadge={2}
        />
    );
}

function AppBarItemDetail() {
    return (
        <SpecimenRow>
            <AppBarItem
                icon={(
                    <Icon
                        glyph={<PaletteOutlineIcon/>}
                        size='20'
                    />
                )}
                label='Default'
            />
            <AppBarItem
                icon={(
                    <Icon
                        glyph={<PaletteOutlineIcon/>}
                        size='20'
                    />
                )}
                label='Selected'
                state='selected'
            />
            <AppBarItem
                icon={(
                    <Icon
                        glyph={<PaletteOutlineIcon/>}
                        size='20'
                    />
                )}
                label='Unread'
                unreadBadge={true}
            />
        </SpecimenRow>
    );
}

export const navigationEntries: CatalogEntry[] = [
    {
        id: 'channel-sidebar-item',
        name: 'Channel Sidebar Item',
        category: 'navigation',
        description: 'Channel, DM, or GM row with unread and mention state.',
        Preview: ChannelSidebarItemPreview,
        Detail: ChannelSidebarItemDetail,
    },
    {
        id: 'menu-item',
        name: 'Menu Item',
        category: 'navigation',
        description: 'Selectable row inside a menu or dropdown.',
        Preview: MenuItemPreview,
        Detail: MenuItemDetail,
    },
    {
        id: 'app-bar-item',
        name: 'App Bar Item',
        category: 'navigation',
        description: 'Plugin icon in the right-side app bar.',
        Preview: AppBarItemPreview,
        Detail: AppBarItemDetail,
    },
];
