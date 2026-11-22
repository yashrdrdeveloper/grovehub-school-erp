import type { EnquiryFormData } from '../types/enquiryForm';

export const EMPTY_ENQUIRY: EnquiryFormData = {
  firstName: '', middleName: '', lastName: '', phone: '', email: '',
  type: 'ADMISSION', interestedClass: '', message: '', privacy: false,
};
