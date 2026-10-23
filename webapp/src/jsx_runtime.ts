// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

// Mattermost 11.7 doesn't expose window.ReactJSXRuntime, and React's own runtime only works with
// the React major it ships with (18 on 11.7, 19 on 12.0), so JSX goes through the host's
// createElement instead. Once min_server_version is 12.0, delete this file and its alias in
// webpack.config.js, and map react/jsx-runtime to ReactJSXRuntime in externals.
import {createElement, Fragment} from 'react';
import type {ElementType, Key, ReactNode} from 'react';

type Props = {children?: ReactNode; [prop: string]: unknown};

const withKey = (props: object, key?: Key) => (key === undefined ? props : {...props, key});

export function jsx(type: ElementType, props: Props, key?: Key) {
    return createElement(type, withKey(props, key));
}

// Static children are passed as arguments so React doesn't warn that they need keys.
export function jsxs(type: ElementType, {children, ...props}: Props, key?: Key) {
    return createElement(type, withKey(props, key), ...(children as ReactNode[]));
}

export function jsxDEV(type: ElementType, props: Props, key?: Key, isStaticChildren?: boolean) {
    return isStaticChildren ? jsxs(type, props, key) : jsx(type, props, key);
}

export {Fragment};
