// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// Compass tokens are declared in cascade layers, so this is safe to load next to the host's
// copy. Never import '@mattermost/compass-ui/styles/standalone' here: it overrides the user's theme.
import '@mattermost/compass-ui/styles';

import {reducer, statusesReset, statusReceived} from 'channel_status/state';
import {setSiteURL} from 'client';
import manifest from 'manifest';
import {FormattedMessage} from 'react-intl';
import type {Store} from 'redux';

import FlagOutlineIcon from '@mattermost/compass-icons/components/flag-outline';
import type {Channel} from '@mattermost/types/channels';
import type {GlobalState} from '@mattermost/types/store';

import {getConfig} from 'mattermost-redux/selectors/entities/general';

import ChannelStatusSidebar from 'components/channel_status_sidebar/channel_status_sidebar';
import ChannelStatusTag from 'components/channel_status_tag/channel_status_tag';

import type {PluginRegistry} from 'types/mattermost-webapp';

import en from '../i18n/en.json';

const translations: Record<string, Record<string, string>> = {en};

export default class Plugin {
    public async initialize(registry: PluginRegistry, store: Store<GlobalState>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        setSiteURL(getConfig(store.getState()).SiteURL);

        // Strings fall back to their defaultMessage for locales without a file.
        registry.registerTranslations((locale) => translations[locale] ?? {});

        registry.registerReducer(reducer);

        // The server sends this to the channel's members whenever its status changes.
        registry.registerWebSocketEventHandler<{channel_id: string; status: string}>(
            `custom_${manifest.id}_status_updated`,
            (msg) => store.dispatch(statusReceived(msg.data.channel_id, msg.data.status)),
        );

        // Events sent while disconnected are lost, so drop what's cached and let components refetch.
        registry.registerReconnectHandler(() => store.dispatch(statusesReset()));

        const {showRHSPlugin, toggleRHSPlugin} = registry.registerRightHandSidebarComponent(
            ChannelStatusSidebar,
            <FormattedMessage
                id='channel_status.title'
                defaultMessage='Channel status'
            />,
        );

        registry.registerChannelHeaderButtonAction(
            <FlagOutlineIcon size={18}/>,
            () => store.dispatch(toggleRHSPlugin),
            'Channel status',
            'Channel status',
        );

        registry.registerChannelHeaderIcon(({channel}: {channel: Channel}) => (
            <ChannelStatusTag
                channel={channel}
                onOpen={() => store.dispatch(showRHSPlugin)}
            />
        ));
    }
}

declare global {
    interface Window {
        registerPlugin(pluginId: string, plugin: Plugin): void;
    }
}

window.registerPlugin(manifest.id, new Plugin());
