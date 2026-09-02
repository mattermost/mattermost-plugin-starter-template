import type {ComponentType} from 'react';

export const CATEGORIES = [
    {id: 'all', label: 'All'},
    {id: 'actions', label: 'Actions'},
    {id: 'forms', label: 'Forms'},
    {id: 'feedback', label: 'Feedback'},
    {id: 'images', label: 'Images'},
    {id: 'layout', label: 'Layout'},
    {id: 'progress', label: 'Progress'},
    {id: 'status', label: 'Status'},
    {id: 'cards', label: 'Cards'},
    {id: 'messaging', label: 'Messaging'},
    {id: 'navigation', label: 'Navigation'},
    {id: 'patterns', label: 'Patterns'},
    {id: 'banners', label: 'Banners'},
    {id: 'chrome', label: 'Chrome'},
] as const;

export type CategoryId = typeof CATEGORIES[number]['id'];

export type CatalogEntry = {
    id: string;
    name: string;
    category: Exclude<CategoryId, 'all'>;
    description: string;
    layoutExcerpt?: boolean;
    Preview: ComponentType;
    Detail: ComponentType;
};
