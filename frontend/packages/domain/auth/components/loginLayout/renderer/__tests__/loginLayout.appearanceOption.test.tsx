import { ThemePreference } from '@Components/hooks/useWorkspaceTheme';
import { fireEvent, render, screen } from '@testing-library/react';
import { Sun } from '@Components/icons';
import { renderLoginLayoutAppearanceOption } from '../loginLayout.appearanceOption';

describe('renderLoginLayoutAppearanceOption', () => {
  it('indicates selected appearance and delegates the selection event', () => {
    const select = jest.fn();
    render(<>{renderLoginLayoutAppearanceOption({ value:ThemePreference.LIGHT, label:'Light', Icon:Sun, preference:ThemePreference.LIGHT, setPreference:select })}</>);
    const button = screen.getByRole('button', {name:'Light theme'});
    expect(button.getAttribute('aria-pressed')).toBe('true');
    fireEvent.click(button);
    expect(select).toHaveBeenCalledWith(ThemePreference.LIGHT);
  });
});
