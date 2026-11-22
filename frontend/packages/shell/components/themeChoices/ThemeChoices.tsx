import { Box } from '@Components/components/Box';

import { THEME_CHOICES } from './constants/themeChoices.options';
import { ThemeChoice } from './components/themeChoice';

import type { ThemeChoicesComponent } from './types/themeChoices';

/** Theme option rendering stays stateless; one controlled preference is passed down. */
export const ThemeChoices: ThemeChoicesComponent = (props) => (
  <Box as="div" className="theme-choices" role="group" aria-label="Appearance">
    {THEME_CHOICES.map((choice) => (
      <ThemeChoice key={choice.value} {...choice} {...props} />
    ))}
  </Box>
);
