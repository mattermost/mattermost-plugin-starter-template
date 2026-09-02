import React, {useState} from 'react';

import {catalog} from './catalog';
import Detail from './detail';
import Gallery from './gallery';
import type {CategoryId} from './types';

import './showcase.scss';

export default function ShowcaseRHS() {
    const [selectedId, setSelectedId] = useState<string | null>(null);
    const [query, setQuery] = useState('');
    const [category, setCategory] = useState<CategoryId>('all');

    const selected = catalog.find((entry) => entry.id === selectedId);

    if (selected) {
        return (
            <Detail
                entry={selected}
                onBack={() => setSelectedId(null)}
            />
        );
    }

    return (
        <Gallery
            query={query}
            onQueryChange={setQuery}
            category={category}
            onCategoryChange={setCategory}
            onSelect={setSelectedId}
        />
    );
}
