import type { BUTTON_VARIANT } from '../button.variant';

export type ButtonVariant = (typeof BUTTON_VARIANT)[keyof typeof BUTTON_VARIANT];
