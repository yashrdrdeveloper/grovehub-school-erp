import { AuthProofPurpose } from '../../constants/auth.proofScheme';
import { getAuthProofHeader } from '../auth.getProofHeader';

describe('getAuthProofHeader', () => {
  it('uses the expected backend recovery authorization scheme', () => {
    expect(getAuthProofHeader(AuthProofPurpose.RECOVERY, 'abc')).toEqual({
      Authorization: 'Recovery abc',
    });
  });

  it('keeps invitation credentials distinct from recovery credentials', () => {
    expect(getAuthProofHeader(AuthProofPurpose.INVITATION, 'abc')).toEqual({
      Authorization: 'Invitation abc',
    });
  });
});
