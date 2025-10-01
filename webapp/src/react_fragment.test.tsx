// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';

function generateSomething() {
    return <div>{'something'}</div>;
}

test('Can test React fragments', () => {
    expect(React.version).toEqual('17.0.2');
    expect(generateSomething()).toBeDefined();
});
