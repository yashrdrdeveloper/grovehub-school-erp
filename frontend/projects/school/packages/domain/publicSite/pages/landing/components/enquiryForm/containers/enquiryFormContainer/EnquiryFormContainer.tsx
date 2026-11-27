import { EnquiryForm } from '../../enquiryForm';
import { useEnquiryForm } from '../../hooks/useEnquiryForm';
import type { EnquiryFormProps } from '../../types/enquiryForm';

export const EnquiryFormContainer = (props: EnquiryFormProps) => (
  <EnquiryForm {...props} model={useEnquiryForm()} />
);
