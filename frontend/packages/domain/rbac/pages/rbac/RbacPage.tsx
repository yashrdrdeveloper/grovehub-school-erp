import { Box } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { Sheet } from '@Components/components/Sheet';
import { useRbacPage } from './hooks/useRbacPage';
import styles from './page.module.css';

export const RbacPage = () => {
  const state = useRbacPage();
  const rows = state.query.data?.items ?? [];
  return <Box className={styles.page} data-workspace-fill="true">
    <h1>{__('Roles and permissions')}</h1>
    
    <Box className={styles.toolbar}>
      <input data-testid="rbac-search" aria-label={__('Search')} placeholder={__('Search')} value={state.searchText} onChange={(event) => state.search(event.target.value)} />
      <Button onClick={() => void state.query.refetch()}>{__('Refresh')}</Button>
    </Box>
    {state.query.isPending && <p role="status">{__('Loading…')}</p>}
    {state.query.isError && <p role="alert">{__('Unable to load records. Check your access and try again.')}</p>}
    {state.query.isSuccess && rows.length === 0 && <p>{__('No records found')}</p>}
    <Box className={styles.list}>
      {rows.map((row) => <button className={styles.row} key={row.roleId} data-testid="rbac-row" onClick={() => state.setSelected(row)}><strong>{row.name}</strong><span>{row.status}</span><span>{row.description}</span></button>)}
    </Box>
    <Box className={styles.toolbar}><Button disabled={state.page === 1 || state.query.isFetching} onClick={() => state.setPage(state.page - 1)}>{__('Previous')}</Button><span>{state.page}</span><Button disabled={!state.query.data?.pageInfo.hasNext || state.query.isFetching} onClick={() => state.setPage(state.page + 1)}>{__('Next')}</Button></Box>
    {state.selected && <Sheet open title={__('Details')} dataTestId="rbac-details" onOpenChange={(open) => { if (!open) state.setSelected(null); }}>
      <p><strong>{__('roleId')}</strong> {state.selected.roleId}</p>
      <p><strong>{__('name')}</strong> {state.selected.name}</p>
      <p><strong>{__('description')}</strong> {state.selected.description}</p>
      <p><strong>{__('status')}</strong> {state.selected.status}</p>
      <p><strong>{__('version')}</strong> {state.selected.version}</p>
      <h2>{__('Permissions')}</h2>{state.permissions.isPending && <p>{__('Loading…')}</p>}{state.permissions.isError && <p role="alert">{__('Unable to load permissions')}</p>}{state.permissions.data?.permissionCodes.map((code) => <p key={code}>{code}</p>)}
    </Sheet>}
  </Box>;
};
