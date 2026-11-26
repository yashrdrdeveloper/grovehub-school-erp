import { Box } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Sheet } from '@Components/components/Sheet';
import { useAuditPage } from './hooks/useAuditPage';
import styles from './page.module.css';

export const AuditPage = () => {
  const state = useAuditPage();
  const rows = state.query.data?.items ?? [];
  return <Box className={styles.page} data-workspace-fill="true">
    <h1>{__('Audit logs')}</h1>
    <p>{__('Your activity')}</p>
    <Box className={styles.toolbar}>
      <input data-testid="audit-search" aria-label={__('Search')} placeholder={__('Search')} value={state.searchText} onChange={(event) => state.search(event.target.value)} />
      <Button onClick={() => void state.query.refetch()}>{__('Refresh')}</Button>
    </Box>
    {state.query.isPending && <p role="status">{__('Loading…')}</p>}
    {state.query.isError && <p role="alert">{__('Unable to load records. Check your access and try again.')}</p>}
    {state.query.isSuccess && rows.length === 0 && <p>{__('No records found')}</p>}
    <Box className={styles.list}>
      {rows.map((row) => <button className={styles.row} key={row.auditId} data-testid="audit-row" onClick={() => state.setSelected(row)}><strong>{row.operation}</strong><span>{row.source}</span><span>{new Date(row.occurredAt).toLocaleString()}</span></button>)}
    </Box>
    <Box className={styles.toolbar}><Button disabled={state.page === 1 || state.query.isFetching} onClick={() => state.setPage(state.page - 1)}>{__('Previous')}</Button><span>{state.page}</span><Button disabled={!state.query.data?.pageInfo.hasNext || state.query.isFetching} onClick={() => state.setPage(state.page + 1)}>{__('Next')}</Button></Box>
    {state.selected && <Sheet open title={__('Details')} dataTestId="audit-details" onOpenChange={(open) => { if (!open) state.setSelected(null); }}>
      <p><strong>{__('auditId')}</strong> {state.selected.auditId}</p>
      <p><strong>{__('source')}</strong> {state.selected.source}</p>
      <p><strong>{__('assetType')}</strong> {state.selected.assetType}</p>
      <p><strong>{__('assetId')}</strong> {state.selected.assetId}</p>
      <p><strong>{__('operation')}</strong> {state.selected.operation}</p>
      <p><strong>{__('occurredAt')}</strong> {state.selected.occurredAt}</p>
      <pre className={styles.payload}>{JSON.stringify(state.selected.data, null, 2)}</pre>
    </Sheet>}
  </Box>;
};
