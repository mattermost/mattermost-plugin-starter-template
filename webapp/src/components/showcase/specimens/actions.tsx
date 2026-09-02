import React from 'react';

import PaletteOutlineIcon from '@mattermost/compass-icons/components/palette-outline';
import {ActionButton} from '@mattermost/compass-ui/components/action-button';
import {Button} from '@mattermost/compass-ui/components/button';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {IconButton} from '@mattermost/compass-ui/components/icon-button';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function PaletteIcon() {
    return (
        <Icon
            glyph={<PaletteOutlineIcon/>}
            size='16'
        />
    );
}

function ButtonPreview() {
    return <Button size='small'>{'Save'}</Button>;
}

function ButtonDetail() {
    return (
        <>
            <SpecimenRow label='Emphasis'>
                <Button emphasis='primary'>{'Primary'}</Button>
                <Button emphasis='secondary'>{'Secondary'}</Button>
                <Button emphasis='tertiary'>{'Tertiary'}</Button>
                <Button emphasis='quaternary'>{'Quaternary'}</Button>
            </SpecimenRow>
            <SpecimenRow label='Size'>
                <Button size='x-small'>{'XS'}</Button>
                <Button size='small'>{'Small'}</Button>
                <Button size='medium'>{'Medium'}</Button>
                <Button size='large'>{'Large'}</Button>
            </SpecimenRow>
            <SpecimenRow label='Destructive'>
                <Button destructive={true}>{'Delete'}</Button>
                <Button
                    emphasis='secondary'
                    destructive={true}
                >{'Cancel'}</Button>
            </SpecimenRow>
            <SpecimenRow label='Disabled'>
                <Button disabled={true}>{'Disabled'}</Button>
            </SpecimenRow>
        </>
    );
}

function IconButtonPreview() {
    return (
        <IconButton
            icon={<PaletteIcon/>}
            aria-label='Palette'
        />
    );
}

function IconButtonDetail() {
    return (
        <>
            <SpecimenRow label='Default'>
                <IconButton
                    icon={<PaletteIcon/>}
                    aria-label='Palette'
                />
                <IconButton
                    icon={<PaletteIcon/>}
                    toggled={true}
                    aria-label='Toggled'
                />
                <IconButton
                    icon={<PaletteIcon/>}
                    destructive={true}
                    aria-label='Destructive'
                />
            </SpecimenRow>
            <SpecimenRow label='Size'>
                <IconButton
                    icon={<PaletteIcon/>}
                    size='x-small'
                    aria-label='Extra small'
                />
                <IconButton
                    icon={<PaletteIcon/>}
                    size='small'
                    aria-label='Small'
                />
                <IconButton
                    icon={<PaletteIcon/>}
                    size='medium'
                    aria-label='Medium'
                />
                <IconButton
                    icon={<PaletteIcon/>}
                    size='large'
                    aria-label='Large'
                />
            </SpecimenRow>
        </>
    );
}

function ActionButtonPreview() {
    return (
        <ActionButton
            icon={<PaletteIcon/>}
            label='Edit'
        />
    );
}

function ActionButtonDetail() {
    return (
        <SpecimenRow>
            <ActionButton
                icon={<PaletteIcon/>}
                label='Edit'
            />
            <ActionButton
                icon={<PaletteIcon/>}
                label='Selected'
                active={true}
            />
            <ActionButton
                icon={<PaletteIcon/>}
                label='Remove'
                destructive={true}
            />
        </SpecimenRow>
    );
}

export const actionEntries: CatalogEntry[] = [
    {
        id: 'button',
        name: 'Button',
        category: 'actions',
        description: 'Primary actions with emphasis, size, and destructive variants.',
        Preview: ButtonPreview,
        Detail: ButtonDetail,
    },
    {
        id: 'icon-button',
        name: 'Icon Button',
        category: 'actions',
        description: 'Compact, label-free actions for tight surfaces.',
        Preview: IconButtonPreview,
        Detail: IconButtonDetail,
    },
    {
        id: 'action-button',
        name: 'Action Button',
        category: 'actions',
        description: 'Equal-weight actions next to the content they act on.',
        Preview: ActionButtonPreview,
        Detail: ActionButtonDetail,
    },
];
