import { ChevronLeft, ChevronRight } from '@Components/icons';
import { _map } from '@Base/utilities/map';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Select } from '../Select';
import { Tooltip } from '../Tooltip';
import { PageJump } from './components/pageJump';

import { PaginationType } from './constants/pagination.type';
import styles from './Pagination.module.css';
import type { PaginationProps } from './types/pagination.props';

export const Pagination = (props: PaginationProps) => {
  const disabled = props.disabled || props.loading;
  if (props.type === PaginationType.INFINITE) {
    return <Button disabled={disabled || !props.hasNext} onClick={() => props.onChange({ type: 'LOAD_MORE' })}>{__('Load more')}</Button>;
  }
  const hasNext = () => {
    if (props.type === PaginationType.CURSOR || props.type === PaginationType.HYBRID) return props.hasNext;
    return props.currentPage * props.pageSize < props.totalRecords;
  };
  const hasPrevious = () => {
    if (props.type === PaginationType.CURSOR || props.type === PaginationType.HYBRID) return props.hasPrevious;
    return props.currentPage > 1;
  };
  const getTotalPages = () => {
    if (props.type === PaginationType.CURSOR) return undefined;
    if (props.type === PaginationType.HYBRID && props.totalPages !== undefined) return props.totalPages;
    if (props.totalRecords !== undefined) return Math.ceil(props.totalRecords / props.pageSize);
    return undefined;
  };
  const totalPages = getTotalPages();
  return (
    <Box className={styles.footer} dataTestId={props.dataTestId}>
      <Box className={styles.size}>
        <Box>{__('Rows per page')}</Box>
        <Select<number> hideLabel label={__('Rows per page')} value={props.pageSize} options={_map(props.pageSizeOptions ?? [10, 25, 50, 100], (size) => ({ value: size, label: String(size) }))} searchable={false} clearable={false} disabled={disabled} onChange={(pageSize) => {
          if (pageSize) props.onChange({ type: 'PAGE_SIZE_CHANGE', pageSize });
        }} />
      </Box>
      <Box className={styles.navigation}>
        <Tooltip label={!hasPrevious() && __('You are on the first page.') || undefined} focusable={!hasPrevious()}><ButtonWithIcon icon={ChevronLeft} iconSize={15} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Previous')} disabled={disabled || !hasPrevious()} onClick={() => props.onChange({ type: 'PREVIOUS' })} /></Tooltip>
        {totalPages !== undefined && <PageJump key={props.currentPage} currentPage={props.currentPage} totalPages={totalPages} disabled={Boolean(disabled)} onChange={(page) => props.onChange({ type: 'PAGE_CHANGE', page })} />}
        {totalPages === undefined && <Box className={styles.page} aria-live="polite" aria-label={__('Page {{page}}', { page: props.currentPage })}>{props.currentPage}</Box>}
        <Tooltip label={!hasNext() && __('There is no next page.') || undefined} focusable={!hasNext()}><ButtonWithIcon icon={ChevronRight} iconSize={15} size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.GHOST} aria-label={__('Next')} disabled={disabled || !hasNext()} onClick={() => props.onChange({ type: 'NEXT' })} /></Tooltip>
      </Box>
    </Box>
  );
};
