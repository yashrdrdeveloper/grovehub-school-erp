import { Eye, EyeOff } from '@Components/icons';
import { getPasswordFieldVisibility } from '../passwordField.visibility';
describe('getPasswordFieldVisibility', () => {
  it('selects accessible visible-password presentation', () => {
    expect(getPasswordFieldVisibility(true)).toEqual({
      VisibilityIcon: EyeOff, toggleLabel: 'Hide password', passwordType: 'text',
    });
    expect(getPasswordFieldVisibility(false)).toEqual({
      VisibilityIcon: Eye, toggleLabel: 'Show password', passwordType: 'password',
    });
  });
});
