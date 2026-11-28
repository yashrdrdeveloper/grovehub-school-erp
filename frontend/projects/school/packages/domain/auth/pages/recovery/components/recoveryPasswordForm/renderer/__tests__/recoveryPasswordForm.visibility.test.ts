import { Eye, EyeOff } from '@Components/icons';
import { getRecoveryPasswordVisibility } from '../recoveryPasswordForm.visibility';
describe('getRecoveryPasswordVisibility', () => {
  it('uses matching visibility labels and input types', () => {
    expect(getRecoveryPasswordVisibility(false).VisibilityIcon).toBe(Eye);
    expect(getRecoveryPasswordVisibility(false).confirmationLabel).toBe(
      'Show confirmation password',
    );
    expect(getRecoveryPasswordVisibility(true).VisibilityIcon).toBe(EyeOff);
    expect(getRecoveryPasswordVisibility(true).passwordType).toBe('text');
  });
});
