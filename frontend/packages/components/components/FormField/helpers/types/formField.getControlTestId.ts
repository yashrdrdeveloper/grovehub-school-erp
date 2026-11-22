export interface FormFieldControlTestIdInput {
  readonly label: string;
  readonly id?: string;
  readonly dataTestId?: string;
}

export type GetFormFieldControlTestId = (input: FormFieldControlTestIdInput) => string;
