import type { FormEventHandler, ReactElement } from 'react';

export interface EnquiryFormData {
  firstName: string;
  middleName: string;
  lastName: string;
  phone: string;
  email: string;
  type: string;
  interestedClass: string;
  message: string;
  privacy: boolean;
}

export interface EnquiryFormProps {
  readonly compact?: boolean;
}

export interface EnquiryFormModel {
  readonly fieldId: string;
  readonly data: EnquiryFormData;
  readonly step: number;
  readonly error: string;
  readonly ready: boolean;
  readonly mailto: string;
  readonly patch: <K extends keyof EnquiryFormData>(key: K, value: EnquiryFormData[K]) => void;
  readonly next: FormEventHandler<HTMLFormElement>;
  readonly submit: FormEventHandler<HTMLFormElement>;
  readonly edit: () => void;
  readonly back: () => void;
}

export type EnquiryFormComponent = (
  props: EnquiryFormProps & { model: EnquiryFormModel },
) => ReactElement;
