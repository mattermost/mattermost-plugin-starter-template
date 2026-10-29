// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

module.exports = {
    presets: [
        ['@babel/preset-env', {

            // Matches the browsers supported by the Mattermost web app.
            targets: {
                chrome: 110,
                firefox: 102,
                edge: 110,
                safari: '16.4',
            },
            corejs: 3,
            useBuiltIns: 'usage',
            shippedProposals: true,
        }],

        // The automatic runtime imports react/jsx-runtime, which webpack aliases to src/jsx_runtime.ts.
        ['@babel/preset-react', {
            runtime: 'automatic',
        }],
        ['@babel/preset-typescript', {
            allExtensions: true,
            isTSX: true,
        }],
    ],
};
