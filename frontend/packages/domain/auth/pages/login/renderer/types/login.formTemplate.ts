import type { LoginProps } from '../../types';
import type { ReactElement } from 'react';

export type GetLoginFormTemplate = (isMobile: boolean) => (props: LoginProps) => ReactElement;
