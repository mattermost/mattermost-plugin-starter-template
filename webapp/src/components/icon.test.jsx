import React from 'react';
import {shallow} from 'enzyme';

import Icon from './icon';

test('icon test', () => {
    const wrapper = shallow(<Icon/>);
    expect(wrapper).toMatchSnapshot();
});
