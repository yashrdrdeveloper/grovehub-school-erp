import { _includes } from '../includes';

describe('_includes', () => {
  it('uses exact membership semantics', () => {
    expect(_includes(['a', 'b'], 'a')).toBe(true);
    expect(_includes(['a', 'b'], 'c')).toBe(false);
  });
});
