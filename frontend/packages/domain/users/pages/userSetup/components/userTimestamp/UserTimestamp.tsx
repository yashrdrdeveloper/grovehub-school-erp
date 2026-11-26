import { getFormattedDate } from '@Base/localization/getFormattedDate';
import { Box } from '@Components/components/Box';

import type { UserTimestampProps } from './types/userTimestamp.props';
import styles from './UserTimestamp.module.css';

export const UserTimestamp = ({ value }: UserTimestampProps) => {
  const date = getFormattedDate(value, { includeTime: false });
  const time = getFormattedDate(value, { timeOnly: true });
  return (
    <Box className={styles.root} title={getFormattedDate(value)}>
      <Box className={styles.date}>{date}</Box>
      {time !== '—' && <Box className={styles.time}>{time}</Box>}
    </Box>
  );
};
