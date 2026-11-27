import type { ReactElement } from 'react';

import type { ThemeChoicesProps } from '../../../types/themeChoices';
import type { ThemeChoiceOption } from '../../../constants/themeChoices.options';

export interface ThemeChoiceProps extends ThemeChoiceOption, ThemeChoicesProps {}

export type ThemeChoiceComponent = (props: ThemeChoiceProps) => ReactElement;
