import { _map } from '@Base/utilities/map';
import { Box, BOX_TAG } from '@Components/components/Box';

import styles from './ThemeChoices.module.css';

import { THEME_CHOICES } from './constants/themeChoices.options';
import { ThemeChoice } from './components/themeChoice';

import type { ThemeChoicesComponent } from './types/themeChoices';

/** Theme option rendering stays stateless; one controlled preference is passed down. */
export const ThemeChoices: ThemeChoicesComponent = (props) => (
  <Box
    as={BOX_TAG.DIV}
    className={styles.root}
    data-compact={props.compact}
    role="group"
    aria-label={__('Appearance')}
  >
    {_map(THEME_CHOICES, (choice) => (
      <ThemeChoice key={choice.value} {...choice} {...props} />
    ))}
  </Box>
);
