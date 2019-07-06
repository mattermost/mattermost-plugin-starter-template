import {id as pluginId} from './manifest';

import Icon from './components/icon';

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/

        // registerChannelHeaderButtonAction demonstrates a plugin that add a channel header button
        registry.registerChannelHeaderButtonAction(
            <Icon/>,
            () => {
                alert('Hello World'); // eslint-disable-line no-alert
            },
            'Hello World',
        );
    }
}

window.registerPlugin(pluginId, new Plugin());
