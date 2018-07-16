import React from 'react';
import PropTypes from 'prop-types';

export default class UserActionsComponent extends React.PureComponent {
    static propTypes = {
        openRootModal: PropTypes.func.isRequired,
        theme: PropTypes.object.isRequired,
    }

    onClick = () => {
        this.props.openRootModal();
    }

    render() {
        const style = getStyle(this.props.theme);

        return (
            <div>
                { 'Sample Plugin: '}
                <button
                    style={style.button}
                    onClick={this.onClick}
                >
                    {'Action'}
                </button>
            </div>
        );
    }
}

const getStyle = (theme) => ({
    button: {
        color: theme.buttonColor,
        backgroundColor: theme.buttonBg,
    },
});
