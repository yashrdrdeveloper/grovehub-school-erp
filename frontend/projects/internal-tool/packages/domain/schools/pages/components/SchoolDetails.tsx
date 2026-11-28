import { Link } from '@Components/router';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import type { SchoolAction, SchoolRecord } from '../types/schools.types';

function actionLabel(action: SchoolAction, record: SchoolRecord): string {
 if (action.kind === 'rename') return __('Request rename');
 if (record.status === 'ACTIVE') return __('Request deactivation');
 return __('Request activation');
}

export function SchoolDetails({ record, close, requestChange, canViewHistory }: { record: SchoolRecord; canViewHistory: boolean; close: () => void; requestChange: (record: SchoolRecord, action: SchoolAction) => void }) {
 return <section className="action-card" data-testid="school-details">
  <header className="page-heading"><div><h2>{record.name}</h2><p className="muted">{record.environment.name}</p></div><Button variant={BUTTON_VARIANT.GHOST} onClick={close} dataTestId="school-details-close">{__('Close')}</Button></header>
  <dl className="school-detail-grid">
   <div><dt>{__('School ID')}</dt><dd>{record.schoolId}</dd></div>
   <div><dt>{__('School code')}</dt><dd>{record.code || '—'}</dd></div>
   <div><dt>{__('Registry status')}</dt><dd>{record.snapshotReady && record.status || __('Unavailable')}</dd></div>
   <div><dt>{__('Registry version')}</dt><dd>{record.snapshotReady && record.version || '—'}</dd></div>
  </dl>
  {canViewHistory && <Link to={'/audit/schools?schoolId='+encodeURIComponent(record.schoolId)} className="school-link" dataTestId="school-mutation-history">{__('View school mutation history')}</Link>}
  {!record.snapshotReady && <p role="alert" className="error">{__('Current school details are unavailable or inconsistent. Refresh before requesting changes.')}</p>}
  {record.snapshotReady && !record.actions.length && <p className="muted">{__('No school mutation request grants for this school and environment.')}</p>}
  {record.actions.map((action) => <div key={action.template.id} className="school-action"><div><h3>{actionLabel(action, record)}</h3><p className="muted small">{action.template.code} · v{action.template.version}</p></div><Button variant={BUTTON_VARIANT.OUTLINE} onClick={() => requestChange(record, action)} dataTestId={'school-request-'+action.template.id}>{actionLabel(action, record)}</Button></div>)}
  <p className="muted small">{__('Changes require independent approval and execution. Disabling a school invalidates existing sessions and credential proofs.')}</p>
 </section>;
}
