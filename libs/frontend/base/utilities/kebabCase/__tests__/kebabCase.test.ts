import { _kebabCase } from '../kebabCase';

describe('_kebabCase', () => {
  it('delegates the documented normalization to Lodash', () => {
    expect(_kebabCase("Student's ID")).toBe('students-id');
  });
});
