import { getEnquiryFormView } from '../enquiryForm.getView';
import type { EnquiryFormModel } from '../../types/enquiryForm';

const model = {
  step: 0,
  ready: false,
  next: jest.fn(),
  submit: jest.fn(),
} as unknown as EnquiryFormModel;

describe('getEnquiryFormView', () => {
  it('chooses the compact form and first-step action without mutation', () => {
    const view = getEnquiryFormView(true, model);
    expect(view.rootId).toBeUndefined();
    expect(view.handleSubmit).toBe(model.next);
    expect(view.heading).toBe('Have a question? We’re here.');
  });
  it('shows the final enquiry action and accurate unsubmitted-email guidance', () => {
    const view = getEnquiryFormView(false, { ...model, step: 1, ready: true });
    expect(view.rootId).toBe('enquiry');
    expect(view.handleSubmit).toBe(model.submit);
    expect(view.heading).toBe('Your enquiry is ready');
  });
});
