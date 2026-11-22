import { EMPTY_ENQUIRY } from '../../../config/emptyEnquiry.config';
import { getEnquiryMailto } from '../getEnquiryMailto';

describe('getEnquiryMailto', () => {
  it('generates an encoded email draft with contact details', () => {
    const mailto = getEnquiryMailto({ ...EMPTY_ENQUIRY, firstName: 'Parent', email: 'parent@example.com' });
    expect(mailto).toMatch(/^mailto:/);
    expect(decodeURIComponent(mailto)).toContain('Guardian: Parent');
    expect(decodeURIComponent(mailto)).toContain('Email: parent@example.com');
  });
});
