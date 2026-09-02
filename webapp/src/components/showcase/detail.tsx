import React from 'react';

import ChevronLeftIcon from '@mattermost/compass-icons/components/chevron-left';
import {Icon} from '@mattermost/compass-ui/components/icon';
import {IconButton} from '@mattermost/compass-ui/components/icon-button';
import {Scrollbar} from '@mattermost/compass-ui/components/scrollbar';

import type {CatalogEntry} from './types';

type DetailProps = {
    entry: CatalogEntry;
    onBack: () => void;
};

export default function Detail({entry, onBack}: DetailProps) {
    const Preview = entry.Detail;

    return (
        <div className='CompassShowcase CompassShowcase--detail'>
            <div className='CompassShowcase__detail-header'>
                <IconButton
                    icon={(
                        <Icon
                            glyph={<ChevronLeftIcon/>}
                            size='16'
                        />
                    )}
                    aria-label='Back to gallery'
                    onClick={onBack}
                />
                <div className='CompassShowcase__detail-heading'>
                    <h2 className='CompassShowcase__detail-title'>
                        {entry.name}
                    </h2>
                    <p className='CompassShowcase__detail-description'>
                        {entry.description}
                    </p>
                </div>
            </div>
            <Scrollbar className='CompassShowcase__scroll'>
                <div className='CompassShowcase__detail-body'>
                    <Preview/>
                </div>
            </Scrollbar>
        </div>
    );
}
