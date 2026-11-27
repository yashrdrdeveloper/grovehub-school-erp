import { LOGIN_EMAIL_PATTERN } from '../constants/login.validation';

export const isValidLoginEmail = (value: string): boolean => LOGIN_EMAIL_PATTERN.test(value);
