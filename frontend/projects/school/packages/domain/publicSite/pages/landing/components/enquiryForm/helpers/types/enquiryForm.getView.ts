import type { EnquiryFormModel } from '../../types/enquiryForm';
import type { getEnquiryFormView } from '../enquiryForm.getView';
export type GetEnquiryFormView = (
  compact: boolean,
  model: EnquiryFormModel,
) => ReturnType<typeof getEnquiryFormView>;
