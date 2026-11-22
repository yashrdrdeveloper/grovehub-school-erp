import { getRecoveryArtworkDescriptor } from '../recovery.getArtworkDescriptor';

describe('getRecoveryArtworkDescriptor', () => {
  it('assigns the success style and status unambiguously', () => {
    expect(getRecoveryArtworkDescriptor('success')).toMatchObject({statusKind:'success',visualClass:'success'});
  });
  it('preserves the exhausted and transient error distinction', () => {
    expect(getRecoveryArtworkDescriptor('attempts-exhausted').statusKind).toBe('attempts-exhausted');
    expect(getRecoveryArtworkDescriptor('unavailable').statusKind).toBe('unavailable');
  });
});
