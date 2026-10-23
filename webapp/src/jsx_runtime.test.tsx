// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import {render, screen} from '@testing-library/react';
import {Fragment, jsx, jsxs} from 'jsx_runtime';
import {createRef} from 'react';

test('renders children, keys and refs like the React runtime without warnings', () => {
    const consoleError = jest.spyOn(console, 'error');
    const ref = createRef<HTMLUListElement>();
    const items = ['one', 'two'].map((label) => jsx('li', {children: label}, label));

    render(jsxs(Fragment, {
        children: [
            jsx('h1', {children: 'Title'}),
            jsx('ul', {ref, 'aria-label': 'Items', children: items}),
        ],
    }));

    expect(screen.getByRole('heading', {name: 'Title'})).toBeInTheDocument();
    expect(screen.getAllByRole('listitem').map((li) => li.textContent)).toEqual(['one', 'two']);
    expect(ref.current).toBe(screen.getByRole('list', {name: 'Items'}));
    expect(items.map((item) => item.key)).toEqual(['one', 'two']);
    expect(consoleError).not.toHaveBeenCalled();
});
