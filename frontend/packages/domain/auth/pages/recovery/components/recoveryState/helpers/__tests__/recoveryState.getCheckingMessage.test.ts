import { getRecoveryCheckingMessage } from '../recoveryState.getCheckingMessage';
describe('getRecoveryCheckingMessage', () => {
  it('prioritizes slow verification over the invitation purpose', () => {
    expect(getRecoveryCheckingMessage({ invitation: true, slow: true })).toBe('Still checking. This may take a moment…');
    expect(getRecoveryCheckingMessage({ invitation: true, slow: false })).toBe('Verifying your invitation securely…');
    expect(getRecoveryCheckingMessage({ invitation: false, slow: false })).toBe('Verifying your reset link securely…');
  });
});
