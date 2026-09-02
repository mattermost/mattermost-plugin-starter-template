import React, {useState} from 'react';

import {Checkbox} from '@mattermost/compass-ui/components/checkbox';
import {Chip} from '@mattermost/compass-ui/components/chip';
import {Combobox} from '@mattermost/compass-ui/components/combobox';
import {DateRangePicker} from '@mattermost/compass-ui/components/date-range-picker';
import {Dropdown} from '@mattermost/compass-ui/components/dropdown';
import {Radio} from '@mattermost/compass-ui/components/radio';
import {SearchInput} from '@mattermost/compass-ui/components/search-input';
import {Select} from '@mattermost/compass-ui/components/select';
import {Switch} from '@mattermost/compass-ui/components/switch';
import {TextArea} from '@mattermost/compass-ui/components/text-area';
import {TextInput} from '@mattermost/compass-ui/components/text-input';

import {SpecimenRow} from '../helpers';
import type {CatalogEntry} from '../types';

const TEAM_OPTIONS = [
    {value: 'alpha', label: 'Alpha'},
    {value: 'bravo', label: 'Bravo'},
    {value: 'charlie', label: 'Charlie'},
];

function CheckboxPreview() {
    return <Checkbox defaultChecked={true}>{'Notify me'}</Checkbox>;
}

function CheckboxDetail() {
    return (
        <>
            <SpecimenRow label='States'>
                <Checkbox>{'Unchecked'}</Checkbox>
                <Checkbox defaultChecked={true}>{'Checked'}</Checkbox>
                <Checkbox indeterminate={true}>{'Indeterminate'}</Checkbox>
                <Checkbox disabled={true}>{'Disabled'}</Checkbox>
            </SpecimenRow>
            <SpecimenRow label='Size'>
                <Checkbox size='small'>{'Small'}</Checkbox>
                <Checkbox size='medium'>{'Medium'}</Checkbox>
                <Checkbox size='large'>{'Large'}</Checkbox>
            </SpecimenRow>
        </>
    );
}

function RadioPreview() {
    return (
        <Radio
            name='preview-radio'
            defaultChecked={true}
        >
            {'Email'}
        </Radio>
    );
}

function RadioDetail() {
    return (
        <SpecimenRow>
            <Radio
                name='detail-radio'
                defaultChecked={true}
            >
                {'Email'}
            </Radio>
            <Radio name='detail-radio'>{'Push'}</Radio>
            <Radio
                name='detail-radio'
                disabled={true}
            >
                {'SMS'}
            </Radio>
        </SpecimenRow>
    );
}

function SwitchPreview() {
    return <Switch defaultChecked={true}>{'Desktop'}</Switch>;
}

function SwitchDetail() {
    return (
        <>
            <SpecimenRow>
                <Switch>{'Off'}</Switch>
                <Switch defaultChecked={true}>{'On'}</Switch>
                <Switch disabled={true}>{'Disabled'}</Switch>
            </SpecimenRow>
            <SpecimenRow label='Size'>
                <Switch size='small'>{'Small'}</Switch>
                <Switch size='medium'>{'Medium'}</Switch>
                <Switch size='large'>{'Large'}</Switch>
            </SpecimenRow>
        </>
    );
}

function TextInputPreview() {
    return (
        <TextInput
            label='Channel name'
            defaultValue='Town Square'
        />
    );
}

function TextInputDetail() {
    return (
        <>
            <SpecimenRow>
                <TextInput label='Label'/>
            </SpecimenRow>
            <SpecimenRow>
                <TextInput
                    label='Invalid'
                    invalid={true}
                    defaultValue='oops'
                />
            </SpecimenRow>
            <SpecimenRow>
                <TextInput
                    label='Disabled'
                    disabled={true}
                />
            </SpecimenRow>
        </>
    );
}

function TextAreaPreview() {
    return (
        <TextArea
            label='Purpose'
            defaultValue='Share updates'
        />
    );
}

function TextAreaDetail() {
    return (
        <SpecimenRow>
            <TextArea
                label='Description'
                maxLength={80}
                showCharacterCount={true}
                defaultValue='Channel purpose'
            />
        </SpecimenRow>
    );
}

function SearchInputPreview() {
    return (
        <SearchInput
            label='Search'
            defaultValue='button'
        />
    );
}

function SearchInputDetail() {
    return (
        <>
            <SpecimenRow>
                <SearchInput label='Search components'/>
            </SpecimenRow>
            <SpecimenRow label='Size'>
                <SearchInput
                    size='small'
                    label='Small'
                />
                <SearchInput
                    size='medium'
                    label='Medium'
                />
            </SpecimenRow>
        </>
    );
}

