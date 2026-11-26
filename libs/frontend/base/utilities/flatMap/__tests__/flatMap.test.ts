import { _flatMap } from '../flatMap';

describe('_flatMap', () => {
  it('preserves order while flattening one level of mapped collections', () => {
    expect(_flatMap([1, 2], (value) => [value, value + 10])).toEqual([1, 11, 2, 12]);
  });
});
