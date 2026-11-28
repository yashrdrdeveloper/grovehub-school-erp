import { _join } from '../join';

describe('_join', () => {
  it('delegates the documented normalization to Lodash', () => {
    expect(_join(['first', 'second'], ' ')).toBe('first second');
  });
});
