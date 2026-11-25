import { _map } from '@Base/utilities/map';
import { Box } from '../Box';
import type { ListProps } from './types/list.props';
import styles from './List.module.css';

export const List = <T,>({ items, getItemId, renderItem, label, emptyContent, dataTestId }: ListProps<T>) => (
  <Box role="list" aria-label={label} className={styles.list} dataTestId={dataTestId}>
    {_map(items, (item, index) => (
      <Box key={getItemId(item)} role="listitem" className={styles.item}>
        {renderItem(item, index)}
      </Box>
    ))}
    {!items.length && emptyContent}
  </Box>
);
