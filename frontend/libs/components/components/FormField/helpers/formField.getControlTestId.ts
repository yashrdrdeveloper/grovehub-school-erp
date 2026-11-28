import { _kebabCase } from '@Base/utilities/kebabCase';

import type { GetFormFieldControlTestId } from './types/formField.getControlTestId';

/** A stable field selector is explicit whenever possible, otherwise Lodash normalizes the label. */
export const getFormFieldControlTestId: GetFormFieldControlTestId = ({ label, id, dataTestId }) =>
  dataTestId ?? id ?? _kebabCase(label);
