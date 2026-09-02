import React from 'react';

import type {CatalogEntry, CategoryId} from './types';

export const PLACEHOLDER_AVATAR = 'data:image/svg+xml,' + encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48"><rect fill="#1c58d9" width="48" height="48"/><text x="50%" y="54%" text-anchor="middle" fill="#fff" font-size="20" font-family="sans-serif">AB</text></svg>',
);

type SpecimenRowProps = {
    label?: string;
    children: React.ReactNode;
};

export function SpecimenRow({label, children}: SpecimenRowProps) {
    return (
        <div className='CompassShowcase__row'>
            {label && (
                <div className='CompassShowcase__row-label'>
                    {label}
                </div>
            )}
            <div className='CompassShowcase__row-body'>
                {children}
            </div>
        </div>
    );
}

type LayoutExcerptProps = {
    children: React.ReactNode;
};

export function LayoutExcerpt({children}: LayoutExcerptProps) {
    return (
        <div className='CompassShowcase__excerpt'>
            <p className='CompassShowcase__excerpt-note'>
                {'Too wide for the RHS; showing a compact excerpt.'}
            </p>
            <div className='CompassShowcase__excerpt-preview'>
                {children}
            </div>
        </div>
    );
}

export function filterCatalog(
    entries: CatalogEntry[],
    query: string,
    category: CategoryId,
): CatalogEntry[] {
    const normalized = query.trim().toLowerCase();

    return entries.filter((entry) => {
        if (category !== 'all' && entry.category !== category) {
            return false;
        }
        if (!normalized) {
            return true;
        }

        return entry.name.toLowerCase().includes(normalized) ||
            entry.description.toLowerCase().includes(normalized) ||
            entry.category.toLowerCase().includes(normalized);
    });
}

export function countByCategory(
    entries: CatalogEntry[],
    query: string,
    category: CategoryId,
): number {
    return filterCatalog(entries, query, category).length;
}
