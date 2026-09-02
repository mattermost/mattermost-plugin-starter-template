import React from 'react';

import {ErrorMessage} from '@mattermost/compass-ui/components/error-message';
import {PopoverNotice} from '@mattermost/compass-ui/components/popover-notice';
import {SectionNotice} from '@mattermost/compass-ui/components/section-notice';
import {Toast} from '@mattermost/compass-ui/components/toast';
import {Tooltip} from '@mattermost/compass-ui/components/tooltip';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function ErrorMessagePreview() {
    return <ErrorMessage message='Name is required'/>;
}

function ErrorMessageDetail() {
    return (
        <SpecimenRow>
            <ErrorMessage message='This field is required'/>
        </SpecimenRow>
    );
}

function SectionNoticePreview() {
    return (
        <SectionNotice
            type='info'
            title='Hint'
            description='Try another filter.'
        />
    );
}

function SectionNoticeDetail() {
    return (
        <>
            <SpecimenRow>
                <SectionNotice
                    type='info'
                    title='Info'
                    description='Something to know.'
                />
            </SpecimenRow>
            <SpecimenRow>
                <SectionNotice
                    type='success'
                    title='Success'
                    description='Saved.'
                />
            </SpecimenRow>
            <SpecimenRow>
                <SectionNotice
                    type='warning'
                    title='Warning'
                    description='Check this setting.'
                />
            </SpecimenRow>
            <SpecimenRow>
                <SectionNotice
                    type='danger'
                    title='Danger'
                    description='This cannot be undone.'
                />
            </SpecimenRow>
        </>
    );
}

function ToastPreview() {
    return (
        <Toast
            type='success'
            message='Copied to clipboard'
        />
    );
}

function ToastDetail() {
    return (
        <>
            <SpecimenRow>
                <Toast
                    type='general'
                    message='Working…'
                />
            </SpecimenRow>
            <SpecimenRow>
                <Toast
                    type='info'
                    message='New version available'
                />
            </SpecimenRow>
            <SpecimenRow>
                <Toast
                    type='success'
                    message='Invitation sent'
                />
            </SpecimenRow>
            <SpecimenRow>
                <Toast
                    type='danger'
                    message='Could not save'
                />
            </SpecimenRow>
        </>
    );
}

function TooltipPreview() {
    return <Tooltip label='Edit channel'/>;
}

function TooltipDetail() {
    return (
        <SpecimenRow>
            <Tooltip
                label='Search'
                hint='Find messages in this channel'
                shortcutKeys={[{label: 'Ctrl'}, {label: 'F'}]}
                arrow='bottom'
            />
        </SpecimenRow>
    );
}

function PopoverNoticePreview() {
    return (
        <PopoverNotice
            title='Tip'
            variant='info'
        >
            {'Pin important messages for the team.'}
        </PopoverNotice>
    );
}

function PopoverNoticeDetail() {
    return (
        <SpecimenRow>
            <PopoverNotice
                title='Saved'
                variant='success'
                actions={[{label: 'Undo', emphasis: 'tertiary'}]}
                onClose={() => undefined}
            >
                {'This message was added to Saved.'}
            </PopoverNotice>
        </SpecimenRow>
    );
}

export const feedbackEntries: CatalogEntry[] = [
    {
        id: 'error-message',
        name: 'Error Message',
        category: 'feedback',
        description: 'Inline validation text under a form field.',
        Preview: ErrorMessagePreview,
        Detail: ErrorMessageDetail,
    },
    {
        id: 'section-notice',
        name: 'Section Notice',
        category: 'feedback',
        description: 'In-context alert inside a settings or admin panel.',
        Preview: SectionNoticePreview,
        Detail: SectionNoticeDetail,
    },
    {
        id: 'toast',
        name: 'Toast',
        category: 'feedback',
        description: 'Brief confirmation that sits on top of the workflow.',
        Preview: ToastPreview,
        Detail: ToastDetail,
    },
    {
        id: 'tooltip',
        name: 'Tooltip',
        category: 'feedback',
        description: 'Chrome-only hover label. Host owns trigger and positioning.',
        Preview: TooltipPreview,
        Detail: TooltipDetail,
    },
    {
        id: 'popover-notice',
        name: 'Popover Notice',
        category: 'feedback',
        description: 'Floating tip or confirmation without taking over the screen.',
        Preview: PopoverNoticePreview,
        Detail: PopoverNoticeDetail,
    },
];
