import PluginId from './plugin_id';
import {STATUS_CHANGE, OPEN_ROOT_MODAL, CLOSE_ROOT_MODAL} from './action_types';

export const openRootModal = () => (dispatch) => {
    dispatch({
        type: OPEN_ROOT_MODAL,
    });
};

export const closeRootModal = () => (dispatch) => {
    dispatch({
        type: CLOSE_ROOT_MODAL,
    });
};

export const mainMenuAction = openRootModal;
export const channelHeaderButtonAction = openRootModal;

export const getStatus = () => (dispatch) => {
    fetch('/plugins/' + PluginId + '/').then((r) => r.json()).then((r) => {
        dispatch({
            type: STATUS_CHANGE,
            data: r.enabled,
        });
    });
};

export const websocketStatusChange = (message) => (dispatch) => dispatch({
    type: STATUS_CHANGE,
    data: message.data.enabled,
});
