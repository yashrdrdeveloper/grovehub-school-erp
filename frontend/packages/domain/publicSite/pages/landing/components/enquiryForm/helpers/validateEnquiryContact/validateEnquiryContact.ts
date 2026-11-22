import type { ValidateEnquiryContact } from './types/validateEnquiryContact';

export const validateEnquiryContact: ValidateEnquiryContact = data => {
  if (!data.firstName.trim()) return 'Please enter the guardian’s first name.';
  if (!data.phone.trim() && !data.email.trim()) return 'Enter a phone number or email address.';
  if (data.email.trim() && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(data.email.trim())) {
    return 'Enter a valid email address.';
  }
  if (data.phone.trim() && data.phone.replace(/\D/g, '').length < 10) {
    return 'Enter a valid phone number.';
  }
  return null;
};
