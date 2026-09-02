import React from 'react';

import {PaginationDots} from '@mattermost/compass-ui/components/pagination-dots';
import {ProgressBar} from '@mattermost/compass-ui/components/progress-bar';
import {Spinner} from '@mattermost/compass-ui/components/spinner';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function SpinnerPreview() {
    return (
        <Spinner
            size={20}
            aria-label='Loading'
        />
    );
}

function SpinnerDetail() {
    return (
        <SpecimenRow>
            <Spinner
                size={12}
                aria-label='12'
            />
            <Spinner
                size={16}
                aria-label='16'
            />
            <Spinner
                size={24}
                aria-label='24'
            />
            <Spinner
                size={32}
                aria-label='32'
            />
        </SpecimenRow>
    );
}

function ProgressBarPreview() {
    return (
        <ProgressBar
            value={64}
            aria-label='Upload 64 percent'
        />
    );
}

function ProgressBarDetail() {
    return (
        <>
            <SpecimenRow>
                <ProgressBar
                    value={35}
                    size='small'
                    aria-label='Small'
                />
            </SpecimenRow>
            <SpecimenRow>
                <ProgressBar
                    value={80}
                    semanticColors={true}
                    aria-label='Semantic'
                />
            </SpecimenRow>
        </>
    );
}

function PaginationDotsPreview() {
    return (
        <PaginationDots
            pages={4}
            activePage={2}
        />
    );
}

function PaginationDotsDetail() {
    return (
        <SpecimenRow>
            <PaginationDots
                pages={5}
                activePage={3}
            />
        </SpecimenRow>
    );
}

export const progressEntries: CatalogEntry[] = [
    {
        id: 'spinner',
        name: 'Spinner',
        category: 'progress',
        description: 'Indeterminate wait indicator.',
        Preview: SpinnerPreview,
        Detail: SpinnerDetail,
    },
    {
        id: 'progress-bar',
        name: 'Progress Bar',
        category: 'progress',
        description: 'Determinate progress from 0 to 100.',
        Preview: ProgressBarPreview,
        Detail: ProgressBarDetail,
    },
    {
        id: 'pagination-dots',
        name: 'Pagination Dots',
        category: 'progress',
        description: 'Step indicator for a fixed sequence.',
        Preview: PaginationDotsPreview,
        Detail: PaginationDotsDetail,
    },
];
