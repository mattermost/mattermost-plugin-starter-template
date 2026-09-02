import {actionEntries} from './specimens/actions';
import {bannerEntries} from './specimens/banners';
import {cardEntries} from './specimens/cards';
import {chromeEntries} from './specimens/chrome';
import {feedbackEntries} from './specimens/feedback';
import {formEntries} from './specimens/forms';
import {imageEntries} from './specimens/images';
import {layoutEntries} from './specimens/layout';
import {messagingEntries} from './specimens/messaging';
import {navigationEntries} from './specimens/navigation';
import {patternEntries} from './specimens/patterns';
import {progressEntries} from './specimens/progress';
import {statusEntries} from './specimens/status';
import type {CatalogEntry} from './types';

export const catalog: CatalogEntry[] = [
    ...actionEntries,
    ...formEntries,
    ...feedbackEntries,
    ...imageEntries,
    ...layoutEntries,
    ...progressEntries,
    ...statusEntries,
    ...cardEntries,
    ...messagingEntries,
    ...navigationEntries,
    ...patternEntries,
    ...bannerEntries,
    ...chromeEntries,
];