function SelectPreview() {
    return (
        <Select
            label='Team'
            options={TEAM_OPTIONS}
            defaultValue='alpha'
        />
    );
}

function SelectDetail() {
    return (
        <SpecimenRow>
            <Select
                label='Team'
                placeholder='Choose a team'
                options={TEAM_OPTIONS}
            />
        </SpecimenRow>
    );
}

function ComboboxPreview() {
    return (
        <Combobox
            label='Channel'
            options={TEAM_OPTIONS}
            defaultValue='bravo'
        />
    );
}

function ComboboxDetail() {
    return (
        <SpecimenRow>
            <Combobox
                label='People'
                placeholder='Search'
                options={TEAM_OPTIONS}
                multiple={true}
            />
        </SpecimenRow>
    );
}

function ChipPreview() {
    return <Chip>{'Design'}</Chip>;
}

function ChipDetail() {
    return (
        <SpecimenRow>
            <Chip>{'Default'}</Chip>
            <Chip colored={true}>{'Colored'}</Chip>
            <Chip onRemove={() => undefined}>{'Removable'}</Chip>
            <Chip size='small'>{'Small'}</Chip>
        </SpecimenRow>
    );
}

function DropdownPreview() {
    return <Dropdown>{'Town Square'}</Dropdown>;
}

function DropdownDetail() {
    return (
        <SpecimenRow>
            <Dropdown>{'Closed'}</Dropdown>
            <Dropdown isOpen={true}>{'Open'}</Dropdown>
        </SpecimenRow>
    );
}

function DateRangePickerPreview() {
    return <DateRangePicker value='2026-09-02'/>;
}

function DateRangePickerDetail() {
    const [date, setDate] = useState('2026-09-02');

    return (
        <SpecimenRow>
            <DateRangePicker
                value={date}
                onChange={setDate}
            />
        </SpecimenRow>
    );
}

export const formEntries: CatalogEntry[] = [
    {
        id: 'checkbox',
        name: 'Checkbox',
        category: 'forms',
        description: 'Independent on/off choices, including indeterminate.',
        Preview: CheckboxPreview,
        Detail: CheckboxDetail,
    },
    {
        id: 'radio',
        name: 'Radio',
        category: 'forms',
        description: 'Pick exactly one option from a small set.',
        Preview: RadioPreview,
        Detail: RadioDetail,
    },
    {
        id: 'switch',
        name: 'Switch',
        category: 'forms',
        description: 'Immediate binary settings that take effect on toggle.',
        Preview: SwitchPreview,
        Detail: SwitchDetail,
    },
    {
        id: 'text-input',
        name: 'Text Input',
        category: 'forms',
        description: 'Single-line text with a floating label.',
        Preview: TextInputPreview,
        Detail: TextInputDetail,
    },
    {
        id: 'text-area',
        name: 'Text Area',
        category: 'forms',
        description: 'Multi-line text with optional character count.',
        Preview: TextAreaPreview,
        Detail: TextAreaDetail,
    },
    {
        id: 'search-input',
        name: 'Search Input',
        category: 'forms',
        description: 'Text field with a leading magnifier and clear control.',
        Preview: SearchInputPreview,
        Detail: SearchInputDetail,
    },
    {
        id: 'select',
        name: 'Select',
        category: 'forms',
        description: 'Single-select menu with a labeled trigger.',
        Preview: SelectPreview,
        Detail: SelectDetail,
    },
    {
        id: 'combobox',
        name: 'Combobox',
        category: 'forms',
        description: 'Searchable select, including multi-select chips.',
        Preview: ComboboxPreview,
        Detail: ComboboxDetail,
    },
    {
        id: 'chip',
        name: 'Chip',
        category: 'forms',
        description: 'Compact tokens for a selected person, filter, or value.',
        Preview: ChipPreview,
        Detail: ChipDetail,
    },
    {
        id: 'dropdown',
        name: 'Dropdown',
        category: 'forms',
        description: 'Heading-weight trigger for a popover menu.',
        Preview: DropdownPreview,
        Detail: DropdownDetail,
    },
    {
        id: 'date-range-picker',
        name: 'Date Range Picker',
        category: 'forms',
        description: 'Single date or range selection from a calendar.',
        Preview: DateRangePickerPreview,
        Detail: DateRangePickerDetail,
    },
];
