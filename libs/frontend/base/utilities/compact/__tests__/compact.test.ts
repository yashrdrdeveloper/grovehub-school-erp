import { _compact } from '../compact';

describe('_compact', () => {
  it('delegates the documented normalization to Lodash', () => {
    expect(_compact(['one', undefined, '', 'two'])).toEqual(['one', 'two']);
  });
});
