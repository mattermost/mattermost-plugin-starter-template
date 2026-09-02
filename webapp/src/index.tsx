// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import manifest from 'manifest';
import React from 'react';
import type {Store} from 'redux';

import PaletteOutlineIcon from '@mattermost/compass-icons/components/palette-outline';
import {Icon} from '@mattermost/compass-ui/components/icon';
import type {GlobalState} from '@mattermost/types/store';

import ShowcaseRHS from 'components/showcase/rhs';

import type {PluginRegistry} from 'types/mattermost-webapp';

import '@mattermost/compass-ui/component-styles';
import '@mattermost/compass-ui/styles';

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState>) {
        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(
            ShowcaseRHS,
            'Compass UI',
        );

        registry.registerChannelHeaderButtonAction(
            <Icon
                glyph={<PaletteOutlineIcon/>}
                size='24'
            />,
            () => store.dispatch(toggleRHSPlugin),
            'Compass UI',
            'Browse Compass UI components',
        );

        registry.registerAppBarComponent(
            `/plugins/${manifest.id}/public/app-bar-icon.svg`,
            () => store.dispatch(toggleRHSPlugin),
            'Browse Compass UI components',
            null,
        );
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
