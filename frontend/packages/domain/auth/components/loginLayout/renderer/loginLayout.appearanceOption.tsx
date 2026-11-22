import type { ReactElement } from 'react';
import { LOGIN_LAYOUT_TEST_IDS } from '../constants/loginLayout.testIds';
import type { RenderLoginLayoutAppearanceOption } from './types/loginLayout.appearanceOption';

/** A single selectable appearance option; no selector logic in the layout JSX. */
export const renderLoginLayoutAppearanceOption: RenderLoginLayoutAppearanceOption = ({ value, label, Icon, preference, setPreference }): ReactElement => {
  let className = '';
  if (preference === value) className = 'is-active';
  const onSelect = () => setPreference(value);
  return (
    <button key={value} type="button" className={className} data-testid={LOGIN_LAYOUT_TEST_IDS.APPEARANCE_PREFIX + value}
      aria-label={label + ' theme'} aria-pressed={preference === value}
      title={label + ' theme'} onClick={onSelect}>
      <Icon size={16}/>
    </button>
  );
};
