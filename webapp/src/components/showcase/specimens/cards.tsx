import React from 'react';

import {AttachmentCard} from '@mattermost/compass-ui/components/attachment-card';
import {LinkPreview} from '@mattermost/compass-ui/components/link-preview';
import {PermalinkPreview} from '@mattermost/compass-ui/components/permalink-preview';

import {PLACEHOLDER_AVATAR, SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function AttachmentCardPreview() {
    return (
        <AttachmentCard
            fileName='spec.pdf'
            fileMeta='PDF 128KB'
            fileType='pdf'
        />
    );
}

function AttachmentCardDetail() {
    return (
        <>
            <SpecimenRow>
                <AttachmentCard
                    fileName='notes.txt'
                    fileMeta='TXT 15KB'
                    fileType='text'
                />
            </SpecimenRow>
            <SpecimenRow>
                <AttachmentCard
                    fileName='upload.zip'
                    fileMeta='ZIP 2MB'
                    fileType='zip'
                    state='uploading'
                    progress={42}
                />
            </SpecimenRow>
        </>
    );
}

function LinkPreviewPreview() {
    return (
        <LinkPreview
            siteName='mattermost.com'
            title='Compass Design System'
            description='Components for Mattermost products.'
            imageSize='none'
        />
    );
}

function LinkPreviewDetail() {
    return (
        <SpecimenRow>
            <LinkPreview
                siteName='developers.mattermost.com'
                title='Plugin documentation'
                description='Extend Mattermost with plugins.'
                imageSize='none'
            />
        </SpecimenRow>
    );
}

function PermalinkPreviewPreview() {
    return (
        <PermalinkPreview
            authorName='Ada Lovelace'
            avatarSrc={PLACEHOLDER_AVATAR}
            timestamp='2:14 PM'
            messageText='Ship the RHS showcase today.'
            originalChannel='~design'
        />
    );
}

function PermalinkPreviewDetail() {
    return (
        <SpecimenRow>
            <PermalinkPreview
                authorName='Grace Hopper'
                avatarSrc={PLACEHOLDER_AVATAR}
                timestamp='Yesterday'
                messageText='Let’s review the gallery layout.'
                originalChannel='~ux'
            />
        </SpecimenRow>
    );
}

export const cardEntries: CatalogEntry[] = [
    {
        id: 'attachment-card',
        name: 'Attachment Card',
        category: 'cards',
        description: 'File identity and actions in a compact card.',
        Preview: AttachmentCardPreview,
        Detail: AttachmentCardDetail,
    },
    {
        id: 'link-preview',
        name: 'Link Preview',
        category: 'cards',
        description: 'Rich URL metadata under a message.',
        Preview: LinkPreviewPreview,
        Detail: LinkPreviewDetail,
    },
    {
        id: 'permalink-preview',
        name: 'Permalink Preview',
        category: 'cards',
        description: 'Inline card for a linked message.',
        Preview: PermalinkPreviewPreview,
        Detail: PermalinkPreviewDetail,
    },
];
