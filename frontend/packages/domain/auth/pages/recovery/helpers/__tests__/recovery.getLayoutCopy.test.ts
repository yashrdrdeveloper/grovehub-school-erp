import { getRecoveryLayoutCopy } from '../recovery.getLayoutCopy';
describe('getRecoveryLayoutCopy', () => {
  it('shows invitation-specific titles and hides mobile form copy during checking', () => {
    expect(getRecoveryLayoutCopy({ invitation: true, formVisible: false }).mobileTitle).toBe(
      'Set password',
    );
    expect(getRecoveryLayoutCopy({ invitation: false, formVisible: true }).mobileTitle).toBe(
      'Reset password',
    );
    expect(getRecoveryLayoutCopy({ invitation: false, formVisible: false }).mobileDescription).toBe(
      '',
    );
  });
});
