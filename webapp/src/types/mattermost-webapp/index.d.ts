// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

/* eslint-disable max-lines */

import type {Reducer} from 'redux';

import type {WebSocketMessage} from '@mattermost/client';
import type {Channel} from '@mattermost/types/channels';
import type {FileInfo} from '@mattermost/types/files';
import type {Post, PostEmbed} from '@mattermost/types/posts';
import type {ProductScope} from '@mattermost/types/products';

export type UniqueIdentifier = string;
export type ReactResolvable = React.ReactNode | React.ElementType;

export type PluginComponent = {
    id: string;
    pluginId: string;
    title?: string;

    /** @default null - which means 'channels'*/
    supportedProductIds?: ProductScope;
    component?: React.ComponentType;
    subMenu?: Menu[];
    text?: string;
    dropdownText?: string;
    tooltipText?: string;
    button?: React.ReactElement;
    dropdownButton?: React.ReactElement;
    icon?: React.ReactElement;
    iconUrl?: string;
    mobileIcon?: React.ReactElement;
    filter?: (id: string) => boolean;
    action?: (...args: unknown) => void;
    shouldRender?: (state: GlobalState) => boolean;
    hook?: (post: Post, message?: string) => string;
};

export type PluginSiteStatsHandler = () => Promise<Record<string, PluginAnalyticsRow>>;

export type DesktopNotificationArgs = {
    title: string;
    body: string;
    silent: boolean;
    soundName: string;
    url: string;
    notify: boolean;
};

export type NewPostMessageProps = {
    channel_type: ChannelType;
    channel_display_name: string;
    channel_name: string;
    sender_name: string;
    set_online: boolean;
    mentions?: string;
    followers?: string;
    team_id: string;
    should_ack: boolean;
    otherFile?: 'true';
    image?: 'true';
    post: string;
}

export type PluginConfiguration = {

    /** Plugin ID  */
    id: string;

    /** Name of the plugin to show in the UI. We recommend to use manifest.name */
    uiName: string;

    /** URL to the icon to show in the UI. No icon will show the plug outline icon. */
    icon?: string;

    /** Action that will appear at the beginning of the plugin settings tab */
    action?: PluginConfigurationAction;
    sections: Array<PluginConfigurationSection | PluginConfigurationCustomSection>;
}

export type PluginConfigurationAction = {

    /** Text shown as the title of the action */
    title: string;

    /** Text shown as the body of the action */
    text: string;

    /** Text shown at the button */
    buttonText: string;

    /** This function is called when the button on the action is clicked */
    onClick: () => void;
}

export type PluginConfigurationSection = {
    settings: PluginConfigurationSetting[];

    /** The title of the section. All titles must be different. */
    title: string;

    /** Whether the section is disabled. */
    disabled?: boolean;

    /**
     * This function will be called whenever a section is saved.
     *
     * The configuration will be automatically saved in the user preferences,
     * so use this function only in case you want to add some side effect
     * to the change.
    */
    onSubmit?: (changes: {[name: string]: string}) => void;
}

export type PluginConfigurationCustomSection = {

    /** The title of the section. All titles must be different. */
    title: string;

    /** A React component used to render the custom section. */
    component: React.ComponentType;
}

export interface PluginRegistry {

