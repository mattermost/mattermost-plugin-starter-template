import React, {useState} from 'react';

import {Divider} from '@mattermost/compass-ui/components/divider';
import {EmptyState} from '@mattermost/compass-ui/components/empty-state';
import {Scrollbar} from '@mattermost/compass-ui/components/scrollbar';
import {Tabs} from '@mattermost/compass-ui/components/tabs';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function DividerPreview() {
    return <Divider/>;
}

function DividerDetail() {
    return (
        <SpecimenRow>
            <div style={{width: '100%'}}>
                <p>{'Above'}</p>
                <Divider/>
                <p>{'Below'}</p>
            </div>
        </SpecimenRow>
    );
}

function EmptyStatePreview() {
    return (
        <EmptyState
            title='Nothing here'
            description='Try another filter.'
        />
    );
}

function EmptyStateDetail() {
    return (
        <SpecimenRow>
            <EmptyState
                title='No saved messages'
                description='Save a message to find it later.'
                action={{children: 'Browse channels', emphasis: 'secondary'}}
            />
        </SpecimenRow>
    );
}

function ScrollbarPreview() {
    return (
        <Scrollbar style={{height: 72, width: '100%'}}>
            <p>{'Scroll to see more content in this region. Compass uses a thin thumb over a transparent track.'}</p>
            <p>{'Second paragraph to force overflow.'}</p>
            <p>{'Third paragraph to force overflow.'}</p>
        </Scrollbar>
    );
}

function ScrollbarDetail() {
    return (
        <SpecimenRow>
            <Scrollbar
                alwaysVisible={true}
                style={{height: 96, width: '100%'}}
            >
                <p>{'Always-visible thumb when content overflows.'}</p>
                <p>{'Line two.'}</p>
                <p>{'Line three.'}</p>
                <p>{'Line four.'}</p>
            </Scrollbar>
        </SpecimenRow>
    );
}

function TabsPreview() {
    const [active, setActive] = useState('one');
    return (
        <Tabs
            activeKey={active}
            onChange={setActive}
            tabs={[
                {key: 'one', label: 'One', countBadge: 2},
                {key: 'two', label: 'Two'},
            ]}
        />
    );
}

function TabsDetail() {
    const [active, setActive] = useState('mentions');
    return (
        <SpecimenRow>
            <Tabs
                activeKey={active}
                onChange={setActive}
                tabs={[
                    {key: 'mentions', label: 'Mentions', countBadge: 4},
                    {key: 'threads', label: 'Threads', unreadBadge: true},
                    {key: 'saved', label: 'Saved'},
                ]}
            />
        </SpecimenRow>
    );
}

export const layoutEntries: CatalogEntry[] = [
    {
        id: 'divider',
        name: 'Divider',
        category: 'layout',
        description: 'Thin horizontal rule between related groups.',
        Preview: DividerPreview,
        Detail: DividerDetail,
    },
    {
        id: 'empty-state',
        name: 'Empty State',
        category: 'layout',
        description: 'Placeholder when a view has no content yet.',
        Preview: EmptyStatePreview,
        Detail: EmptyStateDetail,
    },
    {
        id: 'scrollbar',
        name: 'Scrollbar',
        category: 'layout',
        description: 'Minimal overflow indicator for scroll regions.',
        Preview: ScrollbarPreview,
        Detail: ScrollbarDetail,
    },
    {
        id: 'tabs',
        name: 'Tabs',
        category: 'layout',
        description: 'Switch sibling views; used here as a category filter.',
        Preview: TabsPreview,
        Detail: TabsDetail,
    },
];
