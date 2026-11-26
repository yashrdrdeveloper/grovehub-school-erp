import { _trim } from '../trim';

describe('_trim', () => {
  it('trims whitespace while preserving case', () => {
    expect(_trim('  Parent@School.edu  ')).toBe('Parent@School.edu');
  });
  it('does not modify interior spaces', () => {
    expect(_trim('a  b')).toBe('a  b');
  });
});
