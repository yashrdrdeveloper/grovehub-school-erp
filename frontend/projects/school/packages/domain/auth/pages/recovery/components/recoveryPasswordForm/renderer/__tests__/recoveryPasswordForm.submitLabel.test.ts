import { getRecoveryPasswordSubmitLabel } from '../recoveryPasswordForm.submitLabel';
describe('getRecoveryPasswordSubmitLabel', () => {
  it('resolves saving versus default literal message copy', () => {
    expect(getRecoveryPasswordSubmitLabel(false)).toBe('Update password');
    expect(getRecoveryPasswordSubmitLabel(true)).toBe('Securing your account…');
  });
});
