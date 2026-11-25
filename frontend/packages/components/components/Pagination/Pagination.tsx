import { ChevronLeft, ChevronRight } from '@Components/icons';
import { _map } from '@Base/utilities/map';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Select } from '../Select';

import { PaginationType } from './constants/pagination.type';
import styles from './Pagination.module.css';
import type { PaginationProps } from './types/pagination.props';

export const Pagination = (props: PaginationProps) => {
  const disabled = props.disabled || props.loading;
  if (props.type === PaginationType.INFINITE) {
    return <Button disabled={disabled || !props.hasNext} onClick={() => props.onChange({ type: 'LOAD_MORE' })}>{__('Load more')}</Button>;
  }
  const hasNext = () => {
    if (props.type === PaginationType.CURSOR) return props.hasNext;
    return props.currentPage * props.pageSize < props.totalRecords;
  };
  const hasPrevious = () => {
    if (props.type === PaginationType.CURSOR) return props.hasPrevious;
    return props.currentPage > 1;
  };
  return (
    <Box className={styles.footer} dataTestId={props.dataTestId}>
      <Box className={styles.size}>
        <Box>{__('Rows per page')}</Box>
        <Select<number> hideLabel label={__('Rows per page')} value={props.pageSize} options={_map(props.pageSizeOptions ?? [10, 25, 50, 100], (size) => ({ value: size, label: String(size) }))} searchable={false} clearable={false} disabled={disabled} onChange={(pageSize) => {
          if (pageSize) props.onChange({ type: 'PAGE_SIZE_CHANGE', pageSize });
        }} />
      </Box>
      <Box className={styles.navigation}>
        <ButtonWithIcon icon={ChevronLeft} iconSize={15} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Previous')} disabled={disabled || !hasPrevious()} onClick={() => props.onChange({ type: 'PREVIOUS' })} />
        <Box className={styles.page} aria-live="polite" aria-label={__('Page {{page}}', { page: props.currentPage })}>{props.currentPage}</Box>
        <ButtonWithIcon icon={ChevronRight} iconSize={15} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Next')} disabled={disabled || !hasNext()} onClick={() => props.onChange({ type: 'NEXT' })} />
      </Box>
    </Box>
  );
};
