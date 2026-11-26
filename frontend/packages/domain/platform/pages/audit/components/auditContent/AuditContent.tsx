import { Box } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Sheet } from '@Components/components/Sheet';
import type { AuditControls } from '../../hooks/useAuditPage/types/useAuditPage';
import { AUDIT_COPY } from '../../constants/audit.copy';
import { AUDIT_TEST_IDS } from '../../constants/audit.testIds';
import styles from '../../page.module.css';

export const AuditContent = ({ controls: state }: { controls: AuditControls }) => {
  const rows = state.query.data?.items ?? [];
  return <Box className={styles.page} data-workspace-fill="true">
    <h1>{AUDIT_COPY.label1}</h1>
    <p>{AUDIT_COPY.label2}</p>
    <Box className={styles.toolbar}>
      <input data-testid={AUDIT_TEST_IDS.search} aria-label={AUDIT_COPY.label3} placeholder={AUDIT_COPY.label3} value={state.searchText} onChange={(event) => state.search(event.target.value)} />
      <Button onClick={() => void state.query.refetch()}>{AUDIT_COPY.label4}</Button>
    </Box>
    {state.query.isPending && <p role="status">{AUDIT_COPY.label5}</p>}
    {state.query.isError && <p role="alert">{AUDIT_COPY.label6}</p>}
    {state.query.isSuccess && rows.length === 0 && <p>{AUDIT_COPY.label7}</p>}
    <Box className={styles.list}>
      {rows.map((row) => <button className={styles.row} key={row.auditId} data-testid={AUDIT_TEST_IDS.row} onClick={() => state.setSelected(row)}><strong>{row.operation}</strong><span>{row.source}</span><span>{new Date(row.occurredAt).toLocaleString()}</span></button>)}
    </Box>
    <Box className={styles.toolbar}><Button disabled={state.page === 1 || state.query.isFetching} onClick={() => state.setPage(state.page - 1)}>{AUDIT_COPY.label8}</Button><span>{state.page}</span><Button disabled={!state.query.data?.pageInfo.hasNext || state.query.isFetching} onClick={() => state.setPage(state.page + 1)}>{AUDIT_COPY.label9}</Button></Box>
    {state.selected && <Sheet open title={AUDIT_COPY.label10} dataTestId={AUDIT_TEST_IDS.details} onOpenChange={(open) => { if (!open) state.setSelected(null); }}>
      <p><strong>{AUDIT_COPY.label11}</strong> {state.selected.auditId}</p>
      <p><strong>{AUDIT_COPY.label12}</strong> {state.selected.source}</p>
      <p><strong>{AUDIT_COPY.label13}</strong> {state.selected.assetType}</p>
      <p><strong>{AUDIT_COPY.label14}</strong> {state.selected.assetId}</p>
      <p><strong>{AUDIT_COPY.label15}</strong> {state.selected.operation}</p>
      <p><strong>{AUDIT_COPY.label16}</strong> {state.selected.occurredAt}</p>
      <pre className={styles.payload}>{JSON.stringify(state.selected.data, null, 2)}</pre>
    </Sheet>}
  </Box>;
};
