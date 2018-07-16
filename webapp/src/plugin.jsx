import React from 'react';

import PluginId from './plugin_id';

import Root from './components/root';
import BottomTeamSidebar from './components/bottom_team_sidebar';
import LeftSidebarHeader from './components/left_sidebar_header';
import UserAttributes from './components/user_attributes';
import UserActions from './components/user_actions';
import PostType from './components/post_type';
import {
    MainMenuMobileIcon,
    ChannelHeaderButtonIcon,
} from './components/icons';
import {
    mainMenuAction,
    channelHeaderButtonAction,
    websocketStatusChange,
    getStatus,
} from './actions';
import reducer from './reducer';

export default class SamplePlugin {
    initialize(registry, store) {
        registry.registerRootComponent(Root);
        registry.registerPopoverUserAttributesComponent(UserAttributes);
        registry.registerPopoverUserActionsComponent(UserActions);
        registry.registerLeftSidebarHeaderComponent(LeftSidebarHeader);
        registry.registerBottomTeamSidebarComponent(
            BottomTeamSidebar,
        );

        registry.registerChannelHeaderButtonAction(
            <ChannelHeaderButtonIcon/>,
            () => store.dispatch(channelHeaderButtonAction()),
            'Sample Plugin',
        );

        registry.registerPostTypeComponent('custom_sample_plugin', PostType);

        registry.registerMainMenuAction(
            'Sample Plugin',
            () => store.dispatch(mainMenuAction()),
            <MainMenuMobileIcon/>,
        );

        registry.registerWebSocketEventHandler(
            'custom_' + PluginId + '_status_change',
            (message) => {
                store.dispatch(websocketStatusChange(message));
            },
        );

        registry.registerReducer(reducer);

        // Immediately fetch the current plugin status.
        store.dispatch(getStatus());

        // Fetch the current status whenever we recover an internet connection.
        registry.registerReconnectHandler(() => {
            store.dispatch(getStatus());
        });
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(PluginId + '::uninitialize()');
    }
}
