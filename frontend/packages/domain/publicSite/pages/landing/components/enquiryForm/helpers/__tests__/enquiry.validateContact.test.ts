import { EMPTY_ENQUIRY } from '../../constants/enquiry.initial';
import { validateEnquiryContact } from '../enquiry.validateContact';

describe('validateEnquiryContact', () => {
  it('requires the guardian name', () => {
    expect(validateEnquiryContact(EMPTY_ENQUIRY)).toBe('Please enter the guardian’s first name.');
  });
  it('requires one valid contact method', () => {
    expect(validateEnquiryContact({ ...EMPTY_ENQUIRY, firstName: 'RDR' })).toBe(
      'Enter a phone number or email address.',
    );
    expect(validateEnquiryContact({ ...EMPTY_ENQUIRY, firstName: 'RDR', email: 'bad' })).toBe(
      'Enter a valid email address.',
    );
    expect(validateEnquiryContact({ ...EMPTY_ENQUIRY, firstName: 'RDR', phone: '123' })).toBe(
      'Enter a valid phone number.',
    );
  });
  it('accepts an email without phone', () => {
    expect(
      validateEnquiryContact({ ...EMPTY_ENQUIRY, firstName: 'RDR', email: 'a@example.com' }),
    ).toBeNull();
  });
});
