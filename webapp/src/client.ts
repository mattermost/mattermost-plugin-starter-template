// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import manifest from 'manifest';

import type {Options} from '@mattermost/types/client4';

import {Client4} from 'mattermost-redux/client';

let basePath = '';

// Keep only the path of the SiteURL so requests stay same-origin even when the
// browser reached the server through a different hostname.
export function setSiteURL(siteURL: string | undefined) {
    basePath = siteURL ? new URL(siteURL).pathname.replace(/\/$/, '') : '';
}

async function doFetch<T>(path: string, options: Options = {}): Promise<T> {
    // getOptions adds X-Requested-With and, for non-GET requests, the session's CSRF
    // token. Without the token the server forwards the request as unauthenticated.
    const response = await fetch(
        `${basePath}/plugins/${manifest.id}/api/v1${path}`,
        Client4.getOptions(options),
    );

    const body = await response.json().catch(() => null);
    if (!response.ok) {
        throw new Error(body?.error ?? `Request failed with status ${response.status}`);
    }
    return body as T;
}

export async function getGreetingCount(channelId: string): Promise<number> {
    const {count} = await doFetch<{count: number}>(`/channels/${encodeURIComponent(channelId)}/greetings`);
    return count;
}

export async function resetGreetingCount(channelId: string): Promise<void> {
    await doFetch(`/channels/${encodeURIComponent(channelId)}/greetings`, {method: 'delete'});
}
