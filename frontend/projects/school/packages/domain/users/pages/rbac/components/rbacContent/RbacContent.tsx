import { _map } from '@Base/utilities/map';
import { Typography, TYPOGRAPHY_TAG } from '@Components/components/Typography';
import { Box } from '@Components/components/Box';
import { Sheet } from '@Components/components/Sheet';
import { Table, TableType, TableDataMode } from '@Components/components/Table';
import type { RbacControls } from '../../hooks/useRbacPage/types/useRbacPage';
import { RBAC_COPY } from '../../constants/rbac.copy';
import { RBAC_TEST_IDS } from '../../constants/rbac.testIds';
import { RBAC_COLUMNS } from '../../constants/rbac.columns';
import { RBAC_FILTERS } from '../../constants/rbac.filters';
import { RBAC_SORT_FIELDS } from '../../constants/rbac.sort';
import { renderRbacMobileRow } from '../../renderer/rbac.mobileRow';
import styles from '../../page.module.css';

export const RbacContent = ({ controls: state }: { controls: RbacControls }) => {
  const getError = () => {
    if (state.error) return RBAC_COPY.label5;
    return undefined;
  };
  return (
    <Box className={styles.page} data-workspace-fill="true">
      <Table
        id={RBAC_TEST_IDS.row}
        fillHeight
        type={TableType.DEFAULT}
        dataMode={TableDataMode.SERVER}
        data={state.rows}
        columns={RBAC_COLUMNS}
        getRowId={(row) => row.roleId}
        headerProps={{ title: RBAC_COPY.label1 }}
        searchProps={{ value: state.searchText, placeholder: RBAC_COPY.label2 }}
        filterProps={{ fields: RBAC_FILTERS, value: state.filters }}
        sortProps={{ fields: RBAC_SORT_FIELDS, value: state.sorting }}
        paginationProps={state.pagination}
        mobileRowProps={{ render: renderRbacMobileRow, onClick: state.setSelected }}
        mobileInfiniteProps={state.mobileInfiniteProps}
        actions={state.actions}
        loading={state.loading}
        refreshing={state.refreshing}
        error={getError()}
      />
      {state.selected && (
        <Sheet
          open
          title={RBAC_COPY.label9}
          dataTestId={RBAC_TEST_IDS.details}
          onOpenChange={(open) => {
            if (!open) state.setSelected(null);
          }}
        >
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{RBAC_COPY.label10}</Typography>{' '}
            {state.selected.roleId}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{RBAC_COPY.label11}</Typography>{' '}
            {state.selected.name}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{RBAC_COPY.label12}</Typography>{' '}
            {state.selected.description}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{RBAC_COPY.label13}</Typography>{' '}
            {state.selected.status}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.P}>
            <Typography as={TYPOGRAPHY_TAG.STRONG}>{RBAC_COPY.label14}</Typography>{' '}
            {state.selected.version}
          </Typography>
          <Typography as={TYPOGRAPHY_TAG.H2}>{RBAC_COPY.label15}</Typography>
          {state.permissions.isPending && (
            <Typography as={TYPOGRAPHY_TAG.P}>{RBAC_COPY.label4}</Typography>
          )}
          {state.permissions.isError && (
            <Typography as={TYPOGRAPHY_TAG.P} role="alert">
              {RBAC_COPY.label16}
            </Typography>
          )}
          {_map(state.permissions.data?.permissionCodes ?? [], (code) => (
            <Typography as={TYPOGRAPHY_TAG.P} key={code}>
              {code}
            </Typography>
          ))}
        </Sheet>
      )}
    </Box>
  );
};
