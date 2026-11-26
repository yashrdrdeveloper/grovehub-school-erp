import { Box } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Sheet } from '@Components/components/Sheet';
import type { RbacControls } from '../../hooks/useRbacPage/types/useRbacPage';
import { RBAC_COPY } from '../../constants/rbac.copy';
import { RBAC_TEST_IDS } from '../../constants/rbac.testIds';
import styles from '../../page.module.css';

export const RbacContent = ({ controls: state }: { controls: RbacControls }) => {
  const rows = state.query.data?.items ?? [];
  return <Box className={styles.page} data-workspace-fill="true">
    <h1>{RBAC_COPY.label1}</h1>
    
    <Box className={styles.toolbar}>
      <input data-testid={RBAC_TEST_IDS.search} aria-label={RBAC_COPY.label2} placeholder={RBAC_COPY.label2} value={state.searchText} onChange={(event) => state.search(event.target.value)} />
      <Button onClick={() => void state.query.refetch()}>{RBAC_COPY.label3}</Button>
    </Box>
    {state.query.isPending && <p role="status">{RBAC_COPY.label4}</p>}
    {state.query.isError && <p role="alert">{RBAC_COPY.label5}</p>}
    {state.query.isSuccess && rows.length === 0 && <p>{RBAC_COPY.label6}</p>}
    <Box className={styles.list}>
      {rows.map((row) => <button className={styles.row} key={row.roleId} data-testid={RBAC_TEST_IDS.row} onClick={() => state.setSelected(row)}><strong>{row.name}</strong><span>{row.status}</span><span>{row.description}</span></button>)}
    </Box>
    <Box className={styles.toolbar}><Button disabled={state.page === 1 || state.query.isFetching} onClick={() => state.setPage(state.page - 1)}>{RBAC_COPY.label7}</Button><span>{state.page}</span><Button disabled={!state.query.data?.pageInfo.hasNext || state.query.isFetching} onClick={() => state.setPage(state.page + 1)}>{RBAC_COPY.label8}</Button></Box>
    {state.selected && <Sheet open title={RBAC_COPY.label9} dataTestId={RBAC_TEST_IDS.details} onOpenChange={(open) => { if (!open) state.setSelected(null); }}>
      <p><strong>{RBAC_COPY.label10}</strong> {state.selected.roleId}</p>
      <p><strong>{RBAC_COPY.label11}</strong> {state.selected.name}</p>
      <p><strong>{RBAC_COPY.label12}</strong> {state.selected.description}</p>
      <p><strong>{RBAC_COPY.label13}</strong> {state.selected.status}</p>
      <p><strong>{RBAC_COPY.label14}</strong> {state.selected.version}</p>
      <h2>{RBAC_COPY.label15}</h2>{state.permissions.isPending && <p>{RBAC_COPY.label4}</p>}{state.permissions.isError && <p role="alert">{RBAC_COPY.label16}</p>}{state.permissions.data?.permissionCodes.map((code) => <p key={code}>{code}</p>)}
    </Sheet>}
  </Box>;
};
