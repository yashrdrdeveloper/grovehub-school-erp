import { _replace } from '@Base/utilities/replace';
import { _toLower } from '@Base/utilities/toLower';

import type { GetStatusPillDisplayLabel } from './types/statusPill.getDisplayLabel';

/** Match the existing public label: replace underscores first, then lower-case via Lodash. */
export const getStatusPillDisplayLabel: GetStatusPillDisplayLabel = (value) =>
  _toLower(_replace(value, /_/g, ' '));
