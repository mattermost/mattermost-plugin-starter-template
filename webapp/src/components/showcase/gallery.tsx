import React, {useMemo} from 'react';

import {Divider} from '@mattermost/compass-ui/components/divider';
import {EmptyState} from '@mattermost/compass-ui/components/empty-state';
import {Scrollbar} from '@mattermost/compass-ui/components/scrollbar';
import {SearchInput} from '@mattermost/compass-ui/components/search-input';
import {Tabs} from '@mattermost/compass-ui/components/tabs';

import {catalog} from './catalog';
import {countByCategory, filterCatalog} from './helpers';
import {CATEGORIES, type CategoryId} from './types';

type GalleryProps = {
    query: string;
    onQueryChange: (value: string) => void;
    category: CategoryId;
    onCategoryChange: (value: CategoryId) => void;
    onSelect: (id: string) => void;
};

export default function Gallery({
    query,
    onQueryChange,
    category,
    onCategoryChange,
    onSelect,
}: GalleryProps) {
    const matches = useMemo(
        () => filterCatalog(catalog, query, category),
        [query, category],
    );

    const tabs = useMemo(
        () => CATEGORIES.map((item) => ({
            key: item.id,
            label: item.label,
            countBadge: countByCategory(catalog, query, item.id),
        })),
        [query],
    );

    const sections = useMemo(() => {
        const grouped: Array<{id: CategoryId; label: string; entries: typeof matches}> = [];
        for (const item of CATEGORIES) {
            if (item.id === 'all') {
                continue;
            }
            const entries = matches.filter((entry) => entry.category === item.id);
            if (entries.length) {
                grouped.push({id: item.id, label: item.label, entries});
            }
        }
        return grouped;
    }, [matches]);

    const resetFilters = () => {
        onQueryChange('');
        onCategoryChange('all');
    };

    return (
        <div className='CompassShowcase'>
            <div className='CompassShowcase__filters'>
                <SearchInput
                    label='Search components'
                    value={query}
                    onChange={(event) => onQueryChange(event.target.value)}
                    onClear={() => onQueryChange('')}
                />
                <div className='CompassShowcase__tabs'>
                    <Tabs
                        tabs={tabs}
                        activeKey={category}
                        onChange={(key) => onCategoryChange(key as CategoryId)}
                    />
                </div>
            </div>
            <Scrollbar className='CompassShowcase__scroll'>
                {matches.length === 0 ? (
                    <EmptyState
                        title='No matching components'
                        description='Try a different search or category.'
                        action={{
                            children: 'Clear filters',
                            emphasis: 'secondary',
                            onClick: resetFilters,
                        }}
                    />
                ) : (
                    sections.map((section, index) => (
                        <section
                            key={section.id}
                            className='CompassShowcase__section'
                        >
                            {index > 0 && <Divider/>}
                            <h2 className='CompassShowcase__section-title'>
                                {section.label}
                            </h2>
                            <div className='CompassShowcase__cards'>
                                {section.entries.map((entry) => {
                                    const Preview = entry.Preview;
                                    return (
                                        <button
                                            key={entry.id}
                                            type='button'
                                            className='CompassShowcase__card'
                                            onClick={() => onSelect(entry.id)}
                                        >
                                            <div className='CompassShowcase__card-preview'>
                                                <Preview/>
                                            </div>
                                            <div className='CompassShowcase__card-meta'>
                                                <span className='CompassShowcase__card-name'>
                                                    {entry.name}
                                                </span>
                                                {entry.layoutExcerpt && (
                                                    <span className='CompassShowcase__card-flag'>
                                                        {'Excerpt'}
                                                    </span>
                                                )}
                                            </div>
                                        </button>
                                    );
                                })}
                            </div>
                        </section>
                    ))
                )}
            </Scrollbar>
        </div>
    );
}
