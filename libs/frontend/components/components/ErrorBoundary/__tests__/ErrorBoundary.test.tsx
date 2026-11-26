import { fireEvent, render, screen } from '@testing-library/react';

import { ErrorBoundary } from '../ErrorBoundary';

const FaultyView = ({ crash }: { crash: boolean }) => {
  if (crash) {
    throw new Error('Private diagnostics must not leak');
  }

  return <div>Ready</div>;
};

describe('ErrorBoundary', () => {
  const consoleError = jest.spyOn(console, 'error').mockImplementation(() => undefined);

  afterAll(() => {
    consoleError.mockRestore();
  });

  it('keeps failed page content isolated behind a safe accessible fallback', () => {
    const onError = jest.fn();
    render(
      <ErrorBoundary onError={onError}>
        <FaultyView crash />
      </ErrorBoundary>,
    );

    expect(screen.getByRole('alert')).toBeTruthy();
    expect(screen.getByText('This section could not be displayed.')).toBeTruthy();
    expect(screen.queryByText('Private diagnostics must not leak')).toBeNull();
    expect(onError).toHaveBeenCalledTimes(1);
  });

  it('resets on a new route key without losing the containing application', () => {
    const { rerender } = render(
      <ErrorBoundary resetKey="/one">
        <FaultyView crash />
      </ErrorBoundary>,
    );

    rerender(
      <ErrorBoundary resetKey="/two">
        <FaultyView crash={false} />
      </ErrorBoundary>,
    );

    expect(screen.getByText('Ready')).toBeTruthy();
  });

  it('provides an accessible local retry action', () => {
    const renderState = { broken: true };
    const Recoverable = () => <FaultyView crash={renderState.broken} />;
    render(
      <ErrorBoundary>
        <Recoverable />
      </ErrorBoundary>,
    );

    renderState.broken = false;
    fireEvent.click(screen.getByRole('button', { name: 'Try again' }));
    expect(screen.getByText('Ready')).toBeTruthy();
  });
});
