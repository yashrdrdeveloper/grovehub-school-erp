import { fireEvent, render, screen } from '@testing-library/react';

import { Sun } from '@Components/icons';
import { ThemePreference } from '@Components/hooks/useWorkspaceTheme';

import { ThemeChoice } from '../ThemeChoice';

describe('ThemeChoice', () => {
  it('shows selected state and calls the supplied theme preference action', () => {
    const onPreferenceChange = jest.fn();

    render(
      <ThemeChoice
        value={ThemePreference.LIGHT}
        preference={ThemePreference.LIGHT}
        Icon={Sun}
        label="Light"
        onPreferenceChange={onPreferenceChange}
      />,
    );

    const button = screen.getByRole('button', { name: 'Light' });
    expect(button.getAttribute('aria-pressed')).toBe('true');

    fireEvent.click(button);
    expect(onPreferenceChange).toHaveBeenCalledWith(ThemePreference.LIGHT);
  });
});
