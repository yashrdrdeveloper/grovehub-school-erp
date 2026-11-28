import { createShellQueryClient } from '../shellContainer.createQueryClient';

describe('createShellQueryClient', () => {
  it('returns an isolated client configured with shared options', () => {
    const first = createShellQueryClient();
    const second = createShellQueryClient();
    expect(first).not.toBe(second);
    expect(first.getDefaultOptions().queries?.staleTime).toBe(30000);
    expect(first.getDefaultOptions().queries?.retry).toBe(1);
    expect(first.getDefaultOptions().queries?.refetchOnWindowFocus).toBe(false);
  });
});
