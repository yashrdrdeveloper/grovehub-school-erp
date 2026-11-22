import { getAuthProofHeader } from '../auth.getProofHeader';

describe('getAuthProofHeader', () => {
  it('uses the expected backend recovery authorization scheme', () => {
    expect(getAuthProofHeader('RECOVERY', 'abc')).toEqual({
      Authorization: 'Recovery abc',
    });
  });

  it('keeps invitation credentials distinct from recovery credentials', () => {
    expect(getAuthProofHeader('INVITATION', 'abc')).toEqual({
      Authorization: 'Invitation abc',
    });
  });
});
