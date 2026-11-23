import { render, screen } from '@testing-library/react';

import { __ as resolveMessage } from '@Base/i18n';

import { InitI18nStatus } from '../../../hooks/useInitI18n';
import { I18nBootstrapFallback } from '../I18nBootstrapFallback';

describe('I18nBootstrapFallback', () => {
  afterEach(() => {
    window.__ = resolveMessage;
  });

  it('renders an accessible loading status without calling the global formatter', () => {
    Reflect.deleteProperty(window, '__');

    render(<I18nBootstrapFallback status={InitI18nStatus.INITIALIZING} />);

    expect(screen.getByRole('status').textContent).toContain('Loading school workspace');
  });

  it('renders an accessible failure instruction without calling the global formatter', () => {
    Reflect.deleteProperty(window, '__');

    render(<I18nBootstrapFallback status={InitI18nStatus.FAILED} />);

    expect(screen.getByRole('status').textContent).toContain('Please refresh the page');
  });
});
