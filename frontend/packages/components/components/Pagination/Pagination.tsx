import { useState } from 'react';
import { ChevronLeft, ChevronRight } from '@Components/icons';
import { _map } from '@Base/utilities/map';

import { Box } from '../Box';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '../Button';
import { ButtonWithIcon } from '../ButtonWithIcon';
import { Select } from '../Select';
import { Tooltip } from '../Tooltip';
import { Popover } from '../Popover';
import { PageJump } from './components/pageJump';

import { PaginationType } from './constants/pagination.type';
import styles from './Pagination.module.css';
import type { PaginationProps } from './types/pagination.props';

export const Pagination = (props: PaginationProps) => {
  const [jumpOpen, setJumpOpen] = useState(false);
  const disabled = props.disabled || props.loading;
  if (props.type === PaginationType.INFINITE) {
    return (
      <Button
        disabled={disabled || !props.hasNext}
        onClick={() => props.onChange({ type: 'LOAD_MORE' })}
      >
        {__('Load more')}
      </Button>
    );
  }
  const hasNext = () => {
    if (props.type === PaginationType.CURSOR || props.type === PaginationType.HYBRID)
      return props.hasNext;
    return props.currentPage * props.pageSize < props.totalRecords;
  };
  const hasPrevious = () => {
    if (props.type === PaginationType.CURSOR || props.type === PaginationType.HYBRID)
      return props.hasPrevious;
    return props.currentPage > 1;
  };
  const getTotalPages = () => {
    if (props.type === PaginationType.CURSOR) return undefined;
    if (props.type === PaginationType.HYBRID && props.totalPages !== undefined)
      return props.totalPages;
    if (props.totalRecords !== undefined) return Math.ceil(props.totalRecords / props.pageSize);
    return undefined;
  };
  const totalPages = getTotalPages();
  return (
    <Box className={styles.footer} dataTestId={props.dataTestId}>
      <Box className={styles.size}>
        <Box>{__('Rows per page')}</Box>
        <Select<number>
          hideLabel
          label={__('Rows per page')}
          value={props.pageSize}
          options={_map(props.pageSizeOptions ?? [10, 25, 50, 100], (size) => ({
            value: size,
            label: String(size),
          }))}
          searchable={false}
          clearable={false}
          disabled={disabled}
          onChange={(pageSize) => {
            if (pageSize) props.onChange({ type: 'PAGE_SIZE_CHANGE', pageSize });
          }}
        />
      </Box>
      <Box className={styles.navigation}>
        {totalPages !== undefined && totalPages <= 1 && (
          <Box className={styles.page} aria-live="polite">
            {__('{{page}} of {{count}}', {
              page: Math.min(props.currentPage, totalPages),
              count: totalPages,
            })}
          </Box>
        )}
        {totalPages !== undefined && totalPages > 1 && (
          <Popover
            open={jumpOpen}
            onOpenChange={setJumpOpen}
            title={__('Go to page')}
            width={280}
            compact
            trigger={
              <Button
                unstyled
                className={styles.pageButton}
                disabled={disabled}
                aria-expanded={jumpOpen}
                aria-label={__('Go to page, {{page}} of {{count}}', {
                  page: props.currentPage,
                  count: totalPages,
                })}
                onClick={() => setJumpOpen(true)}
              >
                {__('{{page}} of {{count}}', { page: props.currentPage, count: totalPages })}
              </Button>
            }
          >
            <PageJump
              key={props.currentPage}
              currentPage={props.currentPage}
              totalPages={totalPages}
              disabled={Boolean(disabled)}
              onChange={(page) => {
                setJumpOpen(false);
                props.onChange({ type: 'PAGE_CHANGE', page });
              }}
            />
          </Popover>
        )}
        {totalPages === undefined && (
          <Box
            className={styles.page}
            aria-live="polite"
            aria-label={__('Page {{page}}', { page: props.currentPage })}
          >
            {props.currentPage}
          </Box>
        )}
        <Tooltip
          label={(!hasPrevious() && __('You are on the first page.')) || undefined}
          focusable={!hasPrevious()}
        >
          <ButtonWithIcon
            className={styles.arrow}
            icon={ChevronLeft}
            iconSize={14}
            size={BUTTON_SIZE.SMALL}
            variant={BUTTON_VARIANT.OUTLINE}
            aria-label={__('Previous')}
            disabled={disabled || !hasPrevious()}
            onClick={() => props.onChange({ type: 'PREVIOUS' })}
          />
        </Tooltip>
        <Tooltip
          label={(!hasNext() && __('There is no next page.')) || undefined}
          focusable={!hasNext()}
        >
          <ButtonWithIcon
            className={styles.arrow}
            icon={ChevronRight}
            iconSize={14}
            size={BUTTON_SIZE.SMALL}
            variant={BUTTON_VARIANT.OUTLINE}
            aria-label={__('Next')}
            disabled={disabled || !hasNext()}
            onClick={() => props.onChange({ type: 'NEXT' })}
          />
        </Tooltip>
      </Box>
    </Box>
  );
};
