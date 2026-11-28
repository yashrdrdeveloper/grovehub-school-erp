import { _compact } from '@Base/utilities/compact';
import { _join } from '@Base/utilities/join';

import type { GetFormFieldCombinedDescription } from './types/formField.getCombinedDescription';

/** Assemble ARIA descriptions without empty separators or a duplicate feature implementation. */
export const getFormFieldCombinedDescription: GetFormFieldCombinedDescription = (
  describedBy,
  descriptionId,
) => {
  const value = _join(_compact([describedBy, descriptionId]), ' ');

  if (value.length > 0) {
    return value;
  }

  return undefined;
};
