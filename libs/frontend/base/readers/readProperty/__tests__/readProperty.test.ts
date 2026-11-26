import { readProperty } from '../readProperty';

describe('readProperty', () => {
  it('selects the requested property with its original type', () => {
    const readName = readProperty<{ name: string; age: number }, 'name'>('name');
    expect(readName({ name: 'RDR', age: 12 })).toBe('RDR');
  });
});
