import { Typography, TYPOGRAPHY_TAG } from '@Components/components/Typography';
import { Button } from '@Components/components/Button';
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
  const getError = () => {
    if (state.error) return AUDIT_COPY.label6;
    return undefined;
  };
  return (
    <Box className={styles.page} data-workspace-fill="true">
      {state.canViewTenant && (
        <Box>
          <Button
            dataTestId="audit-personal-scope"
            disabled={!state.tenantScope}
            onClick={() => state.changeScope(false)}
          >
            {AUDIT_COPY.personalScope}
          </Button>
          <Button
            dataTestId="audit-tenant-scope"
            disabled={state.tenantScope}
            onClick={() => state.changeScope(true)}
          >
            {AUDIT_COPY.tenantScope}
          </Button>
        </Box>
      )}
      <Table
        id={AUDIT_TEST_IDS.row}
        fillHeight
        type={TableType.DEFAULT}
        dataMode={TableDataMode.SERVER}
        data={state.rows}
        columns={AUDIT_COLUMNS}
        getRowId={(row) => row.auditId}
        headerProps={{ title: AUDIT_COPY.label1 }}
        searchProps={{ value: state.searchText, placeholder: AUDIT_COPY.label3 }}
        filterProps={{ fields: AUDIT_FILTERS, value: state.filters }}
        sortProps={{ fields: AUDIT_SORT_FIELDS, value: state.sorting }}
        paginationProps={state.pagination}
        mobileRowProps={{ render: renderAuditMobileRow, onClick: state.setSelected }}
        mobileInfiniteProps={state.mobileInfiniteProps}
        actions={state.actions}
        loading={state.loading}
        refreshing={state.refreshing}
        error={getError()}
      />
      {state.selected && (
        <Sheet
          open
          title={AUDIT_COPY.label10}
          dataTestId={AUDIT_TEST_IDS.details}
          onOpenChange={(open) => {
            if (!open) state.setSelected(null);
          }}
        >
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{AUDIT_COPY.label11}</Typography>{' '}
            {state.selected.auditId}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{AUDIT_COPY.label12}</Typography>{' '}
            {state.selected.source}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{AUDIT_COPY.label13}</Typography>{' '}
            {state.selected.assetType}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{AUDIT_COPY.label14}</Typography>{' '}
            {state.selected.assetId}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{AUDIT_COPY.label15}</Typography>{' '}
            {state.selected.operation}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{AUDIT_COPY.label16}</Typography>{' '}
            {state.selected.occurredAt}
          </Typography>
          <pre className={styles.payload}>{JSON.stringify(state.selected.data, null, 2)}</pre>
        </Sheet>
      )}
    </Box>
  );
};
