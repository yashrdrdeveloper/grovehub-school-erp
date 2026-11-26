import { fireEvent, render, screen } from '@testing-library/react';

import { Tooltip } from '../Tooltip';

describe('Tooltip', () => {
  it('forwards the test ID to the trigger wrapper', () => {
    render(
      <Tooltip label="Refresh" dataTestId="refresh-tooltip">
        <button type="button">Refresh table</button>
      </Tooltip>,
    );

    expect(screen.getByTestId('refresh-tooltip').contains(screen.getByRole('button'))).toBe(true);
  });

  it('opens on hover above the trigger and dismisses on pointer down and mouse leave', () => {
    render(
      <Tooltip label="Refresh">
        <button type="button">Refresh table</button>
      </Tooltip>,
    );

    const button = screen.getByRole('button', { name: 'Refresh table' });
    expect(screen.queryByRole('tooltip')).toBeNull();

    fireEvent.mouseEnter(button);
    expect(screen.getByRole('tooltip').textContent).toBe('Refresh');
    expect(button.parentElement?.getAttribute('data-tooltip-hover')).toBe('true');
    fireEvent.scroll(window);
    expect(screen.getByRole('tooltip').textContent).toBe('Refresh');
    expect(screen.getByRole('tooltip').style.top).toBe('-6px');

    fireEvent.pointerDown(button);
    fireEvent.click(button);
    expect(screen.queryByRole('tooltip')).toBeNull();

    fireEvent.mouseLeave(button);
    fireEvent.mouseEnter(button);
    expect(screen.getByRole('tooltip')).toBeTruthy();
    fireEvent.mouseLeave(button);
    expect(screen.queryByRole('tooltip')).toBeNull();
  });
});
