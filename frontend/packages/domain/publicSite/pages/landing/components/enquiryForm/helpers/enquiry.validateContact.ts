import { _trim } from '@Base/utilities/trim';
import type { ValidateEnquiryContact } from './types/enquiry.validateContact';

export const validateEnquiryContact: ValidateEnquiryContact = data => {
  if (!_trim(data.firstName)) return 'Please enter the guardian’s first name.';
  if (!_trim(data.phone) && !_trim(data.email)) return 'Enter a phone number or email address.';
  if (_trim(data.email) && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(_trim(data.email))) {
    return 'Enter a valid email address.';
  }
  if (_trim(data.phone) && data.phone.replace(/\D/g, '').length < 10) {
    return 'Enter a valid phone number.';
  }
  return null;
};
