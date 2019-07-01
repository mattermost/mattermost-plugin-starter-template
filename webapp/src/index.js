import {id as pluginId} from './manifest';

import Settings from './components/settings.jsx';

// eslint-disable-next-line no-unused-vars
function addToAdminConsole(registry) {
    registry.registerAdminConsolePlugin((adminDefinition) => {
        // the settings will appear on the authentication section
        adminDefinition.authentication.starter = {
            url: 'plugins/starter',
            icon: 'fa-rocket',
            title: 'com.mattermost.starter-plugin.title',
            title_default: 'Starter',
            isHidden: () => false,
            schema: {
                id: 'StarterPlugin',
                component: Settings,
            },
        };
        return adminDefinition;
    });
}

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/

        // addToAdminConsole(registry);
    }
}

window.registerPlugin(pluginId, new Plugin());
