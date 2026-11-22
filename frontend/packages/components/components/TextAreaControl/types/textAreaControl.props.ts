import type { TextareaHTMLAttributes, ForwardRefExoticComponent, PropsWithoutRef, RefAttributes } from 'react';
export interface TextAreaControlProps extends TextareaHTMLAttributes<HTMLTextAreaElement> {
  readonly dataTestId?: string;
}
export type TextAreaControlComponent = ForwardRefExoticComponent<PropsWithoutRef<TextAreaControlProps> & RefAttributes<HTMLTextAreaElement>>;
