import {id, version} from './manifest';

test('Plugin id and version are defined', () => {
    expect(id).toBeDefined();
    expect(version).toBeDefined();
});
