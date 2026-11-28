import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';

import { renderRecoveryContent } from '../recovery.content';

import type { RecoveryProps } from '../../types/recovery.props';

const props: RecoveryProps = {
  id: 'recovery-renderer',
  purpose: 'recovery',
  invitation: false,
  view: 'ready',
  slow: false,
  password: '',
  confirmation: '',
  visible: false,
  error: '',
  fieldErrors: {},
  onPasswordChange: () => {},
  onConfirmationChange: () => {},
  onToggleVisibility: () => {},
  onSubmit: (event) => event.preventDefault(),
  onRetry: () => {},
};

describe('renderRecoveryContent', () => {
  it('renders only the password form for a valid recovery state', () => {
    render(<MemoryRouter>{renderRecoveryContent({ props, formVisible: true })}</MemoryRouter>);

    expect(screen.getByLabelText('New password')).toBeTruthy();
    expect(screen.getByLabelText('Confirm password')).toBeTruthy();
    expect(screen.queryByRole('link', { name: 'Request new link' })).toBeNull();
  });

  it('renders recovery status rather than password inputs for expired links', () => {
    render(
      <MemoryRouter>
        {renderRecoveryContent({ props: { ...props, view: 'expired' }, formVisible: false })}
      </MemoryRouter>,
    );

    expect(screen.queryByLabelText('New password')).toBeNull();
    expect(screen.getByRole('link', { name: 'Request new link' })).toBeTruthy();
  });
});