    /**
        * Register a component at the root of the channel view of the app.
        * Accepts a React component. Returns a unique identifier.
    */
    registerRootComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component in the user attributes section of the profile popover (hovercard), below the default user attributes.
        * Accepts a React component. Returns a unique identifier.
    */
    registerPopoverUserAttributesComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component in the user actions of the profile popover (hovercard), below the default actions.
        * Accepts a React component. Returns a unique identifier.
    */
    registerPopoverUserActionsComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component fixed to the top of the left-hand channel sidebar.
        * Accepts a React component. Returns a unique identifier.
    */
    registerLeftSidebarHeaderComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component fixed to the bottom of the team sidebar. Does not render if
        * user is only on one team and the team sidebar is not shown.
        * Accepts a React component. Returns a unique identifier.
    */
    registerBottomTeamSidebarComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component fixed to the bottom of the post message.
        * Accepts a React component. Returns a unique identifier.
    */
    registerPostMessageAttachmentComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to show as a tooltip when a user hovers on a link in a post.
        * Accepts a React component. Returns a unique identifier.
        * The component will be passed the following props:
        * - href - The URL for this link
        * - show - A boolean used to signal that the user is currently hovering over this link. Use this value to initialize your component when this boolean is true for the first time, using `componentDidUpdate` or `useEffect`.
    */
    registerLinkTooltipComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component fixed to the bottom of the create new channel modal and also registers a callback function to be called after
        * the channel has been succesfully created
        * Accepts a React component. Returns a unique identifier.
    */
    registerActionAfterChannelCreation(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Add a button to the channel header. If there are more than one buttons registered by any
        * plugin, a dropdown menu is created to contain all the plugin buttons.
        * Accepts the following:
        * - icon - React element to use as the button's icon
        * - action - a function called when the button is clicked, passed the channel and channel member as arguments
        * - dropdownText - string or React element shown for the dropdown button description
        * - tooltipText - string or React element shown for tooltip appear on hover
    */
    registerChannelHeaderButtonAction(
        ...args: [
            icon: ReactResolvable,
            action: () => void,
            dropdownText: string,
            tooltipText: string
        ] | [{
            icon: ReactResolvable;
            action: () => void;
            dropdownText: string;
            tooltipText: string;
        }]
    ): UniqueIdentifier;

    /**
        * Add a button to the channel intro message.
        * Accepts the following:
        * - icon - React element to use as the button's icon
        * - action - a function called when the button is clicked, passed the channel and channel member as arguments
        * - text - a localized string or React element  to use as the button's text
    */
    registerChannelIntroButtonAction(
        ...args: [
            icon: ReactResolvable,
            action: () => void,
            tooltipText: ReactResolvable
        ] | [{
            icon: ReactResolvable;
            action: () => void;
            tooltipText: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Add a "call button" to the channel header. If there is more than one button registered by any
        * plugin, a dropdown menu is created to contain all the call plugin buttons.
        * Accepts the following:
        * - button - A React element to use as the main button to be displayed in case of a single registration.
        * - dropdownButton -A React element to use as the dropdown button to be displayed in case of multiple registrations.
        * - action - A function called when the button is clicked, passed the channel and channel member as arguments.
        * Returns an unique identifier
        * Minimum required version: 6.5
    */
    registerCallButtonAction(
        ...args: [
            button: ReactResolvable,
            dropdownButton: ReactResolvable,
            fn: (channel: Channel) => void
        ] | [{
            button: ReactResolvable;
            dropdownButton: ReactResolvable;
            fn: (channel: Channel) => void;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to render a custom body for posts with a specific type.
        * Custom post types must be prefixed with 'custom_'.
        * Custom post types can also apply for ephemeral posts.
        * Accepts a string type and a component.
        * Returns a unique identifier.
    */
    registerPostTypeComponent(
        ...args: [
            typeName: string,
            component: ReactResolvable
        ] | [{
            typeName: string;
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to render a custom body for post cards with a specific type.
        * Custom post types must be prefixed with 'custom_'.
        * Accepts a string type and a component.
        * Returns a unique identifier.
    */
    registerPostCardTypeComponent(
        ...args: [
            type: string,
            component: ReactResolvable
        ] | [{
            type: string;
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to render a custom embed preview for post links.
        * Accepts the following:
        * - match - A function that receives the embed object and returns a
        *   boolean indicating if the plugin is able to process it.
        *   The embed object contains the embed `type`, the `url` of the post link
        *   and in some cases, a `data` object with information related to the
        *   link (the opengraph or the image details, for example).
        * - component - The component that renders the embed view for the link
        * - toggleable - A boolean indicating if the embed view should be collapsable
        * Returns a unique identifier.
    */
    registerPostWillRenderEmbedComponent(
        ...args: [
            match: (embed: PostEmbed) => void,
            component: ReactResolvable,
            toggleable: boolean
        ] | [{
            match: (embed: PostEmbed) => void;
            component: ReactResolvable;
            toggleable: boolean;
        }]
    ): UniqueIdentifier;

    /**
        * Register a main menu list item by providing some text and an action function.
        * Accepts the following:
        * - text - A string or React element to display in the menu
        * - action - A function to trigger when component is clicked on
        * - mobileIcon - A React element to display as the icon in the menu in mobile view
        * Returns a unique identifier.
    */
    registerMainMenuAction(
        ...args: [
            text: ReactResolvable,
            action: () => void,
            mobileIcon: ReactResolvable
        ] | [{
            text: ReactResolvable;
            action: () => void;
            mobileIcon: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a channel menu list item by providing some text and an action function.
        * Accepts the following:
        * - text - A string or React element to display in the menu
        * - action - A function that receives the channelId and is called when the menu items is clicked.
        * - shouldRender - A function that receives the state before the
        * component is about to render, allowing for conditional rendering.
        * Returns a unique identifier.
    */
    registerChannelHeaderMenuAction(
        ...args: [
            component: ReactResolvable,
            fn: (channelID: string) => void
        ] | [{
            component: ReactResolvable;
            fn: (channelID: string) => void;
        }]
    ): UniqueIdentifier;

    /**
        * Register a files dropdown list item by providing some text and an action function.
        * Accepts the following:
        * - match - A function  that receives the fileInfo and returns a boolean indicating if the plugin is able to process it.
        * - text - A string or React element to display in the menu
        * - action - A function that receives the fileInfo and is called when the menu items is clicked.
        * Returns a unique identifier.
    */
    registerFileDropdownMenuAction(
        ...args: [
            match: (fileInfo: FileInfo) => boolean,
            text: ReactResolvable,
            action: (fileInfo: FileInfo) => void
        ] | [{
            match: (fileInfo: FileInfo) => boolean;
            text: ReactResolvable;
            action: (fileInfo: FileInfo) => void;
        }]
    ): UniqueIdentifier;

    /**
        * Register a user guide dropdown list item by providing some text and an action function.
        * Accepts the following:
        * - text - A string or React element to display in the menu
        * - action - A function that receives the fileInfo and is called when the menu items is clicked.
        * Returns a unique identifier.
    */
    registerUserGuideDropdownMenuAction(
        ...args: [
            text: ReactResolvable,
            action: (fileInfo: FileInfo) => void
        ] | [{
            text: ReactResolvable;
            action: (fileInfo: FileInfo) => void;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to the add to the post message menu shown on hover.
        * Accepts a React component. Returns a unique identifier.
    */
    registerPostActionComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to the add to the post text editor menu.
        * Accepts a React component. Returns a unique identifier.
    */
    registerPostEditorActionComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to the add to the code block header.
        * Accepts a React component. Returns a unique identifier.
    */
    registerCodeBlockActionComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to the add to the new messages separator.
        * Accepts a React component. Returns a unique identifier.
    */
    registerNewMessagesSeparatorActionComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a post menu list item by providing some text and an action function.
        * Accepts the following:
        * - text - A string or React element to display in the menu
        * - action - A function to trigger when component is clicked on
        * - filter - A function whether to apply the plugin into the post' dropdown menu
        * Returns a unique identifier.
    */
    registerPostDropdownMenuAction(
        ...args: [
            text: React.ReactNode,
            action: () => void,
            filter: (post: Post) => boolean
        ] | [{
            text: React.ReactNode;
            action: () => void;
            filter: (post: Post) => boolean;
        }]
    ): UniqueIdentifier;

    /**
        * Register a post sub menu list item by providing some text and an action function.
        * Accepts the following:
        * - text - A string or React element to display in the menu
        * - action - A function to trigger when component is clicked on
        * - filter - A function whether to apply the plugin into the post' dropdown menu
        *
        * Returns a unique identifier for the root submenu, and a function to register submenu items.
        * At this time, only one level of nesting is allowed to avoid rendering issue in the RHS.
    */
    registerPostDropdownSubMenuAction(
        ...args: [
            text: ReactResolvable,
            action: (postId: string) => void,
            filter: (postId: string) => boolean
        ] | [{
            text: ReactResolvable;
            action: (postId: string) => void;
            filter: (postId: string) => boolean;
        }]
    ): {
        id: UniqueIdentifier;
        rootRegisterMenuItem: (
            text: React.ReactNode,
            action: () => void,
            filter: (post: Post) => boolean
        ) => void;
    };

    /**
        * Register a component at the bottom of the post dropdown menu.
        * Accepts a React component. Returns a unique identifier.
    */
    registerPostDropdownMenuComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * Register a file upload method by providing some text, an icon, and an action function.
        * Accepts the following:
        * - icon - JSX element to use as the button's icon
        * - text - A string or JSX element to display in the file upload menu
        * - action - A function to trigger when the menu item is selected.
        * Returns a unique identifier.
    */
    registerFileUploadMethod(
        ...args: [
            icon: ReactResolvable,
            action: (checkPluginHooksAndUploadFiles: ((files: FileList | File[]) => void)) => void,
            text: ReactResolvable
        ] | [{
            icon: ReactResolvable;
            action: (checkPluginHooksAndUploadFiles: ((files: FileList | File[]) => void)) => void;
            text: ReactResolvable;
        }]
    ): string;

    /**
        * Register a hook to intercept file uploads before they take place.
        * Accepts a function to run before files get uploaded. Receives an array of
        * files and a function to upload files at a later time as arguments. Must
        * return an object that can contain two properties:
        * - message - An error message to display, leave blank or null to display no message
        * - files - Modified array of files to upload, set to null to reject all files
        * Returns a unique identifier.
    */
    registerFilesWillUploadHook(
        ...args: [
            hook: (files: File[], uploadFiles: (files: File[]) => void) => {
                message?: string;
                files?: File[];
            }
        ] | [{
            hook: (files: File[], uploadFiles: (files: File[]) => void) => {
                message?: string;
                files?: File[];
            };
        }]
    ): UniqueIdentifier;

    /**
        * Unregister a component, action or hook using the unique identifier returned after registration.
        * Accepts a string id.
        * Returns undefined in all cases.
    */
    unregisterComponent(
        ...args: [
            componentId: UniqueIdentifier
        ] | [{
            componentId: UniqueIdentifier;
        }]
    ): void;

    /**
        * Unregister a component that provided a custom body for posts with a specific type.
        * Accepts a string id.
        * Returns undefined in all cases.
    */
    unregisterPostTypeComponent(
        ...args: [
            componentId: UniqueIdentifier
        ] | [{
            componentId: UniqueIdentifier;
        }]
    ): void;

    /**
        * Register a reducer against the Redux store. It will be accessible in redux state
        * under "state['plugins-<yourpluginid>']"
        * Accepts a reducer. Returns undefined.
    */
    registerReducer(
        ...args: [
            reducer: Reducer
        ] | [{
            reducer: Reducer;
        }]
    ): string;

    /**
        * Register a handler for WebSocket events.
        * Accepts the following:
        * - event - the event type, can be a regular server event or an event from plugins.
        * Plugin events will have "custom_<pluginid>_" prepended
        * - handler - a function to handle the event, receives the event message as an argument
        * Returns undefined.
    */
    registerWebSocketEventHandler<T = Record<string, string>>(
        ...args: [
            event: string,
            handler: (msg: WebSocketMessage<T>) => void
        ] | [{
            event: string;
            handler: (msg: WebSocketMessage<T>) => void;
        }]
    ): void;

    /**
        * Unregister a handler for a custom WebSocket event.
        * Accepts a string event type.
        * Returns undefined.
    */
    unregisterWebSocketEventHandler(
        ...args: [
            event: string
        ] | [{
            event: string;
        }]
    ): void;

    /**
        * Register a handler that will be called when the app reconnects to the
        * internet after previously disconnecting.
        * Accepts a function to handle the event. Returns undefined.
    */
    registerReconnectHandler(
        ...args: [
            handler: () => void
        ] | [{
            handler: () => void;
        }]
    ): void;

    /**
        * Unregister a previously registered reconnect handler.
        * Returns undefined.
    */
    unregisterReconnectHandler(): void;

    /**
        * Register a hook that will be called when a message is posted by the user before it
        * is sent to the server. Accepts a function that receives the post as an argument.
        *
        * To reject a post, return an object containing an error such as
        *     {error: {message: 'Rejected'}}
        * To modify or allow the post without modification, return an object containing the post
        * such as
        *     {post: {...}}
        *
        * If the hook function is asynchronous, the message will not be sent to the server
        * until the hook returns.
    */
    registerMessageWillBePostedHook(
        ...args: [
            hook: (post: Post) => ({ post: Post } | {error: { message: string }} | Promise<{ post: Post } | { error: { message: string } }>)
        ] | [{
            hook: (post: Post) => ({ post: Post } | {error: { message: string }} | Promise<{ post: Post } | { error: { message: string } }>);
        }]
    ): UniqueIdentifier;

    /**
        * Register a hook that will be called when a slash command is posted by the user before it
        * is sent to the server. Accepts a function that receives the message (string) and the args
        * (object) as arguments.
        * The args object is:
        *        {
        *            channel_id: channelId,
        *            team_id: teamId,
        *            root_id: rootId,
        *        }
        *
        * To reject a command, return an object containing an error:
        *     {error: {message: 'Rejected'}}
        * To ignore a command, return an empty object (to prevent an error from being displayed):
        *     {}
        * To modify or allow the command without modification, return an object containing the new message
        * and args. It is not likely that you will need to change the args, so return the object that was provided:
        *     {message: {...}, args}
        *
        * If the hook function is asynchronous, the command will not be sent to the server
        * until the hook returns.
    */
    registerSlashCommandWillBePostedHook(
        ...args: [
            hook: (message: string, args: ContextArgs) => ({ message: string; args: ContextArgs } | object | Promise<{ message: string; args: ContextArgs } | object>)
        ] | [{
            hook: (message: string, args: ContextArgs) => ({ message: string; args: ContextArgs } | object | Promise<{ message: string; args: ContextArgs } | object>);
        }]
    ): UniqueIdentifier;

    /**
        * Register a hook that will be called before a message is formatted into Markdown.
        * Accepts a function that receives the unmodified post and the message (potentially
        * already modified by other hooks) as arguments. This function must return a string
        * message that will be formatted.
        * Returns a unique identifier.
    */
    registerMessageWillFormatHook(
        ...args: [
            hook: (post: Post, message: string) => string
        ] | [{
            hook: (post: Post, message: string) => string;
        }]
    ): UniqueIdentifier;

    /**
        * Register a component to override file previews. Accepts a function to run before file is
        * previewed and a react component to be rendered as the file preview.
        * - override - A function to check whether preview needs to be overridden. Receives fileInfo and post as arguments.
        * Returns true is preview should be overridden and false otherwise.
        * - component - A react component to display instead of original preview. Receives fileInfo and post as props.
        * Returns a unique identifier.
        * Only one plugin can override a file preview at a time. If two plugins try to override the same file preview, the first plugin will perform the override and the second will not. Plugin precedence is ordered alphabetically by plugin ID.
    */
    registerFilePreviewComponent(
        ...args: [
            override: (fileInfos: FileInfo[], post: Post) => boolean,
                    component: ReactResolvable
        ] | [{
            override: (fileInfos: FileInfo[], post: Post) => boolean;
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    registerTranslations(
        ...args: [
            getTranslationsForLocale: (locale: string) => { [translationId: string]: string }
        ] | [{
            getTranslationsForLocale: (locale: string) => { [translationId: string]: string };
        }]
    ): void;

    /**
        * Register a admin console definitions override function
        * Note that this is a low-level interface primarily meant for internal use, and is not subject
        * to semver guarantees. It may change in the future.
        * Accepts the following:
        * - func - A function that recieve the admin console config definitions and return a new
        *          version of it, which is used for build the admin console.
        * Each plugin can register at most one admin console plugin function, with newer registrations
        * replacing older ones.
    */
    registerAdminConsolePlugin(
        ...args: [
            func: (config: object) => void
        ] | [{
            func: (config: object) => void;
        }]
    ): void;

    /**
        * Unregister a previously registered admin console definition override function.
        * Returns undefined.
    */
    unregisterAdminConsolePlugin(): void;

    /**
        * Register a custom React component to manage the plugin configuration for the given setting key.
        * Accepts the following:
        * - key - A key specified in the settings_schema.settings block of the plugin's manifest.
        * - component - A react component to render in place of the default handling.
        * - options - Object for the following available options to display the setting:
        *     showTitle - Optional boolean that if true the display_name of the setting will be rendered
        * on the left column of the settings page and the registered component will be displayed on the
        * available space in the right column.
    */
    registerAdminConsoleCustomSetting(
        ...args: [
            key: string,
            component: ReactResolvable,
            options?: { showTitle?: boolean }
        ] | [{
            key: string;
            component: ReactResolvable;
            options?: { showTitle?: boolean };
        }]
    ): void;

    /**
        * Register a custom React component to render as a section in the plugin configuration page.
        * Accepts the following:
        * - key - A key specified in the settings_schema.sections block of the plugin's manifest.
        * - component - A react component to render in place of the default handling.
    */
    registerAdminConsoleCustomSection(
        ...args: [
            key: string,
            component: ReactResolvable
        ] | [{
            key: string;
            component: ReactResolvable;
        }]
    ): void;

    /**
        * Register a Right-Hand Sidebar component by providing a title for the right hand component.
        * Accepts the following:
        * - component - A react component to display in the Right-Hand Sidebar.
        * - title - A string or JSX element to display as a title for the RHS.
        * Returns:
        * - id: a unique identifier
        * - showRHSPlugin: the action to dispatch that will open the RHS.
        * - hideRHSPlugin: the action to dispatch that will close the RHS
        * - toggleRHSPlugin: the action to dispatch that will toggle the RHS
    */
    registerRightHandSidebarComponent(
        ...args: [
            component: ReactResolvable,
            title: ReactResolvable
        ] | [{
            component: ReactResolvable;
            title: ReactResolvable;
        }]
    ): {
        id: UniqueIdentifier;
        showRHSPlugin: object;
        hideRHSPlugin: object;
        toggleRHSPlugin: object;
    };

    /**
        * Register a Needs Team component by providing a route past /:team/:pluginId/ to be displayed at.
        * Accepts the following:
        * - route - The route to be displayed at.
        * - component - A react component to display.
        * Returns:
        * - id: a unique identifier
    */
    registerNeedsTeamRoute(
        ...args: [
            route: string,
            component: ReactResolvable
        ] | [{
            route: string;
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    registerCustomRoute(
        ...args: [
            route: string,
                    component: ReactResolvable
        ] | [{
            route: string;
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    registerProduct(
        ...args: [
            baseURL: string,
            switcherIcon: string,
            switcherText: string,
            switcherLinkURL: string,
            mainComponent: ReactResolvable,
            headerCentreComponent: ReactResolvable,
            headerRightComponent?: ReactResolvable,
            showTeamSidebar: boolean
        ] | [{
            baseURL: string;
            switcherIcon: string;
            switcherText: string;
            switcherLinkURL: string;
            mainComponent: ReactResolvable;
            headerCentreComponent: ReactResolvable;
            headerRightComponent?: ReactResolvable;
            showTeamSidebar: boolean;
        }]
    ): UniqueIdentifier;

    /**
        * Register a hook that will be called when a message is edited by the user before it
        * is sent to the server. Accepts a function that receives the post as an argument.
        *
        * To reject a post, return an object containing an error such as
        *     {error: {message: 'Rejected'}}
        * To modify or allow the post without modification, return an object containing the post
        * such as
        *     {post: {...}}
        *
        * If the hook function is asynchronous, the message will not be sent to the server
        * until the hook returns.
    */
    registerMessageWillBeUpdatedHook(
        ...args: [
            hook: (post: Partial<Post>, oldPost: Post) => Promise<{ error: { message: string } } | { post: Post }>
        ] | [{
            hook: (post: Partial<Post>, oldPost: Post) => Promise<{ error: { message: string } } | { post: Post }>;
        }]
    ): UniqueIdentifier;

    /**
        * INTERNAL: Subject to change without notice.
        * Register a component to render in the LHS next to a channel's link label.
        * All parameters are required.
        * Returns a unique identifier.
    */
    registerSidebarChannelLinkLabelComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * INTERNAL: Subject to change without notice.
        * Register a component to render in channel's center view, in place of a channel toast.
        * All parameters are required.
        * Returns a unique identifier.
    */
    registerChannelToastComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    /**
        * INTERNAL: Subject to change without notice.
        * Register a global component at the root of the app that survives across product switches.
        * All parameters are required.
        * Returns a unique identifier.
    */
    registerGlobalComponent(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    registerAppBarComponent(
        ...args: [
            iconUrl: string,
            action: PluginComponent['action'],
            tooltipText: ReactResolvable,
            supportedProductIds: ProductScope,
        ] | [{
            iconUrl: string;
            action: PluginComponent['action'];
            tooltipText: ReactResolvable;
            supportedProductIds: ProductScope;
        }]
    ): UniqueIdentifier;
    registerAppBarComponent(
        ...args: [
            iconUrl: string,
            action: undefined,
            tooltipText: ReactResolvable,
            supportedProductIds: ProductScope,
            rhsComponent: PluginComponent,
            rhsTitle: ReactResolvable,
        ] | [{
            iconUrl: string;
            action: undefined;
            tooltipText: ReactResolvable;
            supportedProductIds: ProductScope;
            rhsComponent: PluginComponent;
            rhsTitle: ReactResolvable;
        }]
    ): {
        id: UniqueIdentifier;
        component: ReturnType<PluginRegistry['registerRightHandSidebarComponent']>;
    };

    /**
        * INTERNAL: Subject to change without notice.
        * Register a handler to retrieve stats that will be displayed on the system console
        * Accepts the following:
        * - handler - Func to be called to retrieve the stats from plugin api. It must be type PluginSiteStatsHandler.
        * Returns undefined
    */
    registerSiteStatisticsHandler(
        ...args: [
            handler: PluginSiteStatsHandler
        ] | [{
            handler: PluginSiteStatsHandler;
        }]
    ): void;

    /**
        * Register a hook to intercept desktop notifications before they occur.
        * Accepts a function to run before the desktop notification is triggered.
        * The function has the following signature:
        *   (post: Post, msgProps: NewPostMessageProps, channel: Channel,
        *    teamId: string, args: DesktopNotificationArgs) => Promise<{
        *         error?: string;
        *         args?: DesktopNotificationArgs;
        *     }>)
        *
        * DesktopNotificationArgs is the following type:
        *   export type DesktopNotificationArgs = {
        *     title: string;
        *     body: string;
        *     silent: boolean;
        *     soundName: string;
        *     url: string;
        *     notify: boolean;
        * };
        *
        * To stop a desktop notification and allow subsequent hooks to process the notification, return:
        *   {args: {...args, notify: false}}
        * To enable a desktop notification and allow subsequent hooks to process the notification, return:
        *   {args: {...args, notify: true}}
        * To stop a desktop notification and prevent subsequent hooks from processing the notification, return either:
        *   {error: 'log this error'}, or {}
        * To allow subsequent hooks to process the notification, return:
        *   {args}, or null or undefined (thanks js)
        *
        * The args returned by the hook will be used as the args for the next hook, until all hooks are
        * completed. The resulting args will be used as the arguments for the `notifyMe` function.
        *
        * Returns a unique identifier.
    */
    registerDesktopNotificationHook(
        ...args: [
            hook: (
                post: Post,
                msgProps: NewPostMessageProps,
                channel: Channel,
                teamId: string,
                args: DesktopNotificationArgs
            ) => Promise<{
                error?: string;
                args?: DesktopNotificationArgs;
            }>
        ] | [{
            hook: (
                post: Post,
                msgProps: NewPostMessageProps,
                channel: Channel,
                teamId: string,
                args: DesktopNotificationArgs
            ) => Promise<{
                error?: string;
                args?: DesktopNotificationArgs;
            }>;
        }]
    ): UniqueIdentifier;

    /**
        * Register a schema for user settings. This will show in the user settings modals
        * and all values will be stored in the preferences with cateogry pp_${pluginId} and
        * the name of the setting.
        *
        * The settings definition can be found in /src/types/plugins/user_settings.ts
        *
        * Malformed settings will be filtered out.
    */
    registerUserSettings(
        ...args: [
            settings: PluginConfiguration
        ] | [{
            settings: PluginConfiguration;
        }]
    ): void;

    /**
        * Register a component to be displayed in the System Console Groups table.
        * Accepts a React component. Returns a unique identifier.
    */
    registerSystemConsoleGroupTable(
        ...args: [
            component: ReactResolvable
        ] | [{
            component: ReactResolvable;
        }]
    ): UniqueIdentifier;

    // The most up-to-date list of methods can be found at https://developers.mattermost.com/extend/plugins/webapp/reference
}
