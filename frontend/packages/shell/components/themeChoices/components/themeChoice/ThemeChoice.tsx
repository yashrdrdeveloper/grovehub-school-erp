import { Button } from '@Components/components/Button';

import { getThemeChoiceClassName } from '../../helpers/themeChoices.getClassName';

import type { ThemeChoiceComponent } from './types/themeChoice.props';

/** One accessible appearance choice; no storage or media-query side effects. */
export const ThemeChoice: ThemeChoiceComponent = ({
  value,
  label,
  Icon,
  preference,
  onPreferenceChange,
}) => {
  const selected = preference === value;

  const handleChange = () => {
    onPreferenceChange(value);
  };

  return (
    <Button
      unstyled
      type="button"
      aria-pressed={selected}
      className={getThemeChoiceClassName(selected)}
      onClick={handleChange}
    >
      <Icon size={16} aria-hidden="true" />
      {label}
    </Button>
  );
};
