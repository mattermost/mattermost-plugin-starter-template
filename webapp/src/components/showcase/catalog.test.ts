import {catalog} from './catalog';
import {countByCategory, filterCatalog} from './helpers';

describe('showcase catalog', () => {
    it('has unique component ids', () => {
        const ids = catalog.map((entry) => entry.id);
        expect(new Set(ids).size).toEqual(ids.length);
    });

    it('filters by name, description, and category', () => {
        const byName = filterCatalog(catalog, 'button', 'all');
        expect(byName.some((entry) => entry.id === 'button')).toBe(true);

        const byCategory = filterCatalog(catalog, '', 'forms');
        expect(byCategory.length).toBeGreaterThan(0);
        expect(byCategory.every((entry) => entry.category === 'forms')).toBe(true);

        const none = filterCatalog(catalog, 'zzzz-not-a-component', 'all');
        expect(none).toHaveLength(0);
    });

    it('counts matches per category including all', () => {
        expect(countByCategory(catalog, '', 'all')).toEqual(catalog.length);
        expect(countByCategory(catalog, '', 'actions')).toEqual(
            catalog.filter((entry) => entry.category === 'actions').length,
        );
    });
});
