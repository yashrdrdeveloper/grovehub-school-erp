import { useCallback, useId, useState, type FormEvent } from 'react';

import { EMPTY_ENQUIRY } from '../../constants/enquiry.initial';
import { validateEnquiryContact } from '../../helpers/enquiry.validateContact';
import { getEnquiryMailto } from '../../helpers/enquiry.getMailto';
import type { EnquiryFormData } from '../../types/enquiryForm';
import type { UseEnquiryForm } from './types/useEnquiryForm';

export const useEnquiryForm: UseEnquiryForm = () => {
  const fieldId = useId();
  const [data, setData] = useState<EnquiryFormData>(EMPTY_ENQUIRY);
  const [step, setStep] = useState(0);
  const [error, setError] = useState('');
  const [ready, setReady] = useState(false);

  const patch = useCallback(
    <K extends keyof EnquiryFormData>(key: K, value: EnquiryFormData[K]) => {
      setData((previous) => ({ ...previous, [key]: value }));
      setError('');
    },
    [],
  );

  const next = useCallback(
    (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      const validationError = validateEnquiryContact(data);
      setError(validationError ?? '');
      if (!validationError) setStep(1);
    },
    [data],
  );

  const submit = useCallback(
    (event: FormEvent<HTMLFormElement>) => {
      event.preventDefault();
      if (!data.privacy) {
        setError('Please confirm you agree to be contacted about your enquiry.');
        return;
      }
      setError('');
      // No public Lead CRM endpoint exists. Do not claim an enquiry was submitted.
      setReady(true);
    },
    [data.privacy],
  );

  const edit = useCallback(() => {
    setReady(false);
    setStep(0);
  }, []);
  const back = useCallback(() => {
    setStep(0);
    setError('');
  }, []);

  return {
    fieldId,
    data,
    step,
    error,
    ready,
    mailto: getEnquiryMailto(data),
    patch,
    next,
    submit,
    edit,
    back,
  };
};
