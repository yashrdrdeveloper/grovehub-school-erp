import { fireEvent, render, screen } from '@testing-library/react';

import { DragList } from '../index';

describe('DragList', () => {
  it('moves a sort rule across multiple rows using a drag handle', () => {
    const onReorder = jest.fn();
    const items = ['Username', 'Status', 'Created at'];
    render(
      <DragList
        compact
        items={items}
        label="Sort priority"
        getItemId={(item) => item}
        onReorder={onReorder}
        renderItem={(item) => item}
      />,
    );
    const handles = screen.getAllByRole('button', { name: 'Reorder item; use up and down arrows' });
    const rows = screen.getAllByRole('listitem');
    const dataTransfer = { effectAllowed: '', setData: jest.fn() };
    fireEvent.dragStart(handles[2], { dataTransfer });
    fireEvent.dragOver(rows[0].firstElementChild!, { dataTransfer });
    fireEvent.drop(rows[0].firstElementChild!, { dataTransfer });
    expect(onReorder).toHaveBeenCalledWith(['Created at', 'Username', 'Status']);
  });

  it('supports keyboard-based ordering through the same handle', () => {
    const onReorder = jest.fn();
    render(
      <DragList
        items={['Username', 'Status']}
        label="Sort priority"
        getItemId={(item) => item}
        onReorder={onReorder}
        renderItem={(item) => item}
      />,
    );
    fireEvent.keyDown(screen.getAllByRole('button', { name: 'Reorder item; use up and down arrows' })[1], { key: 'ArrowUp' });
    expect(onReorder).toHaveBeenCalledWith(['Status', 'Username']);
  });
});
