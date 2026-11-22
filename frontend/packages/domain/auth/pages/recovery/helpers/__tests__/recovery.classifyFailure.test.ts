import { ApiError } from '@Business/api';
import { classifyRecoveryFailure } from '../recovery.classifyFailure';

describe('classifyRecoveryFailure', () => {
  it('recognizes the server lock decision', () => {
    expect(classifyRecoveryFailure(new ApiError(429, 'IDENTITY_SETUP_ATTEMPTS_EXCEEDED', 'Locked')))
      .toBe('attempts-exhausted');
  });
  it('keeps network errors and service failures retryable', () => {
    expect(classifyRecoveryFailure(new Error('network error'))).toBe('retryable');
    expect(classifyRecoveryFailure(new ApiError(503, 'IDENTITY_UNAVAILABLE', 'Unavailable')))
      .toBe('retryable');
  });
  it('distinguishes invalid tokens from service problems', () => {
    expect(classifyRecoveryFailure(new ApiError(401, 'IDENTITY_TOKEN_INVALID', 'Invalid')))
      .toBe('expired');
  });
});
