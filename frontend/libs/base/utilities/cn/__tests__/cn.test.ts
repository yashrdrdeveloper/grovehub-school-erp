import { cn } from '../cn';

describe('cn', () => {
  it('merges compatible utility classes', () => {
    expect(cn('px-2', 'px-4', undefined)).toBe('px-4');
  });
  it('applies conditional classes and resolves display conflicts', () => {
    const hidden = false;
    expect(cn('flex', hidden && 'hidden', { block: true })).toBe('block');
  });
});
