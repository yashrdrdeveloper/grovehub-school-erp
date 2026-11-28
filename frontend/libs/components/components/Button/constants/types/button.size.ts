import type { BUTTON_SIZE } from '../button.size';

export type ButtonSize = (typeof BUTTON_SIZE)[keyof typeof BUTTON_SIZE];
