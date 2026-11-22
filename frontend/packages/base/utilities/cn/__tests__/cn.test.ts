import { cn } from '../cn';

describe('cn', () => {
  it('merges compatible utility classes', () => {
    expect(cn('px-2', 'px-4', undefined)).toBe('px-4');
  });
  it('supports conditional classes', () => {
    expect(cn('flex', false && 'hidden', { block: true })).toBe('flex block');
  });
});
