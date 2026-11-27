import { ENQUIRY_TEST_IDS } from '../../../constants/enquiry.testIds';
import type { EnquiryFormModel } from '../types/enquiryForm';

/** Pure presentation view; no form mutations or ad-hoc JSX copy. */
export const getEnquiryFormView = (compact: boolean, model: EnquiryFormModel) => {
  const rootClass = 'rdr-enquiry' + ((compact && ' rdr-enquiry--compact') || '');
  const rootId = (!compact && 'enquiry') || undefined;
  const rootTestId = (compact && ENQUIRY_TEST_IDS.SHEET) || ENQUIRY_TEST_IDS.PAGE;
  const heading =
    (model.ready && __('Your enquiry is ready')) || __('Have a question? We’re here.');
  const description =
    (model.ready && __('Open your email application to send the enquiry to the school.')) ||
    __('A quick conversation is the first step. We’ll help you find the answers.');
  const handleSubmit = (model.step === 0 && model.next) || model.submit;
  const stepHeading = (model.step === 0 && __('Your details')) || __('A little more detail');
  const submitLabel = (model.step === 0 && __('Continue')) || __('Prepare enquiry');

  return {
    rootClass,
    rootId,
    rootTestId,
    heading,
    description,
    handleSubmit,
    stepHeading,
    submitLabel,
  };
};
