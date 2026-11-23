import { getLoginLayoutHeading } from '../loginLayout.getHeading';

describe('getLoginLayoutHeading', () => {
  it('selects the mobile title only when the viewport is mobile', () => {
    expect(getLoginLayoutHeading(false, 'Set a new password.', 'Set password')).toBe(
      'Set a new password.',
    );
    expect(getLoginLayoutHeading(true, 'Set a new password.', 'Set password')).toBe(
      'Set password',
    );
    expect(getLoginLayoutHeading(true, 'Welcome back.')).toBe('Welcome back.');
  });
});
