import { useState } from 'react';

import { Box, BOX_TAG } from '../../../Box';
import { Button, BUTTON_SIZE, BUTTON_VARIANT } from '../../../Button';
import { InputControl, INPUT_CONTROL_SIZE } from '../../../InputControl';
import { Tooltip } from '../../../Tooltip';

import type { PageJumpProps } from './types/pageJump.props';
import styles from './PageJump.module.css';

export const PageJump = ({ currentPage, totalPages, disabled, onChange }: PageJumpProps) => {
  const [draft, setDraft] = useState(String(currentPage));
  const page = Number(draft);
  const valid = /^\d+$/.test(draft) && Number.isSafeInteger(page) && page >= 1 && page <= totalPages;
  const canJump = !disabled && valid && page !== currentPage;
  const getReason = () => {
    if (disabled) return __('Pagination is loading.');
    if (!totalPages) return __('No pages available.');
    if (!valid) return __('Enter a page from 1 to {{count}}.', { count: totalPages });
    if (page === currentPage) return __('You are already on this page.');
    return undefined;
  };
  const jump = () => {
    if (canJump) onChange(page);
  };
  return (
    <Box className={styles.root}>
      <Box as={BOX_TAG.SPAN}>{__('Page')}</Box>
      <InputControl
        type="number"
        min={1}
        max={Math.max(1, totalPages)}
        step={1}
        aria-label={__('Page number')}
        controlSize={INPUT_CONTROL_SIZE.SMALL}
        className={styles.input}
        disabled={disabled || !totalPages}
        value={draft}
        onChange={(event) => setDraft(event.currentTarget.value)}
        onKeyDown={(event) => {
          if (event.key === 'Enter') {
            event.preventDefault();
            jump();
          }
        }}
      />
      <Box as={BOX_TAG.SPAN} aria-live="polite">{__('of {{count}}', { count: totalPages })}</Box>
      <Tooltip label={getReason()} focusable={!canJump}>
        <Button size={BUTTON_SIZE.SMALL} variant={BUTTON_VARIANT.OUTLINE} disabled={!canJump} onClick={jump}>{__('Go')}</Button>
      </Tooltip>
    </Box>
  );
};
