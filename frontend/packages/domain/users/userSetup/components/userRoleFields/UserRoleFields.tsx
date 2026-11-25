import { _filter } from '@Base/utilities/filter';
import { Box } from '@Components/components/Box';
import { Select } from '@Components/components/Select';
import { MultiSelect } from '@Components/components/MultiSelect';

import styles from './UserRoleFields.module.css';
import type { UserRoleFieldsProps } from './types/userRoleFields.props';

export const UserRoleFields = ({ options, value, onChange, disabled }: UserRoleFieldsProps) => (
  <Box className={styles.fields}>
    <Select label={__('Primary role (required)')} value={value.primaryRoleId || null} options={options} clearable={false} disabled={disabled} onChange={(role) => {
      if (role) onChange({ primaryRoleId: role, secondaryRoleIds: _filter(value.secondaryRoleIds, (id) => id !== role) });
    }} />
    <MultiSelect label={__('Secondary roles')} value={value.secondaryRoleIds} options={_filter(options, (option) => option.value !== value.primaryRoleId)} disabled={disabled} onChange={(roles) => onChange({ ...value, secondaryRoleIds: roles })} />
  </Box>
);
