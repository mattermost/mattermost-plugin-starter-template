// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// Compass tokens are declared in cascade layers, so this is safe to load next to the host's
// copy. Never import '@mattermost/compass-ui/styles/standalone' here: it overrides the user's theme.
import '@mattermost/compass-ui/styles';

import {setSiteURL} from 'client';
import manifest from 'manifest';
import {FormattedMessage} from 'react-intl';
import type {Store} from 'redux';

import HandRightOutlineIcon from '@mattermost/compass-icons/components/hand-right-outline';
import type {GlobalState} from '@mattermost/types/store';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import GreetingsSidebar from 'components/greetings_sidebar/greetings_sidebar';

import type {PluginRegistry} from 'types/mattermost-webapp';

import en from '../i18n/en.json';

const translations: Record<string, Record<string, string>> = {en};

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        setSiteURL(getConfig(store.getState()).SiteURL);

        // Strings fall back to their defaultMessage for locales without a file.
        registry.registerTranslations((locale) => translations[locale] ?? {});

        const {toggleRHSPlugin} = registry.registerRightHandSidebarComponent(
            GreetingsSidebar,
            <FormattedMessage
                id='greetings.title'
                defaultMessage='Greetings'
            />,
        );

        registry.registerChannelHeaderButtonAction(
            <HandRightOutlineIcon size={18}/>,
            () => store.dispatch(toggleRHSPlugin),
            'Greetings',
            'Open greetings',
        );
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
