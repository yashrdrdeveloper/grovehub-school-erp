import { render, screen } from '@testing-library/react';
import type { AuditActivity } from '@Business/domain/platform/audit/types/contracts';
import { renderAuditMobileRow } from '../audit.mobileRow';

const activity: AuditActivity = {
  auditId: 'audit-1',
  source: 'IDENTITY',
  assetType: 'ACCOUNT',
  assetId: 'account-1',
  operation: 'CREATED',
  occurredAt: '2026-11-27T10:00:00Z',
  data: {},
};

describe('renderAuditMobileRow', () => {
  it('renders the row summary without an independent action', () => {
    render(renderAuditMobileRow(activity));

    expect(screen.getByText(activity.operation)).toBeTruthy();
    expect(screen.getByText(activity.source)).toBeTruthy();
    expect(screen.getByText(activity.assetType)).toBeTruthy();
    expect(screen.queryByRole('button')).toBeNull();
  });
});
