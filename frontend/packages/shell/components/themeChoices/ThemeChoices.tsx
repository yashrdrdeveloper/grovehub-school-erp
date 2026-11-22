import { Monitor, Moon, Sun } from 'lucide-react';
import type { ThemePreference } from '@Components/hooks/theme';

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
    <button type="button" aria-pressed={preference === value}
      className={preference === value ? 'theme-choice theme-choice--active' : 'theme-choice'}
      onClick={handleChange}><Icon size={16}/>{label}</button>
  );
};

export const ThemeChoices: ThemeChoicesComponent = (props) => (
  <div className="theme-choices" role="group" aria-label="Appearance">
    {choices.map(choice => <ThemeChoice key={choice.value} {...choice} {...props}/>)}
  </div>
);
