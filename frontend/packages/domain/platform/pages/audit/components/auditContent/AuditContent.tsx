import { Box } from '@Components/components/Box';
import { Sheet } from '@Components/components/Sheet';
import { Table, TableType, TableDataMode } from '@Components/components/Table';
import type { AuditControls } from '../../hooks/useAuditPage/types/useAuditPage';
import { AUDIT_COPY } from '../../constants/audit.copy';
import { AUDIT_TEST_IDS } from '../../constants/audit.testIds';
import { AUDIT_COLUMNS } from '../../constants/audit.columns';
import { AUDIT_FILTERS } from '../../constants/audit.filters';
import { AUDIT_SORT_FIELDS } from '../../constants/audit.sort';
import { renderAuditMobileRow } from '../../renderer/audit.mobileRow';
import styles from '../../page.module.css';

export const AuditContent = ({ controls: state }: { controls: AuditControls }) => {
  const getError = () => { if (state.error) return AUDIT_COPY.label6; return undefined; };
  return <Box className={styles.page} data-workspace-fill="true">
    <Table id={AUDIT_TEST_IDS.row} fillHeight type={TableType.DEFAULT} dataMode={TableDataMode.SERVER} data={state.rows} columns={AUDIT_COLUMNS} getRowId={(row) => row.auditId} headerProps={{ title: AUDIT_COPY.label1 }} searchProps={{ value: state.searchText, placeholder: AUDIT_COPY.label3 }} filterProps={{ fields: AUDIT_FILTERS, value: state.filters }} sortProps={{ fields: AUDIT_SORT_FIELDS, value: state.sorting }} paginationProps={state.pagination} mobileRowProps={{ render: renderAuditMobileRow, onClick: state.setSelected }} mobileInfiniteProps={state.mobileInfiniteProps} actions={state.actions} loading={state.loading} refreshing={state.refreshing} error={getError()} />
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
