import { render, screen } from '@testing-library/react';
import { ExternalAnchor } from '../ExternalAnchor';

describe('ExternalAnchor', () => {
  it('forwards native semantics and a deterministic test ID', () => {
    render(
      <ExternalAnchor href="https://example.com" dataTestId="source-link">
        Source
      </ExternalAnchor>,
    );
    expect(screen.getByRole('link', { name: 'Source' })).toBe(screen.getByTestId('source-link'));
  });
});
