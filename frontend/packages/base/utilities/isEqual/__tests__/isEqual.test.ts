import { isEqual } from '../isEqual';

// lodash-es is an ESM package. Exercise our typed delegation contract without
// asking Jest's CommonJS transformer to parse its node_modules source.
const mockCompare = jest.fn<boolean, [unknown, unknown]>();
jest.mock('lodash-es/isEqual', () => ({
  __esModule: true,
  default: (left: unknown, right: unknown) => mockCompare(left, right),
}));

describe('isEqual', () => {
  beforeEach(() => { mockCompare.mockReset(); });

  it('passes both inputs through to the existing equality implementation', () => {
    const left = { members: ['one'] };
    const right = { members: ['one'] };
    mockCompare.mockReturnValue(true);

    expect(isEqual(left, right)).toBe(true);
    expect(mockCompare).toHaveBeenCalledWith(left, right);
  });
});
