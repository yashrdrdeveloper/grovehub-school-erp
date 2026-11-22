import { _toLower } from '@Base/utilities/toLower';

import type { GetStatusPillDisplayLabel } from './types/statusPill.getDisplayLabel';

/** Delegate underscore-separated backend status labels to Lodash normalization. */
export const getStatusPillDisplayLabel: GetStatusPillDisplayLabel = (value) => _toLower(value);
