import type { GetFormFieldDescriptionId } from './types/formField.getDescriptionId';

/** The native description target only exists when help or validation text exists. */
export const getFormFieldDescriptionId: GetFormFieldDescriptionId = (id, error, hint) => {
  if (error || hint) {
    return `${id}-description`;
  }

  return undefined;
};
