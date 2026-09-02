import React from 'react';

import {AdminConsoleHeader} from '@mattermost/compass-ui/components/admin-console-header';
import {AdminPanel} from '@mattermost/compass-ui/components/admin-panel';
import {ChannelHeader} from '@mattermost/compass-ui/components/channel-header';
import {FeatureDiscoveryPanel} from '@mattermost/compass-ui/components/feature-discovery-panel';
import {GlobalHeader} from '@mattermost/compass-ui/components/global-header';
import {RightSidebarHeader} from '@mattermost/compass-ui/components/right-sidebar';

import {LayoutExcerpt, SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function ChannelHeaderPreview() {
    return (
        <LayoutExcerpt>
            <ChannelHeader
                name='design'
                memberCount={12}
                pinnedCount={2}
            />
        </LayoutExcerpt>
    );
}

function ChannelHeaderDetail() {
    return (
        <SpecimenRow>
            <LayoutExcerpt>
                <ChannelHeader
                    type='channel'
                    name='ux-review'
                    description='Weekly design critique'
                    memberCount={8}
                    favorited={true}
                />
            </LayoutExcerpt>
        </SpecimenRow>
    );
}

function GlobalHeaderPreview() {
    return (
        <LayoutExcerpt>
            <GlobalHeader
                product='channels'
                userAvatarAlt='You'
            />
        </LayoutExcerpt>
    );
}

function GlobalHeaderDetail() {
    return (
        <SpecimenRow>
            <LayoutExcerpt>
                <GlobalHeader
                    product='channels'
                    showChannelsBranding={true}
                    userAvatarAlt='You'
                />
            </LayoutExcerpt>
        </SpecimenRow>
    );
}

function AdminConsoleHeaderPreview() {
    return (
        <LayoutExcerpt>
            <AdminConsoleHeader title='Users'/>
        </LayoutExcerpt>
    );
}

function AdminConsoleHeaderDetail() {
    return (
        <SpecimenRow>
            <LayoutExcerpt>
                <AdminConsoleHeader
                    title='Authentication'
                    showBack={true}
                    enterpriseBadge={true}
                />
            </LayoutExcerpt>
        </SpecimenRow>
    );
}

function AdminPanelPreview() {
    return (
        <LayoutExcerpt>
            <AdminPanel
                title='Notifications'
                subtitle='Choose how you are notified'
            />
        </LayoutExcerpt>
    );
}

function AdminPanelDetail() {
    return (
        <SpecimenRow>
            <LayoutExcerpt>
                <AdminPanel
                    title='Session lengths'
                    subtitle='Control how long users stay signed in'
                    showSwitch={true}
                    switchLabel='Enable'
                    defaultSwitchChecked={true}
                >
                    {'Panel body for grouped admin fields.'}
                </AdminPanel>
            </LayoutExcerpt>
        </SpecimenRow>
    );
}

function RightSidebarHeaderPreview() {
    return (
        <LayoutExcerpt>
            <RightSidebarHeader title='Thread'/>
        </LayoutExcerpt>
    );
}

function RightSidebarHeaderDetail() {
    return (
        <SpecimenRow>
            <LayoutExcerpt>
                <RightSidebarHeader
                    title='Thread'
                    secondaryTitle='design'
                    tag='BETA'
                    onBack={() => undefined}
                    onClose={() => undefined}
                />
            </LayoutExcerpt>
        </SpecimenRow>
    );
}

function FeatureDiscoveryPreview() {
    return (
        <LayoutExcerpt>
            <FeatureDiscoveryPanel
                skuLabel='PROFESSIONAL'
                title='Guest accounts'
                description='Invite guests into specific channels.'
            />
        </LayoutExcerpt>
    );
}

function FeatureDiscoveryDetail() {
    return (
        <SpecimenRow>
            <LayoutExcerpt>
                <FeatureDiscoveryPanel
                    skuLabel='ENTERPRISE'
                    title='Data retention'
                    description='Automatically delete messages after a policy window.'
                    primaryAction={{children: 'Learn more'}}
                />
            </LayoutExcerpt>
        </SpecimenRow>
    );
}

export const chromeEntries: CatalogEntry[] = [
    {
        id: 'channel-header',
        name: 'Channel Header',
        category: 'chrome',
        description: 'Channel title, stats, and header actions.',
        layoutExcerpt: true,
        Preview: ChannelHeaderPreview,
        Detail: ChannelHeaderDetail,
    },
    {
        id: 'global-header',
        name: 'Global Header',
        category: 'chrome',
        description: 'Top bar with product switcher and session controls.',
        layoutExcerpt: true,
        Preview: GlobalHeaderPreview,
        Detail: GlobalHeaderDetail,
    },
    {
        id: 'admin-console-header',
        name: 'Admin Console Header',
        category: 'chrome',
        description: 'System Console page title stripe.',
        layoutExcerpt: true,
        Preview: AdminConsoleHeaderPreview,
        Detail: AdminConsoleHeaderDetail,
    },
    {
        id: 'admin-panel',
        name: 'Admin Panel',
        category: 'chrome',
        description: 'Grouped System Console settings panel.',
        layoutExcerpt: true,
        Preview: AdminPanelPreview,
        Detail: AdminPanelDetail,
    },
    {
        id: 'right-sidebar-header',
        name: 'Right Sidebar Header',
        category: 'chrome',
        description: 'Compass RHS header chrome — not the Mattermost plugin frame.',
        layoutExcerpt: true,
        Preview: RightSidebarHeaderPreview,
        Detail: RightSidebarHeaderDetail,
    },
    {
        id: 'feature-discovery-panel',
        name: 'Feature Discovery Panel',
        category: 'chrome',
        description: 'System Console upsell for a higher-tier SKU.',
        layoutExcerpt: true,
        Preview: FeatureDiscoveryPreview,
        Detail: FeatureDiscoveryDetail,
    },
];
