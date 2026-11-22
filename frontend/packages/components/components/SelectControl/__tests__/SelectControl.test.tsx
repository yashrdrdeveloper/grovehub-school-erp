import { render, screen } from '@testing-library/react';
import { SelectControl } from '../SelectControl';

describe('SelectControl', () => {
  it('forwards native semantics and a deterministic test ID', () => {
    render(
      <SelectControl aria-label="Class" dataTestId="class-select">
        <option value="IX">IX</option>
      </SelectControl>,
    );
    expect(screen.getByRole('combobox', { name: 'Class' })).toBe(
      screen.getByTestId('class-select'),
    );
  });
});
