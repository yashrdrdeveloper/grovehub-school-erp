import { _filter } from '@Base/utilities/filter';
import { _find } from '@Base/utilities/find';
import { enquiryTypes, school } from '../../../constants/school.content';
import type { GetEnquiryMailto } from './types/enquiry.getMailto';

/** Prepares an email draft only; it never claims a CRM submission. */
export const getEnquiryMailto: GetEnquiryMailto = (data) => {
  const enquiry = _find(enquiryTypes, (type) => type.value === data.type)?.label ?? data.type;
  const subject = encodeURIComponent(`RDR School — ${enquiry}`);
  const body = encodeURIComponent(
    [
      `Guardian: ${_filter([data.firstName, data.middleName, data.lastName], Boolean).join(' ')}`,
      `Phone: ${data.phone || 'Not provided'}`,
      `Email: ${data.email || 'Not provided'}`,
      `Enquiry: ${enquiry}`,
      `Interested class: ${data.interestedClass || 'Not provided'}`,
      `Message: ${data.message || '—'}`,
    ].join('\n'),
  );
  return `mailto:${school.email}?subject=${subject}&body=${body}`;
};
