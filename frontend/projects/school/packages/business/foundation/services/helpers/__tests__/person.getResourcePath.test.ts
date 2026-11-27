import { getPersonResourcePath } from '../person.getResourcePath';

describe('getPersonResourcePath', () => {
  it('escapes a person identifier while preserving the protected endpoint', () => {
    expect(getPersonResourcePath('person/a')).toBe('/foundation/protected/persons/person%2Fa');
  });
});
