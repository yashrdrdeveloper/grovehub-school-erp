import { __ } from '../__';

describe('__', () => {
  it('returns unparameterized source copy unchanged', () => {
    expect(__('Check your email')).toBe('Check your email');
  });

  it('interpolates named values including repeated placeholders', () => {
    expect(__('Hello {{name}}, welcome back {{name}}.', { name: 'Yash' })).toBe(
      'Hello Yash, welcome back Yash.',
    );
  });

  it('accepts numeric zero, boolean false, and optional placeholder whitespace', () => {
    expect(__('{{count}} active, {{ enabled }} enabled', { count: 0, enabled: false })).toBe(
      '0 active, false enabled',
    );
  });

  it('does not interpret placeholder-like content from inserted values', () => {
    expect(__('Hello {{name}}', { name: '{{unsafe}}' })).toBe('Hello {{unsafe}}');
  });

  it('rejects missing interpolation arguments rather than exposing placeholder tokens', () => {
    expect(() => __('Student {{student}}')).toThrow('Missing message parameter: student');
    expect(() => __('Student {{student}}', {})).toThrow('Missing message parameter: student');
  });

  it('does not resolve properties inherited through the prototype chain', () => {
    expect(() => __('{{toString}}', {})).toThrow('Missing message parameter: toString');
  });
});
