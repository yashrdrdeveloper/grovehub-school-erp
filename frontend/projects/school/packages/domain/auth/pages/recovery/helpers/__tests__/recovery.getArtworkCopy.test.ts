import { getRecoveryArtworkCopy } from '../recovery.getArtworkCopy';

describe('getRecoveryArtworkCopy', () => {
  it('uses invitation guidance for incomplete invitation links', () => {
    expect(getRecoveryArtworkCopy({ purpose: 'invitation', state: 'missing' }).message).toMatch(
      /school administrator/i,
    );
  });
  it('uses a distinct reset-password success message', () => {
    expect(getRecoveryArtworkCopy({ purpose: 'recovery', state: 'success' }).title).toBe(
      'Password updated!',
    );
  });
});
