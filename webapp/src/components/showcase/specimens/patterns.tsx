import React from 'react';

import {Button} from '@mattermost/compass-ui/components/button';
import {MenuItem} from '@mattermost/compass-ui/components/menu-item';
import {Modal} from '@mattermost/compass-ui/components/modal';
import {PopoverMenu, PopoverMenuTitle} from '@mattermost/compass-ui/components/popover-menu';
import {ProfilePopover} from '@mattermost/compass-ui/components/profile-popover';
import {TourPoint} from '@mattermost/compass-ui/components/tour-point';

import {PLACEHOLDER_AVATAR, SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

function ModalPreview() {
    return (
        <div className='CompassShowcase__modal-frame'>
            <Modal
                size='small'
                title='Rename channel'
                onClose={() => undefined}
                footer={<Button size='small'>{'Save'}</Button>}
            >
                {'Chrome only — host owns open, portal, and focus trap.'}
            </Modal>
        </div>
    );
}

function ModalDetail() {
    return (
        <SpecimenRow>
            <div className='CompassShowcase__modal-frame'>
                <Modal
                    size='small'
                    title='Invite people'
                    subtitle='Add members to this channel'
                    onClose={() => undefined}
                    footer={(
                        <>
                            <Button
                                emphasis='tertiary'
                                size='small'
                            >{'Cancel'}</Button>
                            <Button size='small'>{'Invite'}</Button>
                        </>
                    )}
                >
                    {'This is the modal body. Positioning is up to the product host.'}
                </Modal>
            </div>
        </SpecimenRow>
    );
}

function PopoverMenuPreview() {
    return (
        <PopoverMenu>
            <MenuItem label='Copy text'/>
            <MenuItem label='Pin to channel'/>
        </PopoverMenu>
    );
}

function PopoverMenuDetail() {
    return (
        <SpecimenRow>
            <PopoverMenu>
                <PopoverMenuTitle>{'Actions'}</PopoverMenuTitle>
                <MenuItem label='Reply'/>
                <MenuItem label='Follow'/>
                <MenuItem
                    label='Delete'
                    destructive={true}
                />
            </PopoverMenu>
        </SpecimenRow>
    );
}

function ProfilePopoverPreview() {
    return (
        <div className='CompassShowcase__popover-frame'>
            <ProfilePopover
                avatarSrc={PLACEHOLDER_AVATAR}
                avatarAlt='Ada Lovelace'
                name='Ada Lovelace'
                username='ada.lovelace'
                title='Mathematician'
            />
        </div>
    );
}

function ProfilePopoverDetail() {
    return (
        <SpecimenRow>
            <div className='CompassShowcase__popover-frame'>
                <ProfilePopover
                    avatarSrc={PLACEHOLDER_AVATAR}
                    avatarAlt='Ada Lovelace'
                    name='Ada Lovelace'
                    username='ada.lovelace'
                    title='Mathematician'
                    email='ada@mattermost.com'
                    jobRole='System Admin'
                    onClose={() => undefined}
                />
            </div>
        </SpecimenRow>
    );
}

function TourPointPreview() {
    return (
        <TourPoint
            title='Browse components'
            pointerPosition='none'
            showPulsingDot={false}
        >
            {'Search or pick a category, then open a component to see variants.'}
        </TourPoint>
    );
}

function TourPointDetail() {
    return (
        <SpecimenRow>
            <TourPoint
                title='Channel header'
                pointerPosition='top-left'
                showPulsingDot={false}
                progress={{pages: 3, activePage: 1}}
                primaryAction={{label: 'Next'}}
                onClose={() => undefined}
            >
                {'Open Compass UI from the channel header to browse the catalog.'}
            </TourPoint>
        </SpecimenRow>
    );
}

export const patternEntries: CatalogEntry[] = [
    {
        id: 'modal',
        name: 'Modal',
        category: 'patterns',
        description: 'Dialog chrome only. Host owns portal, focus, and stacking.',
        Preview: ModalPreview,
        Detail: ModalDetail,
    },
    {
        id: 'popover-menu',
        name: 'Popover Menu',
        category: 'patterns',
        description: 'Menu surface composed with MenuItem rows.',
        Preview: PopoverMenuPreview,
        Detail: PopoverMenuDetail,
    },
    {
        id: 'profile-popover',
        name: 'Profile Popover',
        category: 'patterns',
        description: 'User card chrome. Host owns trigger and positioning.',
        Preview: ProfilePopoverPreview,
        Detail: ProfilePopoverDetail,
    },
    {
        id: 'tour-point',
        name: 'Tour Point',
        category: 'patterns',
        description: 'Onboarding callout with optional progress dots.',
        Preview: TourPointPreview,
        Detail: TourPointDetail,
    },
];
