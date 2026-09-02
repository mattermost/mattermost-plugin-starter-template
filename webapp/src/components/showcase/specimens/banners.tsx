import React from 'react';

import {GlobalBanner} from '@mattermost/compass-ui/components/global-banner';
import {MoreUnreadsBanner} from '@mattermost/compass-ui/components/more-unreads-banner';
import {NewMessageBanner} from '@mattermost/compass-ui/components/new-message-banner';
import {SearchTipBanner} from '@mattermost/compass-ui/components/search-tip-banner';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function GlobalBannerPreview() {
    return (
        <GlobalBanner
            type='info'
            message='Scheduled maintenance tonight'
        />
    );
}

function GlobalBannerDetail() {
    return (
        <>
            <SpecimenRow>
                <GlobalBanner
                    type='info'
                    message='Update available'
                />
            </SpecimenRow>
            <SpecimenRow>
                <GlobalBanner
                    type='warning'
                    message='Connection is unstable'
                />
            </SpecimenRow>
            <SpecimenRow>
                <GlobalBanner
                    type='danger'
                    message='License expires in 3 days'
                />
            </SpecimenRow>
        </>
    );
}

function NewMessageBannerPreview() {
    return (
        <NewMessageBanner
            countLabel='12 new messages'
        />
    );
}

function NewMessageBannerDetail() {
    return (
        <>
            <SpecimenRow>
                <NewMessageBanner countLabel='21 new messages since Saturday'/>
            </SpecimenRow>
            <SpecimenRow>
                <NewMessageBanner type='new-replies'/>
            </SpecimenRow>
        </>
    );
}

function MoreUnreadsBannerPreview() {
    return <MoreUnreadsBanner direction='up'/>;
}

function MoreUnreadsBannerDetail() {
    return (
        <SpecimenRow>
            <MoreUnreadsBanner direction='up'/>
            <MoreUnreadsBanner direction='down'/>
        </SpecimenRow>
    );
}

function SearchTipBannerPreview() {
    return <SearchTipBanner/>;
}

function SearchTipBannerDetail() {
    return (
        <SpecimenRow>
            <SearchTipBanner
                shortcutKeys={[{label: 'Ctrl'}, {label: 'F'}]}
                onDismiss={() => undefined}
            />
        </SpecimenRow>
    );
}

export const bannerEntries: CatalogEntry[] = [
    {
        id: 'global-banner',
        name: 'Global Banner',
        category: 'banners',
        description: 'Full-width workspace announcement.',
        Preview: GlobalBannerPreview,
        Detail: GlobalBannerDetail,
    },
    {
        id: 'new-message-banner',
        name: 'New Message Banner',
        category: 'banners',
        description: 'Jump to unread activity in the message stream.',
        Preview: NewMessageBannerPreview,
        Detail: NewMessageBannerDetail,
    },
    {
        id: 'more-unreads-banner',
        name: 'More Unreads Banner',
        category: 'banners',
        description: 'Sidebar pill pointing to unread channels off-screen.',
        Preview: MoreUnreadsBannerPreview,
        Detail: MoreUnreadsBannerDetail,
    },
    {
        id: 'search-tip-banner',
        name: 'Search Tip Banner',
        category: 'banners',
        description: 'Keyboard shortcut nudge for in-channel search.',
        Preview: SearchTipBannerPreview,
        Detail: SearchTipBannerDetail,
    },
];
