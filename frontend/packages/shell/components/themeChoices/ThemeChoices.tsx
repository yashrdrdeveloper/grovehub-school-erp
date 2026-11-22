import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { Monitor, Moon, Sun } from '@Components/icons';
import type { ThemePreference } from '@Components/hooks/theme';

import { getThemeChoiceClassName } from './helpers/themeChoices.getClassName';
import type { ThemeChoicesComponent, ThemeChoicesProps } from './types/themeChoices';

const choices: { value: ThemePreference; label: string; Icon: typeof Sun }[] = [
  { value: 'system', label: 'System', Icon: Monitor },
  { value: 'light', label: 'Light', Icon: Sun },
  { value: 'dark', label: 'Dark', Icon: Moon },
];

const ThemeChoice = ({ value, label, Icon, preference, onPreferenceChange }: {
  value: ThemePreference;
  label: string;
  Icon: typeof Sun;
} & ThemeChoicesProps) => {
  const handleChange = () => onPreferenceChange(value);
  return (
    <Button unstyled type="button" aria-pressed={preference === value}
      className={getThemeChoiceClassName(preference === value)}
      onClick={handleChange}><Icon size={16}/>{label}</Button>
  );
};

export const ThemeChoices: ThemeChoicesComponent = (props) => (
  <Box as="div" className="theme-choices" role="group" aria-label="Appearance">
    {choices.map(choice => <ThemeChoice key={choice.value} {...choice} {...props}/>)}
  </Box>
);
