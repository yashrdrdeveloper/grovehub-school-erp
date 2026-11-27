import { _replace } from '../replace';

describe('_replace', () => {
  it('replaces all occurrences when supplied a global regex', () => {
    expect(_replace('AWAITING_PARENT_SETUP', /_/g, ' ')).toBe('AWAITING PARENT SETUP');
  });
});
